i0_ip=`rvn ip i0`
i1_ip=`rvn ip i1`
e0_ip=`rvn ip e0`
e1_ip=`rvn ip e1`
e2_ip=`rvn ip e2`
v0_ip=`rvn ip v0`
v1_ip=`rvn ip v1`

echo "running tests"
avocado run \
  --remote-username rvn \
  --remote-password rvn \
  --remote-hostname $i0_ip \
  /tmp/nex/tests/basic/test0/i0.py

avocado run \
  --remote-username rvn \
  --remote-password rvn \
  --remote-hostname $i1_ip \
  /tmp/nex/tests/basic/test0/i1.py

avocado run \
  --remote-username rvn \
  --remote-password rvn \
  --remote-hostname $e0_ip \
  /tmp/nex/tests/basic/test0/e0.py

avocado run \
  --remote-username rvn \
  --remote-password rvn \
  --remote-hostname $e1_ip \
  /tmp/nex/tests/basic/test0/e1.py

avocado run \
  --remote-username rvn \
  --remote-password rvn \
  --remote-hostname $e2_ip \
  /tmp/nex/tests/basic/test0/e2.py

avocado run \
  --remote-username rvn \
  --remote-password rvn \
  --remote-hostname $v0_ip \
  /tmp/nex/tests/basic/test0/v0.py

avocado run \
  --remote-username rvn \
  --remote-password rvn \
  --remote-hostname $v1_ip \
  /tmp/nex/tests/basic/test0/v1.py
