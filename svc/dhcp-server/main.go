package dhcpd

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	dhcp "github.com/krolaw/dhcp4"
	log "github.com/sirupsen/logrus"

	nex "gitlab.com/mergetb/tech/nex/pkg"
	"gitlab.com/mergetb/tech/nex/svc/dhcp-server/options"
	"gitlab.com/mergetb/tech/nex/svc/nexd"
)

// Handler handler struct for dhcp server
type Handler struct {
	ip  net.IP
	Ctx context.Context
}

// var (
// 	nexD = options.DefaultNexD
// )

type DhcpD struct {
	nexD        *nexd.NexD
	ifacesChans map[string]*options.RawConn
	mu          *sync.Mutex
}

func New() *DhcpD {
	return &DhcpD{
		nexD:        nexd.DefaultNexD,
		ifacesChans: make(map[string]*options.RawConn),
		mu:          &sync.Mutex{},
	}
}

var DefaultDhcpD = New()

func (d *DhcpD) RunDelete() {
	for ifx := range d.nexD.DelChan {
		fmt.Println(ifx.Name)
		if conn, ok := d.ifacesChans[ifx.Name]; ok {
			conn.Close()
			d.mu.Lock()
			delete(d.ifacesChans, ifx.Name)
			d.mu.Unlock()
		}
	}
}

func (d *DhcpD) Run() {

	log.SetLevel(log.InfoLevel)
	log.Infof("nex-dhcpd: %s", nex.Version)

	// nex.Init()
	go d.RunDelete()
	err := nex.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// if nex.Current.Debug {
	// 	log.SetLevel(log.DebugLevel)
	// }
	// if nex.Current.Trace {
	// 	log.SetLevel(log.TraceLevel)
	// }
	// ifaces, err := nex.GetInterfaces()
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }
	// for _, ifx := range ifaces {
	// 	go d.listenDhcp(ifx, false)
	// }
	for ifx := range d.nexD.AddChan {
		go d.listenDhcp(ifx)

	}
}

func (d *DhcpD) listenDhcp(ifx *nex.InterfaceRequest) {
	cnx, err := options.NewRawListener(ifx.Name)
	if err != nil {
		log.Error(err)
		return
	}
	handler := &Handler{
		ip: cnx.SrcIP,
	}

	d.mu.Lock()
	d.ifacesChans[ifx.Name] = cnx
	d.mu.Unlock()
	// if write {
	// 	nex.Write(&nex.InterfaceObj{
	// 		InterfaceRequest: ifx,
	// 	})
	// }
	if err := dhcp.Serve(cnx, handler); err != nil {
		log.Error(err)
		return
	}
}

// ServeDHCP dhcp server
func (h *Handler) ServeDHCP(
	pkt dhcp.Packet,
	msgType dhcp.MessageType,
	options dhcp.Options,
) dhcp.Packet {

	fields := log.Fields{}
	// for {
	// for {
	// 	select {
	// 	case <- h.:
	// 		// The context is over, stop processing results
	// 		return
	// 	case result := <- resultsCh:
	// 		// Process the results received
	// 	}
	// }
	switch msgType {

	case dhcp.Discover:
		log.Info("discover")
		fields["mac"] = pkt.CHAddr()
		log.WithFields(fields).Debug("discover: start")

		// Collect network information
		network, err := nex.FindMacNetwork(pkt.CHAddr(), h.ip)
		if err != nil {
			log.WithError(err).Error("discover: error")
			return nil
		}
		if network == nil {
			log.WithFields(fields).Warn("discover: has no net")
			return nil
		}
		fields["network"] = network.Name
		log.WithFields(fields).Debug("found network")

		response := func(server, addr net.IP, options []dhcp.Option) dhcp.Packet {

			fields["addr"] = addr
			log.WithFields(fields).Debug("discover: OK")

			opts := ToOpt(network.Options)

			_, subnetCIDR, err := net.ParseCIDR(network.Subnet4)
			if err != nil {
				log.Errorf("bad subnet: %v", err)
				return nil
			}

			opts = append(opts, dhcp.Option{
				Code:  dhcp.OptionSubnetMask,
				Value: subnetCIDR.Mask,
			})

			for _, x := range network.Gateways {
				ip := net.ParseIP(x)
				opts = append(opts, dhcp.Option{
					Code:  dhcp.OptionRouter,
					Value: ip.To4(),
				})
			}

			for _, x := range network.Nameservers {
				ip := net.ParseIP(x)
				opts = append(opts, dhcp.Option{
					Code:  dhcp.OptionDomainNameServer,
					Value: ip.To4(),
				})
			}

			if network.Domain != "" {
				cn := compressedDNSName(network.Domain)
				opts = append(opts, dhcp.Option{
					Code:  dhcp.OptionDomainSearch,
					Value: cn,
				})
			}

			return dhcp.ReplyPacket(
				pkt, dhcp.Offer, server, addr,
				time.Duration(network.LeaseDuration)*time.Second,
				opts)
		}

		// If there is already an address use that
		addr, err := nex.FindMacIpv4(pkt.CHAddr(), network.Name)
		if err != nil && !nex.IsNotFound(err) {
			log.WithError(err).Errorf("discover: mac error")
			return nil
		}
		if addr != nil {
			fields["found"] = addr.String()

			return response(
				net.ParseIP(network.Dhcp4Server).To4(),
				addr,
				ToOpt(network.Options),
			)
		}

		// If no address was found allocate a new one
		addr, err = nex.NewLease4(pkt.CHAddr(), network.Name, pkt)
		if err != nil {
			log.WithError(err).Errorf("discover: lease error")
			return nil
		}
		if addr != nil {
			return response(
				net.ParseIP(network.Dhcp4Server).To4(),
				addr,
				ToOpt(network.Options),
			)
		}

		log.WithFields(fields).Error("address pool depleted")
		//Address is nil, so no discover response will go out

	case dhcp.Request:
		log.Info("request")
		fields["mac"] = pkt.CHAddr()
		log.WithFields(fields).Debug("request: start")
		// options[dhcp.OptionRequestedIPAddress]
		rqAddr := net.IP(pkt.CIAddr())

		network, err := nex.FindMacNetwork(pkt.CHAddr(), h.ip)
		if err != nil {
			log.WithError(err).Error("request: find mac net error")
			return nil
		}
		if network == nil {
			log.WithFields(fields).Warn("request: has no net")
			return nil
		}
		server := net.ParseIP(network.Dhcp4Server).To4()

		addr, err := nex.FindMacIpv4(pkt.CHAddr(), network.Name)
		if err != nil {
			log.WithError(err).Error("request: find mac ipv4 error")
			return nil
		}
		if addr == nil {

			log.WithFields(fields).Warn("request: no address found")
			return dhcp.ReplyPacket(pkt, dhcp.NAK, server, nil, 0, nil)

		}

		// log.Info(addr)
		_, subnetCIDR, err := net.ParseCIDR(network.Subnet4)
		if err != nil {
			log.Errorf("bad subnet: %v", err)
			return nil
		}

		if addr.Equal(rqAddr) || rqAddr.Equal(net.IPv4(0, 0, 0, 0)) {
			log.WithFields(fields).Debug("request: OK")

			nex.RenewLease(pkt.CHAddr(), network.Name)

			opts := ToOpt(network.Options)

			opts = append(opts, dhcp.Option{
				Code:  dhcp.OptionSubnetMask,
				Value: subnetCIDR.Mask,
			})

			for _, x := range network.Gateways {
				ip := net.ParseIP(x)
				opts = append(opts, dhcp.Option{
					Code:  dhcp.OptionRouter,
					Value: ip.To4(),
				})
			}

			for _, x := range network.Nameservers {
				ip := net.ParseIP(x)
				opts = append(opts, dhcp.Option{
					Code:  dhcp.OptionDomainNameServer,
					Value: ip.To4(),
				})
			}

			if network.Domain != "" {
				cn := compressedDNSName(network.Domain)
				opts = append(opts, dhcp.Option{
					Code:  dhcp.OptionDomainSearch,
					Value: cn,
				})
			}

			pkt := dhcp.ReplyPacket(
				pkt, dhcp.ACK, server, addr,
				time.Duration(network.LeaseDuration)*time.Second,
				opts)
			// add next-server option to packet
			if net.ParseIP(network.Siaddr) != nil {
				// does not support ipv6
				pkt.SetSIAddr(net.ParseIP(network.Siaddr))
			}
			return pkt

		}

		// Received dhcp request for unknown address
		fields["rqAddr"] = rqAddr
		fields["addr"] = addr
		log.WithFields(fields).Warn("request: unsolicited IP")
		return dhcp.ReplyPacket(pkt, dhcp.NAK, server, nil, 0, nil)

	case dhcp.Release:
		log.Infof("release: %s", pkt.CHAddr())

	case dhcp.Decline:
		log.Debugf("decline: %s", pkt.CHAddr())

	}

	return nil

}

func compressedDNSName(name string) []byte {
	parts := strings.Split(name, ".")

	var payload []byte

	for _, p := range parts {
		payload = append(payload, byte(len(p)))
		payload = append(payload, []byte(p)...)
	}

	payload = append(payload, byte(0))

	return payload
}

// ToOpt converts nex options into dhcp options
func ToOpt(opts []*nex.Option) []dhcp.Option {
	var result []dhcp.Option
	for _, x := range opts {
		result = append(result, dhcp.Option{
			Code:  dhcp.OptionCode(x.Number),
			Value: []byte(x.Value),
		})
	}
	return result
}
