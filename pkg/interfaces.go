package nex

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
)

func GetInterfaces() ([]*InterfaceRequest, error) {
	var interfaces []*InterfaceRequest
	if err := withEtcd(func(c *etcd.Client) error {

		// get the member macs for the provided network
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err := c.Get(ctx,"/interface",etcd.WithPrefix())
        
		cancel()
		if err != nil {

			return err
		}
		for _, kv := range resp.Kvs {
			iface := &InterfaceRequest{}
			if err := json.Unmarshal(kv.Value, &iface); err != nil {
				continue
			}
			interfaces = append(interfaces, iface)
		}

		return nil

	}); err != nil {
		return nil, err
	}

	return interfaces, nil
}

func GetInterface(name string) bool {
	iface := &InterfaceRequest{}
	err := withEtcd(func(c *etcd.Client) error {

		// get the member macs for the provided network
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		resp, err := c.Get(ctx, "/interface/"+name)
		cancel()
		if err != nil {

			return err
		}

		if resp.Count > 0 {
			if err := json.Unmarshal(resp.Kvs[0].Value, &iface); err != nil {
				fmt.Println("hehehhe")
				return err
			}
		} else {
			return fmt.Errorf("interface with name %s not found", name)
		}

		return nil

	})
	return err == nil
}

func DeleteInterface(name string) error {

	var objects []Object
	objects = append(objects, NewInterfaceRequestObj(&InterfaceRequest{Name: name}))
	_, err := ReadObjects(objects)
	if err != nil {
		return err
	}

	return DeleteObjects(objects)

}
