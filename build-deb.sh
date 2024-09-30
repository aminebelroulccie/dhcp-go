#!/bin/bash

TARGET=${TARGET:-amd64}
DEBUILD_ARGS=${DEBUILD_ARGS:-""}

rm -f build/nex*.build*
rm -f build/nex*.change
rm -f build/nex*.deb

debuild -e V=1 -e prefix=/usr $DEBUILD_ARGS -i -us -uc -b

mv ../nex*.build* build/
mv ../nex*.changes build/
mv ../nex*.deb build/
