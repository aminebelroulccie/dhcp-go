#!/bin/bash

set -e

docker build $BUILD_ARGS -f debian/builder.dock -t nex-builder .
docker run -v `pwd`:/nex nex-builder /nex/build-deb.sh
