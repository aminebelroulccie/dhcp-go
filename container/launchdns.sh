#!/bin/bash

mkdir -p /etc/coredns

cat << EOF > /etc/coredns/Corefile
$DOMAIN {
  nex
}

. {
  forward . 8.8.8.8
}
EOF

/usr/bin/coredns -conf /etc/coredns/Corefile -dns.port 53
