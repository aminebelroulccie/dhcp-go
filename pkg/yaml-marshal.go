package nex

/* These functions implement the YAML marshaler interface for protobuf
 * generated types. This is necessary beacause tthere is a lot of noise
 * inside the protobuf structs as public members that come along for the
 * yaml marshalling ride.  */

func (x *Member) MarshalYAML() (interface{}, error) {
	return map[string]interface{}{
		"name": x.Name,
		"ip4":  x.Ip4,
		"ip6":  x.Ip6,
		"net":  x.Net,
		"mac":  x.Mac,
	}, nil
}

func (x *AddressRange) MarshalYAML() (interface{}, error) {
	return map[string]interface{}{
		"begin": x.Begin,
		"end":   x.End,
	}, nil
}

func (x *Option) MarshalYAML() (interface{}, error) {
	return map[string]interface{}{
		"number": x.Number,
		"value":  x.Value,
	}, nil
}

func (x *Network) MarshalYAML() (interface{}, error) {
	return map[string]interface{}{
		"name":        x.Name,
		"subnet4":     x.Subnet4,
		"subnet6":     x.Subnet6,
		"dhcp4server": x.Dhcp4Server,
		"dhcp6server": x.Dhcp6Server,
		"range4":      x.Range4,
		"range6":      x.Range6,
		"gateways":    x.Gateways,
		"nameservers": x.Nameservers,
		"options":     x.Options,
		"domain":      x.Domain,
		"siaddr":      x.Siaddr,
	}, nil
}
