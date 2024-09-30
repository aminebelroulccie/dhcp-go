package options

import (
	"fmt"
	"net"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
	log "github.com/sirupsen/logrus"
)

type RawConn struct {
	conn  *raw.Conn
	src   net.HardwareAddr
	SrcIP net.IP
	Iface string
}

func NewRawListener(ifx string) (*RawConn, error) {

	ifi, err := net.InterfaceByName(ifx)
	if err != nil {
		return nil, err
	}

	addrs, err := ifi.Addrs()
	if err != nil {
		return nil, err
	}

	var srcIP net.IP = nil
	for _, x := range addrs {
		srcip, _, err := net.ParseCIDR(x.String())
		if err != nil {
			continue
		}
		srcip = srcip.To4()
		if srcip == nil {
			continue
		}
		srcIP = srcip
		break
	}
	if srcIP == nil {
		return nil, fmt.Errorf("interface does not have ipv4 address")
	}

	log.Infof("using source address %s", srcIP)

	c, err := raw.ListenPacket(ifi, uint16(ethernet.EtherTypeIPv4), nil)
	if err != nil {
		return nil, err
	}
	// log.Info(c.LocalAddr().String())
	return &RawConn{conn: c, src: ifi.HardwareAddr, SrcIP: srcIP, Iface: ifx}, nil

}


func(c *RawConn) Close() {
	c.conn.Close()
}

func (c *RawConn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {

	_b := make([]byte, len(b))

	for {

		n, addr, err := c.conn.ReadFrom(_b)
		if err != nil {
			return 0, nil, err
		}
		// log.Info(addr)
		pkt := gopacket.NewPacket(_b[:n], layers.LayerTypeEthernet, gopacket.Default)
		dhcpLayer := pkt.Layer(layers.LayerTypeDHCPv4)
		if dhcpLayer != nil {
            
			
			ipv4, ok := pkt.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
			if !ok {
				return 0, nil, fmt.Errorf("non-ipv4 dhcp packet")
			}

			udp, ok := pkt.Layer(layers.LayerTypeUDP).(*layers.UDP)
			if !ok {
				return 0, nil, fmt.Errorf("non-udp dhcp packet")
			}

			addr = &net.UDPAddr{IP: ipv4.SrcIP, Port: int(udp.SrcPort)}
			// log.Info(c.srcIP)
			// offset = len(eth_hdr) + len(ipv4_hdr) + len(udp_hdr)
			//        = 14 + 20 + 8
			//        = 42
			offset := 42

			copy(b, _b[offset:])
			return n - offset, addr, err

		}

	}

	return 0, nil, nil

}

func (c *RawConn) WriteTo(b []byte, addr net.Addr) (n int, err error) {

	pkt := gopacket.NewPacket(b, layers.LayerTypeDHCPv4, gopacket.Default)
	dhcpLayer := pkt.Layer(layers.LayerTypeDHCPv4)
	if dhcpLayer == nil {
		return 0, fmt.Errorf("not a dhcp packet")
	}
	dhcp, _ := dhcpLayer.(*layers.DHCPv4)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	l2 := &layers.Ethernet{
		SrcMAC:       c.src,
		DstMAC:       dhcp.ClientHWAddr,
		EthernetType: layers.EthernetTypeIPv4,
	}
	l3 := &layers.IPv4{
		Version:  4,
		IHL:      5,
		TOS:      16,
		TTL:      128,
		Protocol: layers.IPProtocolUDP,
		SrcIP:    c.SrcIP,
		DstIP:    net.ParseIP(strings.Split(addr.String(), ":")[0]),
	}
	l4 := &layers.UDP{
		SrcPort: 67,
		DstPort: 68,
	}

	l4.SetNetworkLayerForChecksum(l3)

	err = gopacket.SerializeLayers(buf, opts, l2, l3, l4, gopacket.Payload(b))
	if err != nil {
		log.Error(err)
		return 0, err
	}

	out := buf.Bytes()
	pkt = gopacket.NewPacket(out, layers.LayerTypeEthernet, gopacket.Default)

	return c.conn.WriteTo(out, &raw.Addr{HardwareAddr: dhcp.ClientHWAddr})

}
