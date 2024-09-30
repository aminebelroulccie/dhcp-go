package nex

import (
	"encoding/binary"
	"net"
)

func (a *AddressRange) Size() int {
	begin := binary.BigEndian.Uint32([]byte(net.ParseIP(a.Begin).To4()))
	end := binary.BigEndian.Uint32([]byte(net.ParseIP(a.End).To4()))
	return int(end - begin)
}

func (a *AddressRange) Select(offset int) net.IP {
	begin := binary.BigEndian.Uint32([]byte(net.ParseIP(a.Begin).To4()))
	chosen := begin + uint32(offset)

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, chosen)
	return net.IP(buf)
}

func (a *AddressRange) Offset(ip net.IP) int {
	begin := binary.BigEndian.Uint32([]byte(net.ParseIP(a.Begin).To4()))
	target := binary.BigEndian.Uint32([]byte(ip.To4()))
	return int(target - begin)
}
