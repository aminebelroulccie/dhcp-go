#!/bin/bash

docker run \
  -d -it --rm \
  --name nex \
  -e DOMAIN=test.net \
  -p 6000:6000 -p 5533:53/udp \
  docker.io/mergetb/nex

