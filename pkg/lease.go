package nex

import (
	"net"
	"time"

	dhcp "github.com/krolaw/dhcp4"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Pool struct {
	CountSet
	Net string
}

func NewLease4(mac net.HardwareAddr, network string, pkt dhcp.Packet) (net.IP, error) {

	netW := NewNetworkObj(&Network{Name: network})
	err := Read(netW)
	if err != nil {
		return nil, err
	}

	// pool

	pool := NewPoolObj(&Pool{Net: network})
	err = Read(pool)
	if err != nil {
		return nil, err
	}
	log.Debugf("poolver %d", pool.version)

	index, cs, err := pool.CountSet.Add()
	if err != nil {
		return nil, nil
	}
	pool.CountSet = cs
	ip := netW.Range4.Select(index)
	//if ip.String() == net.
	// mac
	m := NewMacIndex(&Member{
		Mac: mac.String(),
		Net: network,
	})
	err = ReadNew(m)
	if err != nil {
		return nil, err
	}

	// ip4
	ip4 := NewIp4Index(m.Member)
	err = ReadNew(ip4)
	if err != nil {
		return nil, err
	}

	// net
	neti := NewNetIndex(m.Member)
	err = ReadNew(neti)
	if err != nil {
		return nil, err
	}

	opts := pkt.ParseOptions()
	hostname, hasHostname := opts[dhcp.OptionHostName]
	if hasHostname {
		m.ClientName = string(hostname)
	}

	expires := time.Now().Add(time.Duration(netW.LeaseDuration) * time.Second)
	m.Ip4 = &Lease{
		Address: ip.String(),
		Expires: &timestamppb.Timestamp{
			Seconds: expires.Unix(),
			Nanos:   0,
		},
	}

	err = WriteObjects([]Object{pool, m, ip4, neti})
	if err != nil {
		return nil, err
	}

	return ip, err

}

func RenewLease(mac net.HardwareAddr, network string) error {

	m := &Member{Mac: mac.String()}
	obj := NewMacIndex(m)
	err := Read(obj)
	if err != nil {
		return nil
	}

	net := NewNetworkObj(&Network{Name: network})
	err = Read(net)
	if err != nil {
		return err
	}

	// Nothing to do
	if m.Ip4 == nil {
		return nil
	}

	// static lease needs no renewal, bug out early
	if m.Ip4.Expires == nil {
		return nil
	}

	expires := time.Now().Add(time.Duration(net.LeaseDuration) * time.Second)
	m.Ip4.Expires.Seconds = expires.Unix()
	m.Ip4.Expires.Nanos = 0

	return Write(obj)

}
