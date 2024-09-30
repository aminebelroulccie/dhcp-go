package nex

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func withEtcd(f func(*clientv3.Client) error) error {

	c, err := EtcdClient()
	if err != nil {
		return err
	}
	defer c.Close()

	return f(c)
}

func EtcdClient() (*clientv3.Client, error) {
	err := LoadConfig()
	if err != nil {
		return nil, err
	}
	c := Current

	// var tlsc *tls.Config
	// if c.HasTls() {
	// 	capool := x509.NewCertPool()
	// 	capem, err := ioutil.ReadFile(c.CAcert)
	// 	if err != nil {
	// 		return nil, Errorf(fmt.Sprintf("error reading cacert '%s'", c.CAcert), err)
	// 	}
	// 	ok := capool.AppendCertsFromPEM(capem)
	// 	if !ok {
	// 		log.Error("ca invalid")
	// 		return nil, fmt.Errorf("ca invalid")
	// 	}

	// 	cert, err := tls.LoadX509KeyPair(
	// 		c.Cert,
	// 		c.Key,
	// 	)
	// 	if err != nil {
	// 		log.Errorf("error loading keys: %s", err)
	// 		return nil, err
	// 	}

	// 	tlsc = &tls.Config{
	// 		RootCAs:      capool,
	// 		Certificates: []tls.Certificate{cert},
	// 	}
	// }
	dialTimeout, err := time.ParseDuration(c.DialTimeout)
	if err != nil {
		return nil, err
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.Endpoints,
		DialTimeout: dialTimeout,
		// TLS: &tls.Config{
		// 	InsecureSkipVerify: c.InscureTransport,
		// },
	})

	return cli, err
}

func fetchOneKV(format string, args ...interface{}) (*mvccpb.KeyValue, error) {
	c, err := EtcdClient()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	key := fmt.Sprintf(format, args...)
	resp, err := c.Get(context.TODO(), key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, nil
	}

	return resp.Kvs[0], nil
}

func fetchKvs(format string, args ...interface{}) ([]*mvccpb.KeyValue, error) {
	c, err := EtcdClient()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	key := fmt.Sprintf(format, args...)
	resp, err := c.Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	return resp.Kvs, nil
}
