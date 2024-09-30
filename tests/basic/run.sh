#!/bin/bash

set -e

if [[ $UID -ne 0 ]]; then
  echo "must be root"
  exit 1
fi

BLUE="\e[34m"
BLINK="\e[5m"
RESET="\e[0m"

stage() {
  echo -e "$BLUE$1$BLINK 🔨$RESET"
}


stage "destroying any existing topologies"
rvn destroy

stage "building topology"
rvn build

stage "deploying topology"
rvn deploy

stage "waiting for topology to come up"
rvn pingwait server tango foxtrot

stage "running base configuration"
rvn configure
rvn status

stage "setting up nex"
ansible-playbook -i .rvn/ansible-hosts setup.yml

stage "run test"
ansible-playbook -i .rvn/ansible-hosts test.yml
