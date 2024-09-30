package nex

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/mergetb/yaml/v3"
	log "github.com/sirupsen/logrus"
	etcd "go.etcd.io/etcd/client/v3"
)

var Version string = "v1.0"
var Current *Config

var (
	AddInterfaceChan    = make(chan *string, 1)
	DeleteInterfaceChan = make(chan *string, 1)
)

var ConfigPath = flag.String("config", "/opt/vpp-agent/dev/etcd.conf", "config file location")

type Addrs struct {
	Ip4 net.IP
	Ip6 net.IP
}

/* Primary API functions ++++++++++++++++++++++++++++++++++++++++++++++++++++++
+
+ All of the primary API functions exist to modify the nex database. However,
+ they do not modify the database directly. They return transaction operations
+ that can be composed into transactions by higher level calling functions.
+ This is necessary to support database API operations with non-trivial data
+ dependencies.
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/

func FindMacNetwork(mac net.HardwareAddr, sip net.IP) (*Network, error) {

	// First look up static member
	// member := NewMacIndex(&Member{Mac: strings.ToLower(mac.String())})
	// err := ReadNew(member)
	// if err != nil {
	// 	return nil, err
	// }
	// netobj := NewNetworkObj(&Network{Name: member.Net})
	// err = Read(netobj)
	// if err == nil {
	// 	return netobj.Network, nil
	// }
	// if err != nil && !IsNotFound(err) {
	// 	return nil, err
	// }
	// If no static member found, search for dynamic members
	nets, err := GetNetworks()
	if err != nil {
		return nil, err
	}
	for _, x := range nets {
		if sip.String() == x.Siaddr {
			// log.Debugf("%s in %s?", mac, x.Name)
			// if x.MacRange == nil {
			// 	continue
			// }
			hwbegin, err := net.ParseMAC("00:00:00:00:00:00")
			if err != nil {
				log.Warnf("network '%s' has invalid mac_range begin", x)
				continue
			}
			hwend, err := net.ParseMAC("ff:ff:ff:ff:ff:ff")
			if err != nil {
				log.Warnf("network '%s' has invalid mac_range end", x)
				continue
			}

			begin := binary.BigEndian.Uint64(append([]byte{0, 0}, []byte(hwbegin)...))
			end := binary.BigEndian.Uint64(append([]byte{0, 0}, []byte(hwend)...))
			here := binary.BigEndian.Uint64(append([]byte{0, 0}, []byte(mac)...))

			// log.Debugf("lower=%d (%s)", begin, hwbegin)
			// log.Debugf("upper=%d (%s)", end, hwend)
			// log.Debugf("here =%d (%s)", here, mac)

			if begin < here && here <= end {
				return x, nil
			}
		}

	}

	return nil, nil

}

func FindMacIpv4(mac net.HardwareAddr, nets string) (net.IP, error) {

	member := NewMacIndex(&Member{Mac: mac.String(), Net: nets})
	err := Read(member)
	if err != nil {
		return nil, err
	}
	if member.Ip4 == nil {
		return nil, nil
	}

	return net.ParseIP(member.Ip4.Address), nil

}

func ResolveName(name string) ([]*Addrs, error) {

	name = strings.ToLower(name)

	log.WithFields(log.Fields{"name": name}).Info("resolving name")

	var macs []string
	err := withEtcd(func(c *etcd.Client) error {

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		resp, err := c.Get(ctx, "/member/name/"+name, etcd.WithPrefix())
		cancel()
		if err != nil {
			return err
		}

		for _, x := range resp.Kvs {
			macs = append(macs, string(x.Value))
		}

		return nil

	})
	if err != nil {

		log.WithError(err).WithFields(log.Fields{
			"name": name,
		}).Error("name query failed")

		return nil, err

	}

	members := make([]Object, len(macs))
	for i, m := range macs {
		members[i] = NewMacIndex(&Member{Mac: m})
	}

	_, err = ReadObjects(members)
	if err != nil {
		log.WithError(err).Error("failed to read members")
	}

	result := make([]*Addrs, len(members))
	for i, m := range members {

		ip4 := net.ParseIP(m.(*MacIndex).Ip4.Address)

		if ip4 != nil {
			result[i] = &Addrs{Ip4: ip4}
		}

	}

	return result, nil

}

/* types ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

type Opt4 struct {
	Number int
	Value  string
}

type Opt6 struct {
	Number int
	Value  string
}

// type EtcdConfig struct {
// 	Host   string `yaml:"host"`
// 	Port   int    `yaml:"port"`
// 	Cert   string `yaml:"cert"`
// 	Key    string `yaml:"key"`
// 	CAcert string `yaml:"cacert"`
// }
// 

type Config struct {
	Endpoints        []string `yaml:"endpoints"`
	InscureTransport bool     `yaml:"insecure-transport"`
	DialTimeout      string   `yaml:"dial-timeout"`
}

// type DhcpdConfig struct {
// 	Interface      string `yaml:"interface"`
// 	InterfaceIndex int    `yaml:"interface_index"`
// }

// type NexdConfig struct {
// 	Listen string `yaml:"listen"`
// }

// func (c EtcdConfig) HasTls() bool {
// 	return c.CAcert != "" && c.Cert != "" && c.Key != ""
// }

/* helper functions ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/
func init() {
	flag.Parse()
}



func Errorf(message string, err error) error {
	err = fmt.Errorf("%s : %s", message, err)
	log.Error(err)
	return err
}

func LoadConfig() error {

	data, err := ioutil.ReadFile(*ConfigPath)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("could not read configuration file")
	}

	err = yaml.Unmarshal(data, &Current)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("could not parse configuration file")
	}
	return nil

}

// helpers ====================================================================

func poolIndex(key string) (int, error) {
	parts := strings.Split(key, "/")
	index := parts[len(parts)-1]
	return strconv.Atoi(index)
}
