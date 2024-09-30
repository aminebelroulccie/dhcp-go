package nex

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
)

func AddNetwork(n *Network) error {

	ApplyDefaults(n)

	var objs []Object
	if n.Range4 != nil {
		p := &Pool{Net: n.Name}
		p.Size = n.Range4.Size()
		objs = append(objs, NewPoolObj(p))
		for _, addr := range n.Excluded {
			ipAddr := net.ParseIP(addr)
			if ipAddr != nil {
				id := n.Range4.Offset(ipAddr)
				p.Values = append(p.Values, id)
			}
		}
		fmt.Println(p.Values)
	}

	objs = append(objs, NewNetworkObj(n))
	err := WriteObjects(objs)
	if err != nil {
		return err
	}

	return nil

}

func GetNetwork(name string) (*Network, error) {

	//var nets []*Network
	var net *Network
	err := withEtcd(func(c *etcd.Client) error {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		resp, err := c.Get(ctx, "/net/"+name)
		cancel()
		if err != nil {
			return err
		}

		//for _, kv := range resp.Kvs {
		if resp.Count > 0 {
			net = &Network{}
			err = json.Unmarshal(resp.Kvs[0].Value, &net)
			if err != nil {
				return err
			}
			//nets = append(nets, net)
			//}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return net, nil
}

func GetNetworks() ([]*Network, error) {

	var nets []*Network

	err := withEtcd(func(c *etcd.Client) error {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		resp, err := c.Get(ctx, "/net", etcd.WithPrefix())
		cancel()
		if err != nil {
			return err
		}

		for _, kv := range resp.Kvs {
			net := &Network{}
			err := json.Unmarshal(kv.Value, &net)
			if err != nil {
				return err
			}
			nets = append(nets, net)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return nets, nil
}

func UpdateNetwork(r *NetworkUpdateRequest) error {

	net := NewNetworkObj(&Network{Name: r.Name})
	fmt.Println(net)
	err := Read(net)
	if err != nil {
		return err
	}

	ApplyNetworkUpdate(net.Network, r)

	err = Write(net)
	if err != nil {
		return err
	}

	return nil

}

func ApplyNetworkUpdate(net *Network, update *NetworkUpdateRequest) {

	if update.Subnet4 != nil {
		net.Subnet4 = update.Subnet4.Value
	}

	if update.Domain != nil {
		net.Domain = update.Domain.Value
	}

	if update.Dhcp4Server != nil {
		net.Dhcp4Server = update.Dhcp4Server.Value
	}

	if update.Siaddr != nil {
		net.Siaddr = update.Siaddr.Value
	}

	if update.LeaseDuration != nil {
		net.LeaseDuration = update.LeaseDuration.Value
	}

	if update.Range4 != nil {
		net.Range4 = &AddressRange{
			Begin: update.Range4.Begin,
			End:   update.Range4.End,
		}
	}

	// if update.MacRange != nil {
	// 	net.MacRange = &AddressRange{
	// 		Begin: update.MacRange.Begin,
	// 		End:   update.MacRange.End,
	// 	}
	// }

	net.Gateways = updateSet(
		net.Gateways,
		update.GatewaysPresent,
		update.GatewaysAbsent,
	)

	net.Nameservers = updateSet(
		net.Nameservers,
		update.NameserversPresent,
		update.NameserversAbsent,
	)

	net.Options = updateOptionSet(
		net.Options,
		update.OptionsPresent,
		update.OptionsAbsent,
	)

}

func updateSet(original, present, absent []string) []string {

	dedup := make(map[string]bool)

	for _, x := range original {
		dedup[x] = true
	}

	for _, x := range present {
		dedup[x] = true
	}

	for _, x := range absent {
		delete(dedup, x)
	}

	result := make([]string, 0, len(dedup))
	for k, _ := range dedup {
		result = append(result, k)
	}

	return result

}

func updateOptionSet(original, present, absent []*Option) []*Option {

	dedup := make(map[int32]*Option)

	for _, x := range original {
		dedup[x.Number] = x
	}

	for _, x := range present {
		dedup[x.Number] = x
	}

	for _, x := range absent {
		delete(dedup, x.Number)
	}

	result := make([]*Option, 0, len(dedup))
	for _, v := range dedup {
		result = append(result, v)
	}

	return result

}

func DeleteNetwork(name string) error {

	// Get the networks members
	members, err := GetMembers(name)
	if err != nil {
		return err
	}

	// Gather all the member index objects
	var objects []Object
	for _, m := range members {
		objects = append(objects, NewMacIndex(m))
	}
	objects = DeleteMemberObjects(objects)

	// Get the pool if there is one
	objects = append(objects, NewPoolObj(&Pool{Net: name}))

	// Add the network object itself to the list
	objects = append(objects, NewNetworkObj(&Network{Name: name}))

	// Clear out everything in one txn
	return DeleteObjects(objects)

}

func ApplyDefaults(n *Network) {

	if n.LeaseDuration == 0 {
		n.LeaseDuration = 4 * 60 * 60 //4 hour default
	}

}
