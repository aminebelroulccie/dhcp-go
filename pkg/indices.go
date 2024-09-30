package nex

import (
	"context"
	"encoding/json"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
)

func FetchIp4IndexMembers() ([]*Member, error) {

	var members []*Member

	err := withEtcd(func(c *etcd.Client) error {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		resp, err := c.Get(ctx, "/member/ip4/", etcd.WithPrefix())
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

func FetchMacIndexMembers() ([]*Member, error) {

	var members []*Member

	err := withEtcd(func(c *etcd.Client) error {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		resp, err := c.Get(ctx, "/member/mac", etcd.WithPrefix())
		cancel()
		if err != nil {
			return err
		}

		for _, kv := range resp.Kvs {
			member := &Member{}
			err := json.Unmarshal(kv.Value, &member)
			if err != nil {
				return err
			}
			members = append(members, member)
		}

		return nil

	})
	if err != nil {
		return nil, err
	}

	return members, nil

}
