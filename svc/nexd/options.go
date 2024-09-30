package nexd

import (
	"context"
	"fmt"
	"net"
	"strings"

	nex "gitlab.com/mergetb/tech/nex/pkg"
)

var Listen = "0.0.0.0:6000"

type NexD struct {
	nex.UnimplementedNexServer
	AddChan chan *nex.InterfaceRequest
	DelChan chan *nex.InterfaceRequest
}

func New() *NexD {
	return &NexD{
		AddChan: make(chan *nex.InterfaceRequest, 16),
		DelChan: make(chan *nex.InterfaceRequest, 16),
	}
}

var DefaultNexD = New()

func (s *NexD) AddInterface(ctx context.Context, req *nex.InterfaceRequest) (*nex.InterfaceResponse, error) {

	// found := nex.GetInterface(req.Name)
	// // if found {
	// // 	return nil, fmt.Errorf("interface with name %s already exist", req.Name)
	// // }
	// if !found {
	s.AddChan <- req
	//}
	return &nex.InterfaceResponse{}, nil
}

func (s *NexD) DeleteInterface(ctx context.Context, req *nex.InterfaceRequest) (*nex.InterfaceResponse, error) {
	// found := nex.GetInterface(req.Name)
	// if found {
	// nex.(&nex.InterfaceObj{
	// 	InterfaceRequest: req,
	// })
	//nex.DeleteInterface(req.Name)
	s.DelChan <- req
	return &nex.InterfaceResponse{}, nil
	// }
	// return nil, fmt.Errorf("the interface with name %s does not exist", req.Name)

}

/***~~~~ Networks ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

func (s *NexD) GetNetwork(
	ctx context.Context, e *nex.GetNetworkRequest,
) (*nex.GetNetworkResponse, error) {

	net := nex.NewNetworkObj(&nex.Network{Name: e.Name})
	err := nex.Read(net)
	if err != nil {
		return nil, err
	}
	return &nex.GetNetworkResponse{Net: net.Network}, nil

}

func (s *NexD) GetNetworks(
	ctx context.Context, e *nex.GetNetworksRequest,
) (*nex.GetNetworksResponse, error) {

	list, err := nex.GetNetworks()
	if err != nil {
		return nil, err
	}
	var result []string
	for _, net := range list {
		result = append(result, net.Name)
	}

	return &nex.GetNetworksResponse{Nets: result}, nil

}

func (s *NexD) AddNetwork(
	ctx context.Context, e *nex.AddNetworkRequest,
) (*nex.AddNetworkResponse, error) {
	net, _ := nex.GetNetwork(e.Network.Name)
	if net == nil {
		err := nex.AddNetwork(e.Network)
		if err != nil {
			return nil, err
		}
	}

	return &nex.AddNetworkResponse{}, nil

}

func (s *NexD) UpdateNetwork(
	ctx context.Context, e *nex.NetworkUpdateRequest,
) (*nex.NetworkUpdateResponse, error) {

	err := nex.UpdateNetwork(e)
	if err != nil {
		return nil, err
	}

	return &nex.NetworkUpdateResponse{}, nil

}

func (s *NexD) DeleteNetwork(
	ctx context.Context, e *nex.DeleteNetworkRequest,
) (*nex.DeleteNetworkResponse, error) {

	err := nex.DeleteNetwork(e.Name)
	if err != nil {
		return nil, err
	}

	return &nex.DeleteNetworkResponse{}, nil

}

/***~~~~~~~ Members ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

func (s *NexD) GetMembers(
	ctx context.Context, e *nex.GetMembersRequest,
) (*nex.GetMembersResponse, error) {

	members, err := nex.GetMembers(e.Network)
	if err != nil {
		return nil, err
	}

	return &nex.GetMembersResponse{Members: members}, nil

}

func (s *NexD) AddMembers(
	ctx context.Context, e *nex.MemberList,
) (*nex.AddMembersResponse, error) {

	net := nex.NewNetworkObj(&nex.Network{Name: e.Net})
	err := nex.Read(net)
	if err != nil {
		return nil, err
	}

	var objects []nex.Object
	for _, m := range e.List {

		if err := nex.ValidateMac(m.Mac); err != nil {
			return nil, err
		}
		m.Mac = strings.ToLower(m.Mac)

		m.Net = e.Net

		objects = append(objects, []nex.Object{
			nex.NewMacIndex(m),
			nex.NewNetIndex(m),
		}...)

		if m.Ip4 != nil {
			if net.Range4 != nil {
				return nil, fmt.Errorf("cannot assign static IP to pool member")
			}
			objects = append(objects, nex.NewIp4Index(m))
		}
		if m.Name != "" {
			m.Name = m.Name + "." + net.Domain
			m.Name = strings.ToLower(m.Name)
			objects = append(objects, nex.NewNameIndex(m))
		}

	}

	err = nex.CheckDupes(objects)
	if err != nil {
		return nil, err
	}

	var opts nex.WriteOpts
	if e.Force {
		opts = append(opts, nex.NoCheck)
	}
	err = nex.WriteObjects(objects, opts...)
	if err != nil {
		if nex.IsTxnFailed(err) {
			return nil, fmt.Errorf("some or all members already exist")
		}
		return nil, err
	}

	return &nex.AddMembersResponse{}, nil

}

func (s *NexD) UpdateMembers(
	ctx context.Context, e *nex.UpdateList,
) (*nex.UpdateMembersResponse, error) {

	net := nex.NewNetworkObj(&nex.Network{Name: e.Net})
	err := nex.Read(net)
	if err != nil {
		return nil, err
	}

	// Read the current state of the objects being updated in a single shot txn.
	var otx nex.ObjectTx
	for _, u := range e.List {

		otx.Put = append(otx.Put, nex.NewMacIndex(&nex.Member{Mac: u.Mac}))

	}
	_, err = nex.ReadObjects(otx.Put)
	if err != nil {
		return nil, err
	}

	// Update the objects in a single shot txn. The txn will fail if any of the
	// objects have been modified since reading.
	for i, object := range otx.Put {

		m := object.(*nex.MacIndex).Member
		update := e.List[i]
		if update.Name != nil {
			if m.Name != update.Name.GetValue() {
				otx.Delete = append(otx.Delete, nex.NewNameIndex(&nex.Member{Name: m.Name}))
				otx.Put = append(otx.Put, nex.NewNameIndex(m))
			}
			m.Name = update.Name.GetValue() + "." + net.Domain
			m.Name = strings.ToLower(m.Name)
		}
		if update.Ip4 != nil {
			if net.Range4 != nil {
				return nil, fmt.Errorf("cannot assign static IP to pool member")
			}
			if m.Ip4 == nil {
				otx.Delete = append(otx.Delete, nex.NewIp4Index(&nex.Member{Ip4: m.Ip4}))
				otx.Put = append(otx.Put, nex.NewIp4Index(m))
			}
			m.Ip4 = update.Ip4
		}

	}
	err = nex.RunObjectTx(otx)
	if err != nil {
		return nil, err
	}

	return &nex.UpdateMembersResponse{}, nil

}

func (s *NexD) ChangeMemberID(
	ctx context.Context, e *nex.ChangeList,
) (*nex.ChangeMemberIDResponse, error) {

	var objs []nex.Object
	members := make(map[string]*nex.Member)
	var otx nex.ObjectTx

	for _, u := range e.List {

		old := strings.ToLower(u.Old)

		m := &nex.Member{Mac: old}
		members[old] = m
		objs = append(objs, nex.NewMacIndex(m))

	}
	_, err := nex.ReadObjects(objs)
	if err != nil {
		return nil, err
	}

	for _, x := range objs {

		m := x.(*nex.MacIndex)
		if m.GetVersion() == 0 {
			return nil, fmt.Errorf("%s does not exist", m.Member.Mac)
		}

	}

	for _, m := range members {

		// make a copy of the member and collect the indicies to delete
		copy := new(nex.Member)
		*copy = *m
		otx.Delete = append(otx.Delete, nex.NewMacIndex(copy))
		otx.Delete = append(otx.Delete, nex.NewNetIndex(copy))
		if m.Name != "" {
			otx.Delete = append(otx.Delete, nex.NewNameIndex(copy))
		}

	}

	// operations to create 'new' members
	for _, x := range e.List {

		// validate old and new macs
		old := strings.ToLower(x.Old)
		_, err := net.ParseMAC(old)
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid mac: %v", old, err)
		}

		new := strings.ToLower(x.New)
		_, err = net.ParseMAC(new)
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid mac: %v", new, err)
		}

		m, ok := members[old]
		if !ok {
			continue
		}
		m.Mac = new

		otx.Put = append(otx.Put, nex.NewMacIndex(m))
		otx.Put = append(otx.Put, nex.NewNetIndex(m))
		if m.Name != "" {
			otx.Put = append(otx.Put, nex.NewNameIndex(m))
		}
		if m.Ip4 != nil {
			otx.Put = append(otx.Put, nex.NewIp4Index(m))
		}

	}

	err = nex.RunObjectTx(otx)
	if err != nil {
		return nil, err
	}

	return &nex.ChangeMemberIDResponse{}, nil

}

func (s *NexD) DeleteMembers(
	ctx context.Context, e *nex.DeleteMembersRequest,
) (*nex.DeleteMembersResponse, error) {

	err := nex.DeleteMembers(e.List)
	if err != nil {
		return nil, err
	}

	return &nex.DeleteMembersResponse{}, nil
}
