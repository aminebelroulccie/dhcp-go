#!/bin/bash

set -x
set -e

if [[ $(id -u) -ne 0 ]]; then
  echo "must be root to run this script"
  exit 1;
fi

testdir=`pwd`
topdir="$testdir/../../.."
rvndir="$testdir/.."

echo "building software"
cd $topdir
docker build -t nex-builder -f builder.dock --no-cache .
docker run -v "$topdir/build:/build" nex-builder

echo "deploying and configuring raven topology"
cd $rvndir
rvn destroy
rvn build
rvn deploy
rvn pingwait i0 i1 e0 e1 e2 v0 v1 s0 s1 db sw
rvn configure db
rvn configure i0 i1 e0 e1 e2 v0 v1 s0 s1 sw

./test0/runtests.sh

rvn destroy
rm -rf .rvn
