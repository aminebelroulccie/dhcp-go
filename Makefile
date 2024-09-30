prefix ?= /usr/local

artifacts = nexd nexc nex-dhcpd coredns
common=pkg/nex.pb.go pkg/*.go

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

VERSION = $(shell git describe --always --long --dirty)
LDFLAGS = "-X gitlab.com/mergetb/tech/nex/pkg.Version=$(VERSION)"

REGISTRY ?= docker.io
ORG ?= mergetb
TAG ?= $(shell git describe --always)

export PATH := $(HOME)/go/bin:$(PATH)

.PHONY: all
all: $(patsubst %,build/%,$(artifacts)) build/basic-test

.tools:
	$(QUIET) mkdir .tools

protoc-gen-go=.tools/protoc-gen-go
$(protoc-gen-go): | .tools
	$(QUIET) GOBIN=`pwd`/.tools go install github.com/golang/protobuf/protoc-gen-go

$(HOME)/go/bin/protoc-gen-go:
	go install github.com/golang/protobuf/protoc-gen-go

build/nexd: svc/nexd/*.go $(common) | build
	$(go-build)

build/nexc: util/nexc/*.go $(common) | build
	$(go-build)

build/nex-dhcpd: svc/dhcp-server/*.go $(common) | build
	$(go-build)

coredns/plugin/nex/nex.go:
	$(QUIET) git submodule init
	$(QUIET) git submodule update

build/coredns: coredns/plugin/nex/nex.go $(common)
	$(QUIET) cd coredns && go build
	$(QUIET) cp coredns/coredns build/coredns

vendor:
	$(QUIET) go mod vendor

pkg/nex.pb.go: pkg/nex.proto $(protoc-gen-go) vendor
	$(protoc-build)
	$(QUIET) sed -r -i 's/json:"(.*)"/json:"\1" yaml:"\1" mapstructure:"\1"/g' pkg/nex.pb.go


.PHONY: $(ORG)/nex
$(REGISTRY)/$(ORG)/nex: container/Dockerfile build/nexd build/nexc build/nex-dhcpd
	$(docker-build)

.PHONY: container
container: $(REGISTRY)/$(ORG)/nex


build/basic-test: tests/basic/test.go
	$(go-build-file)

build:
	@mkdir -p build

.PHONY: clean
clean:
	$(QUIET) rm -rf build

.PHONY: distclean
distclean: clean
	$(QUIET) find . \
		\( -path ./vendor -o -path ./coredns \) -prune -o -name "*.pb.go" -print \
		| xargs -n 1 rm -f
	$(QUIET) rm -rf .tools
	$(QUIET) rm -rf vendor

.PHONY: install-heavy
install-heavy: $(patsubst %,build/%,$(artifacts))
	$(QUIET) mkdir -p $(prefix)/bin
	$(QUIET) cp build/nexd $(prefix)/bin/nexd
	$(QUIET) cp build/nexc $(prefix)/bin/nexc
	$(QUIET) cp build/nex-dhcpd $(prefix)/bin/nex-dhcpd
	$(QUIET) mkdir -p ${prefix}/../lib/systemd/system
	$(QUIET) cp coredns.service ${prefix}/../lib/systemd/system/coredns.service
	$(QUIET) cp nexd.service ${prefix}/../lib/systemd/system/nexd.service
	$(QUIET) cp nex-dhcpd.service ${prefix}/../lib/systemd/system/nex-dhcpd.service
	$(QUIET) mkdir -p ${prefix}/../etc/merge

.PHONY: install
install: $(patsubst %,build/%,$(artifacts))
	$(QUIET) install -D build/nexd $(DESTDIR)$(prefix)/bin/nexd
	$(QUIET) install -D build/nexc $(DESTDIR)$(prefix)/bin/nex
	$(QUIET) install -D build/nex-dhcpd $(DESTDIR)$(prefix)/bin/nex-dhcpd
	$(QUIET) install -D build/coredns $(DESTDIR)$(prefix)/bin/coredns
	$(QUIET) install -D config/nex.yml $(DESTDIR)/etc/nex/nex.yml



BLUE=\e[34m
GREEN=\e[32m
CYAN=\e[36m
NORMAL=\e[39m

me=$(shell whoami)

.PHONY: cleanbuild
cleanbuild: distclean
	docker run \
		-v `pwd`:/go/src/gitlab.com/mergetb/tech/nex \
		-w /go/src/gitlab.com/mergetb/tech/nex \
		-e "GOPATH=/go" \
		golang:1.11-stretch \
		make build-from-scratch
	sudo chown $(me):$(me) -R .

build-from-scratch:
	apt-get update && apt-get -y install protobuf-compiler
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
	go install ./vendor/github.com/golang/protobuf/protoc-gen-go
	make


QUIET=@
DOCKER_QUIET=-q
ifeq ($(V),1)
	QUIET=
	DOCKER_QUIET=
endif

define build-slug
	@echo "$(BLUE)$1$(GREEN)\t $< $(CYAN)$@$(NORMAL)"
endef

define go-build
	$(call build-slug,go)
	$(QUIET) go build -ldflags=$(LDFLAGS) -o $@ $(dir $<)/*
endef

define go-build-file
	$(call build-slug,go)
	$(QUIET) go build -ldflags=$(LDFLAGS) -o $@ $<
endef

define protoc-build
	$(call build-slug,protoc)
	$(QUIET) PATH=./.tools:$$PATH protoc \
		-I . \
		-I ./vendor \
		-I ./$(dir $@) \
		./$< \
		--go_out=plugins=grpc:.
endef

define docker-build
	$(call build-slug,docker)
	$(QUIET) docker build ${BUILD_ARGS} $(DOCKER_QUIET) -f $< -t $(@):$(TAG) .
	$(if ${PUSH},$(call docker-push))
endef

define docker-push
	$(call build-slug,push)
	$(QUIET) docker push $(@):$(TAG)
endef
