package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"gitlab.com/mergetb/tech/rtnl"
)

var table = map[string]net.IP{
	"tango.basic":   net.IPv4(10, 0, 0, 10),
	"foxtrot.basic": net.IPv4(10, 0, 0, 11),
}

func dhcp(host string) error {

	// bring up interface

	rtx, err := rtnl.OpenDefaultContext()
	if err != nil {
		return err
	}

	eth1, err := rtnl.GetLink(rtx, "eth1")
	if err != nil {
		return err
	}

	err = eth1.Up(rtx)
	if err != nil {
		return err
	}

	// do dhcp

	out, err := exec.Command("dhclient", "-r", "eth1").CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(out))
	}

	out, err = exec.Command("dhclient", "eth1").CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(out))
	}

	addrs, err := eth1.Addrs(rtx)
	if err != nil {
		return err
	}

	log.Printf("found %d addrs for eth1", len(addrs))

	for _, addr := range addrs {
		log.Printf("%s == %s", addr.Info.Address.IP.String(), table[host].String())
		if addr.Info.Address.IP.String() == table[host].String() {
			return nil
		}
	}

	return fmt.Errorf("expected dns address not found for %s. found %v", host, addrs)

}

func dns(host string) error {

	for name, ip := range table {

		dnsIPs, err := net.LookupIP(name)
		if err != nil {
			return err
		}

		found := false
		for _, x := range dnsIPs {
			if x.String() == ip.String() {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("name not found: %s", name)
		}

	}

	return nil

}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("usage: test <host>")
	}

	err := dhcp(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	err = dns(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

}
