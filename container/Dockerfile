FROM debian:buster

RUN apt-get update && apt-get install -y supervisor curl procps

COPY build/nexc /usr/bin/nexc
COPY build/nexd /usr/bin/nexd
COPY build/nex-dhcpd /usr/bin/nex-dhcpd
COPY coredns/coredns /usr/bin/coredns
COPY container/supervisord.conf /etc/supervisord.conf
COPY container/nex.yml /etc/nex/nex.yml
COPY container/hosts /etc/hosts
COPY container/launchdns.sh /usr/bin/launchdns

RUN curl -o /tmp/etcd-v3.3.11-linux-amd64.tar.gz -L \
https://github.com/etcd-io/etcd/releases/download/v3.3.11/etcd-v3.3.11-linux-amd64.tar.gz && \
cd /tmp && \
tar -xzf etcd-v3.3.11-linux-amd64.tar.gz && \
cp etcd-v3.3.11-linux-amd64/etcd /usr/bin/etcd && \
cp etcd-v3.3.11-linux-amd64/etcdctl /usr/bin/etcdctl

RUN mkdir /data

# nex API
EXPOSE 6000 

# dhcp-server
EXPOSE 67   

# dns
EXPOSE 53   

CMD supervisord -n
