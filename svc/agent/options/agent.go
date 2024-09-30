package options

import (
	"sync"

	dhcpd "gitlab.com/mergetb/tech/nex/svc/dhcp-server"
	"gitlab.com/mergetb/tech/nex/svc/nexd"
)

type Agent struct {
	services []ServiceApi
}

func New() *Agent {
	return &Agent{
		services: []ServiceApi{
			nexd.DefaultNexD,
			dhcpd.DefaultDhcpD,
		},
	}
}

func (d *Agent) Run() {
	var wg sync.WaitGroup
	wg.Add(len(d.services))
	for _, svc := range d.services {
		go func() {
			defer wg.Done()
			svc.Run()
		}()
	}
	wg.Wait()
}
