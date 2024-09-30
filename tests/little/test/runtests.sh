c0=`rvn ip c0`
c1=`rvn ip c1`
c2=`rvn ip c2`
c3=`rvn ip c3`

echo "running tests"
avocado run \
  --remote-username rvn \
  --remote-password rvn \
  --remote-hostname $c0 \
  /tmp/nex/tests/little/test/c0.py

avocado run \
  --remote-username rvn \
  --remote-password rvn \
  --remote-hostname $c1 \
  /tmp/nex/tests/little/test/c1.py

avocado run \
  --remote-username rvn \
  --remote-password rvn \
  --remote-hostname $c2 \
  /tmp/nex/tests/little/test/c2.py

avocado run \
  --remote-username rvn \
  --remote-password rvn \
  --remote-hostname $c3 \
  /tmp/nex/tests/little/test/c3.py
