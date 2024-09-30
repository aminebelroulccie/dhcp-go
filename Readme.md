# Nex: Scalable Name and Address Software Infrastructure

Nex is the spirit of dnsmasq meets modern scale out architecture. The design 
is inspired and motivated by the recent success Facebook has had with deploying
Kea using a stateless replicated model. Nex is based on a similar stateless 
service model, built on a simple Go based DHCP server and CoreDNS - with an 
etcd deployment that brings them together into a cohesive system - ala dnsmasq.

The DNS side of things is a fork of CoreDNS plugin and is maintained in this 
repository

https://github.com/mergetb/coredns

## Getting Started

The first thing to do with Nex is describe a starting point for your network.
The most convenient way to do this is through YAML.

```yaml
kind:         Network
name:         mini
subnet4:      10.0.0.0/24
gateways:     [10.0.0.1]
nameservers:  [10.0.0.1]
dhcp4server:  10.0.0.1
domain:       mini.net
range4:
  begin: 10.0.0.10
  end:   10.0.0.254

---
kind:   MemberList
net:    mini
list:
- mac:  00:00:11:11:00:01
  name: whiskey

- mac:  00:00:22:22:00:01
  name: tango

- mac:  00:00:33:33:00:01
  name: foxtrot
```
This document contains two distinct entities. The `Network` describes the basic
properties of a DHCP/DNS network and the `MemberList` describes the devices that
are initially a part of that network. This is just a starting point, devices can
be added or removed at any time, and network properties changed.

This spec can be applied using the Nex command line utility called `nex`.

```shell
$ nex apply mini.yml
```

```shell
$ nex get network mini

name:           mini
subnet4:        10.0.0.0/24
gateways:       10.0.0.1
nameservers:    10.0.0.1
dhcp4server:    10.0.0.1
domain:         mini.net
range4:         10.0.0.10-10.0.0.254
```

```shell
$ nex get members mini

00:00:11:11:00:01    whiskey.mini.net    10.0.0.10
00:00:22:22:00:01    tango.mini.net      10.0.0.11
00:00:33:33:00:01    foxtrot.mini.net    10.0.0.12
```

When devices request DHCP addresses you'll see them show up with the time left
on the lease.

```shell
$ nex get members mini

00:00:11:11:00:01    whiskey.mini.net    10.0.0.10 (3:53:16)
00:00:22:22:00:01    tango.mini.net      10.0.0.11 (3:53:17)
00:00:33:33:00:01    foxtrot.mini.net    10.0.0.12 (3:53:19)
```

## API

Nex has a [gRPC](https://grpc.io) API that allows you to manage members and
networks programmatically.

```protobuf
service Nex {
  /* membership */
  rpc GetMembers(GetMembersRequest) returns (GetMembersResponse);
  rpc AddMembers(MemberList) returns (AddMembersResponse);
  rpc DeleteMembers(DeleteMembersRequest) returns (DeleteMembersResponse);
  rpc UpdateMembers(UpdateList) returns (UpdateMembersResponse);

  /* network */
  rpc GetNetworks(GetNetworksRequest) returns (GetNetworksResponse);
  rpc GetNetwork(GetNetworkRequest) returns (GetNetworkResponse);
  rpc AddNetwork(AddNetworkRequest) returns (AddNetworkResponse);
  rpc UpdateNetwork(NetworkUpdateRequest) returns (NetworkUpdateResponse);
  rpc DeleteNetwork(DeleteNetworkRequest) returns (DeleteNetworkResponse);
}
```
For details see [nex.proto](pkg/nex.proto).

## Scaling and fault tolerance

The Nex storage backend is [etcd](https://etcd.io). You can run multiple
instances of the Nex API, DHCP server and DNS server for scalability or fault
tolerance purposes. All of the API calls are designed for replicated operations
and provide atomic semantics. The system was designed to be run with DHCP relays
anycasting to DHCP server replicas.

## Configuration

Nex requires you to tell it what it should listen to where to find etcd. An
example configuration is the following.

```yaml
dhcpd:
  interface: eth1

etcd:
  host:   db
  port:   2379
  cert:   /etc/nex/db.pem
  key:    /etc/nex/db-key.pem
  cacert: /etc/nex/ca.pem

nexd:
  listen: 0.0.0.0:6000
```

## Installation

Installation automation through an Ansible role is available
[here](https://gitlab.com/mergetb/ansible/nex).

## Hacking

Testing systems are located in tests. Currently only `little` works. These tests
use a technology called [raven](https://gitlab.com/rygoo/raven) the authors
designed for easily materializing testing environments for networked systems.

To get started

```shell
# fetch dependencies and build DHCP bits
go install github.com/golang/protobuf/protoc-gen-go
go mod vendor
make distclean
make

# build DNS bits
git submodule init
git submodule update
cd coredns
make clean
make
cp coredns ../build/

# run the basic test
cd tests/little/test
sudo ./run.sh
```
