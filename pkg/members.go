package nex

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
)

func (m *Member) Clone() *Member {

	c := &Member{
		Mac:  m.Mac,
		Name: m.Name,
		Net:  m.Net,
	}
	if m.Ip4 != nil {
		c.Ip4 = &Lease{}
		*c.Ip4 = *m.Ip4
	}
	if m.Ip6 != nil {
		c.Ip6 = &Lease{}
		*c.Ip6 = *m.Ip6
	}

	return c

}

func GetMembers(net string) ([]*Member, error) {

	var members []*Member

	err := withEtcd(func(c *etcd.Client) error {

		// get the member macs for the provided network
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		resp, err := c.Get(ctx, "/member/net/"+net, etcd.WithPrefix())
		cancel()
		if err != nil {
			return err
		}

		var macs []string
		for _, kv := range resp.Kvs {
			macs = append(macs, string(kv.Value))
		}

		// get the members from the provided macs
		members, err = FetchMembers(macs, c)
		if err != nil {
			return err
		}

		return nil

	})
	if err != nil {
		return nil, err
	}

	return members, nil

}

func FetchMembers(macs []string, c *etcd.Client) ([]*Member, error) {

	var members []*Member

	var ops []etcd.Op
	for _, mac := range macs {
		ops = append(ops, etcd.OpGet("/member/mac/"+mac))
	}

	kvc := etcd.NewKV(c)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	txr, err := kvc.Txn(ctx).Then(ops...).Commit()
	cancel()
	if err != nil {
		return nil, err
	}
	if !txr.Succeeded {
		return nil, TxnFailed("")
	}

	for _, r := range txr.Responses {
		rr := r.GetResponseRange()
		if rr == nil {
			continue
		}

		for _, kv := range rr.Kvs {

			member := &Member{}
			err := json.Unmarshal(kv.Value, &member)
			if err != nil {
				return nil, err
			}
			members = append(members, member)

		}
	}

	return members, nil

}

func DeleteMembers(macs []string) error {

	var objects []Object
	for _, mac := range macs {

		if err := ValidateMac(mac); err != nil {
			return err
		}

		objects = append(objects, NewMacIndex(&Member{Mac: mac}))

	}
	_, err := ReadObjects(objects)
	if err != nil {
		return err
	}

	objects = DeleteMemberObjects(objects)

	return DeleteObjects(objects)

}

func DeleteMemberObjects(objects []Object) []Object {

	for _, o := range objects {

		x := o.(*MacIndex).Member

		objects = append(objects, NewNetIndex(x))

		if x.Ip4 != nil {
			objects = append(objects, NewIp4Index(x))
		}
		if x.Name != "" {
			objects = append(objects, NewNameIndex(x))
		}

	}

	return objects

}

func ValidateMac(mac string) error {
	if _, err := net.ParseMAC(mac); err != nil {
		return fmt.Errorf("invalid MAC: %s", mac)
	}
	return nil
}
