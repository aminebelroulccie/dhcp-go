package nex

import (
	"reflect"
	"sort"
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
)

func TestNetworkUpdate(t *testing.T) {

	net := &Network{
		Name:        "test",
		Subnet4:     "10.47.0.0/24",
		Dhcp4Server: "10.47.0.1",
		Range4: &AddressRange{
			Begin: "10.47.0.10",
			End:   "10.47.0.254",
		},
		// MacRange: &AddressRange{
		// 	Begin: "00:00:00:00:00:01",
		// 	End:   "00:00:00:00:00:99",
		// },
		Gateways:      []string{"10.47.0.1"},
		Nameservers:   []string{"10.47.0.1"},
		Domain:        "test.net",
		Siaddr:        "10.47.0.1",
		LeaseDuration: 10000,
		Options: []*Option{{
			Number: 66,
			Value:  "10.47.0.1",
		}},
	}

	update := &NetworkUpdateRequest{
		Name:          "test",
		Domain:        &wrappers.StringValue{Value: "test.io"},
		Subnet4:       &wrappers.StringValue{Value: "10.99.0.0/24"},
		Dhcp4Server:   &wrappers.StringValue{Value: "10.99.0.0/24"},
		Siaddr:        &wrappers.StringValue{Value: "10.99.0.1"},
		LeaseDuration: &wrappers.UInt64Value{Value: 900},
		Range4: &AddressRange{
			Begin: "10.99.0.10",
			End:   "10.99.0.254",
		},
		MacRange: &AddressRange{
			Begin: "00:00:00:00:00:AA",
			End:   "00:00:00:00:00:FF",
		},
		GatewaysPresent:    []string{"10.99.0.1", "10.99.0.2", "10.99.0.1"},
		GatewaysAbsent:     []string{"10.47.0.1", "10.47.0.2"},
		NameserversPresent: []string{"10.99.0.1"},
		NameserversAbsent:  []string{"10.47.0.1"},
		OptionsPresent:     []*Option{{Number: 47, Value: "test"}},
		OptionsAbsent:      []*Option{{Number: 66}},
	}

	ApplyNetworkUpdate(net, update)

	if net.Subnet4 != update.Subnet4.Value {
		t.Logf("%v != %v", net.Subnet4, update.Subnet4.Value)
		t.Fatal("subnet4 update failed")
	}

	if net.Domain != update.Domain.Value {
		t.Logf("%v != %v", net.Domain, update.Domain.Value)
		t.Fatal("domain update failed")
	}

	if net.Dhcp4Server != update.Dhcp4Server.Value {
		t.Logf("%v != %v", net.Dhcp4Server, update.Dhcp4Server.Value)
		t.Fatal("dhcp4server update failed")
	}

	if net.Siaddr != update.Siaddr.Value {
		t.Logf("%v != %v", net.Siaddr, update.Siaddr.Value)
		t.Fatal("dhcp4server update failed")
	}

	if net.LeaseDuration != update.LeaseDuration.Value {
		t.Logf("%v != %v", net.LeaseDuration, update.LeaseDuration.Value)
		t.Fatal("lease duration update failed")
	}

	if net.Range4.Begin != update.Range4.Begin {
		t.Logf("%v != %v", net.Range4.Begin, update.Range4.Begin)
		t.Fatal("range4 begin update failed")
	}

	if net.Range4.End != update.Range4.End {
		t.Logf("%v != %v", net.Range4.End, update.Range4.End)
		t.Fatal("range4 end update failed")
	}

	// if net.MacRange.Begin != update.MacRange.Begin {
	// 	t.Logf("%v != %v", net.MacRange.Begin, update.MacRange.Begin)
	// 	t.Fatal("range4 begin update failed")
	// }

	// if net.MacRange.End != update.MacRange.End {
	// 	t.Logf("%v != %v", net.MacRange.End, update.MacRange.End)
	// 	t.Fatal("range4 end update failed")
	// }

	expected := []string{"10.99.0.1", "10.99.0.2"}
	sort.Strings(net.Gateways)
	if !reflect.DeepEqual(net.Gateways, expected) {
		t.Logf("%v != %v", net.Gateways, expected)
		t.Fatal("gateway update failed")
	}

	expected = []string{"10.99.0.1"}
	sort.Strings(net.Nameservers)
	if !reflect.DeepEqual(net.Nameservers, expected) {
		t.Logf("%v != %v", net.Nameservers, expected)
		t.Fatal("nameserver update failed")
	}

	opts := []*Option{{Number: 47, Value: "test"}}
	sort.Slice(net.Nameservers, func(i, j int) bool {
		return net.Options[i].Number < net.Options[i].Number
	})
	if !reflect.DeepEqual(net.Options, opts) {
		t.Logf("%v != %v", net.Options, opts)
		t.Fatal("option update failed")
	}

}
