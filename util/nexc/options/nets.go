package options

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/mergetb/yaml/v3"
	"github.com/spf13/cobra"
	"gitlab.com/mergetb/tech/nex/pkg"
)

func NetworkCmds(get, set, add, delete *cobra.Command) {

	getNetworks := &cobra.Command{
		Use:   "networks",
		Short: "Get network list",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			getNetworks()
		},
	}
	get.AddCommand(getNetworks)

	getNetwork := &cobra.Command{
		Use:   "network",
		Short: "Get network list",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			getNetwork(args[0])
		},
	}
	get.AddCommand(getNetwork)

	setnet := &cobra.Command{
		Use:   "network",
		Short: "Set network properties",
	}
	set.AddCommand(setnet)

	deleteNetwork := &cobra.Command{
		Use:   "network [name]",
		Short: "Delete a network and all it's members",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			deleteNetwork(args[0])
		},
	}
	delete.AddCommand(deleteNetwork)

	// updates

	setSubnet4 := &cobra.Command{
		Use:   "subnet4 <network> <subnet>",
		Short: "Set network subnet",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			setNetworkProps(&nex.NetworkUpdateRequest{
				Name:    args[0],
				Subnet4: &wrappers.StringValue{Value: args[1]},
			})
		},
	}
	setnet.AddCommand(setSubnet4)

	setDhcp4Server := &cobra.Command{
		Use:   "dhcp4server <network> <host>",
		Short: "Set dhcp4server",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			setNetworkProps(&nex.NetworkUpdateRequest{
				Name:        args[0],
				Dhcp4Server: &wrappers.StringValue{Value: args[1]},
			})
		},
	}
	setnet.AddCommand(setDhcp4Server)

	setRange4 := &cobra.Command{
		Use:   "range4 <network> <begin> <end>",
		Short: "Set ip4 address range",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			setNetworkProps(&nex.NetworkUpdateRequest{
				Name: args[0],
				Range4: &nex.AddressRange{
					Begin: args[1],
					End:   args[2],
				},
			})
		},
	}
	setnet.AddCommand(setRange4)

	var (
		absent bool
	)
	setgw := &cobra.Command{
		Use:   "gateway <network> [gateway...]",
		Short: "Set gateways",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			if absent {
				setNetworkProps(&nex.NetworkUpdateRequest{
					Name:           args[0],
					GatewaysAbsent: args[1:],
				})
			} else {
				setNetworkProps(&nex.NetworkUpdateRequest{
					Name:            args[0],
					GatewaysPresent: args[1:],
				})
			}

		},
	}
	setgw.Flags().BoolVarP(&absent, "absent", "a", false, "remove gateways")
	setnet.AddCommand(setgw)

	setns := &cobra.Command{
		Use:   "nameserver <network> [nameserver...]",
		Short: "Set nameservers",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			if absent {
				setNetworkProps(&nex.NetworkUpdateRequest{
					Name:              args[0],
					NameserversAbsent: args[1:],
				})
			} else {
				setNetworkProps(&nex.NetworkUpdateRequest{
					Name:               args[0],
					NameserversPresent: args[1:],
				})
			}

		},
	}
	setns.Flags().BoolVarP(&absent, "absent", "a", false, "remove nameservers")
	setnet.AddCommand(setns)

	setopt := &cobra.Command{
		Use:   "option <network> <number> [value]",
		Short: "Set an option",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			value := ""
			if len(args) > 2 {
				value = args[2]
			}

			number, err := strconv.Atoi(args[1])
			if err != nil || number <= 0 {
				log.Fatal("number must be a positive iteger")
			}

			if absent {

				setNetworkProps(&nex.NetworkUpdateRequest{
					Name: args[0],
					OptionsAbsent: []*nex.Option{{
						Number: int32(number),
					}},
				})

			} else {

				setNetworkProps(&nex.NetworkUpdateRequest{
					Name: args[0],
					OptionsPresent: []*nex.Option{{
						Number: int32(number),
						Value:  value,
					}},
				})

			}

		},
	}
	setopt.Flags().BoolVarP(&absent, "absent", "a", false, "remove nameservers")
	setnet.AddCommand(setopt)

	setDomain := &cobra.Command{
		Use:   "domain <network> <domain>",
		Short: "Set DNS domain",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			setNetworkProps(&nex.NetworkUpdateRequest{
				Name:   args[0],
				Domain: &wrappers.StringValue{Value: args[1]},
			})
		},
	}
	setnet.AddCommand(setDomain)

	setMacRange := &cobra.Command{
		Use:   "macrange <network> <begin> <end>",
		Short: "Set mac address range",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			setNetworkProps(&nex.NetworkUpdateRequest{
				Name: args[0],
				MacRange: &nex.AddressRange{
					Begin: args[1],
					End:   args[2],
				},
			})
		},
	}
	setnet.AddCommand(setMacRange)

	setSiaddr := &cobra.Command{
		Use:   "siaddr <network> <addr>",
		Short: "Set DHCP siaddr",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			setNetworkProps(&nex.NetworkUpdateRequest{
				Name:   args[0],
				Siaddr: &wrappers.StringValue{Value: args[1]},
			})
		},
	}
	setnet.AddCommand(setSiaddr)

	setLeaseDuration := &cobra.Command{
		Use:   "duration <network> <seconds>",
		Short: "Set DHCP lease duration",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			value, err := strconv.Atoi(args[1])
			if err != nil || value <= 0 {
				log.Fatal("duration must be positive integer")
			}

			setNetworkProps(&nex.NetworkUpdateRequest{
				Name:          args[0],
				LeaseDuration: &wrappers.UInt64Value{Value: uint64(value)},
			})
		},
	}
	setnet.AddCommand(setLeaseDuration)

}

func setNetworkProps(u *nex.NetworkUpdateRequest) {

	withClient(func(cli nex.NexClient) error {

		_, err := cli.UpdateNetwork(ctx, u)
		if err != nil {
			grpcFatal(err)
		}

		return nil

	})

}

func getNetworks() {
	withClient(func(cli nex.NexClient) error {

		resp, err := cli.GetNetworks(ctx, &nex.GetNetworksRequest{})
		if err != nil {
			grpcFatal(err)
		}

		for _, n := range resp.Nets {
			log.Println(n)
		}

		return nil

	})
}

func getNetwork(name string) {
	withClient(func(cli nex.NexClient) error {

		resp, err := cli.GetNetwork(ctx, &nex.GetNetworkRequest{
			Name: name,
		})
		if err != nil {
			grpcFatal(err)
		}

		fmt.Fprintf(tw, "name:\t%s\n", resp.Net.Name)
		fmt.Fprintf(tw, "subnet4:\t%s\n", resp.Net.Subnet4)
		fmt.Fprintf(tw, "gateways:\t%s\n", strings.Join(resp.Net.Gateways, " "))
		fmt.Fprintf(tw, "nameservers:\t%s\n", strings.Join(resp.Net.Nameservers, " "))
		fmt.Fprintf(tw, "dhcp4server:\t%s\n", resp.Net.Dhcp4Server)
		fmt.Fprintf(tw, "domain:\t%s\n", resp.Net.Domain)
		if resp.Net.Siaddr != "" {
			fmt.Fprintf(tw, "siaddr:\t%s\n", resp.Net.Siaddr)
		}
		if resp.Net.Range4 != nil {
			fmt.Fprintf(tw, "range4:\t%s-%s\n", resp.Net.Range4.Begin, resp.Net.Range4.End)
		}
		// if resp.Net.MacRange != nil {
		// 	fmt.Fprintf(tw, "mac_range:\t%s-%s\n", resp.Net.MacRange.Begin, resp.Net.MacRange.End)
		// }
		fmt.Fprintf(tw, "lease_duration:\t%ds\n", resp.Net.LeaseDuration)
		tw.Flush()
		if len(resp.Net.Options) > 0 {
			fmt.Printf("options:\t\n")
			for _, opt := range resp.Net.Options {
				fmt.Fprintf(tw, "  %d\t%s\n", opt.Number, opt.Value)
			}
			tw.Flush()
		}

		return nil

	})
}

func deleteNetwork(name string) {
	withClient(func(cli nex.NexClient) error {

		_, err := cli.DeleteNetwork(ctx, &nex.DeleteNetworkRequest{Name: name})
		if err != nil {
			grpcFatal(err)
		}

		return nil

	})
}

func loadSpec(file string) (*nex.Network, error) {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	net := &nex.Network{}
	err = yaml.Unmarshal(data, net)
	if err != nil {
		return nil, err
	}

	return net, nil
}
