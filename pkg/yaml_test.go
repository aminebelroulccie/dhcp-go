package nex

import (
	"reflect"
	"testing"
)

func TestSpecRead(t *testing.T) {

	objs, err := ReadSpec("../tests/little/networks/mini.yml")
	if err != nil {
		t.Fatal(err)
	}

	if len(objs) != 2 {
		t.Fatal("expected 2 objects")
	}

	net, ok := objs[0].Object.(*Network)
	if !ok {
		t.Fatal("expected first object to be a Network")
	}
	expectedNet := Network{
		Name:        "mini",
		Subnet4:     "10.0.0.0/24",
		Gateways:    []string{"10.0.0.1"},
		Nameservers: []string{"10.0.0.1"},
		Dhcp4Server: "10.0.0.1",
		Domain:      "mini.net",
		Range4:      &AddressRange{Begin: "10.0.0.10", End: "10.0.0.254"},
		Siaddr:      "",
	}
	if !reflect.DeepEqual(*net, expectedNet) {
		t.Error("nets do not match")
		t.Logf("expected")
		t.Logf("%#v", expectedNet)
		t.Logf("found")
		t.Logf("%#v", *net)
	}

	list, ok := objs[1].Object.(*MemberList)
	if !ok {
		t.Fatal("expected second object to be a MemberList")
	}
	expectedList := MemberList{
		Net: "mini",
		List: []*Member{
			{Mac: "00:00:11:11:00:01", Name: "whiskey"},
			{Mac: "00:00:22:22:00:01", Name: "tango"},
			{Mac: "00:00:33:33:00:01", Name: "foxtrot"},
		},
	}
	if !reflect.DeepEqual(*list, expectedList) {
		t.Error("member lists do not match")
		t.Logf("expected")
		t.Logf("%#v", expectedList)
		t.Logf("found")
		t.Logf("%#v", *list)
	}

}
