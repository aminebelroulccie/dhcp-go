package options

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/spf13/cobra"

	"gitlab.com/mergetb/tech/nex/pkg"
)

func MemberCmds(root, get, set, add, delete *cobra.Command) {

	getMembers := &cobra.Command{
		Use:   "members [network-name]",
		Short: "Get a networks members",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			getMembers(args[0])
		},
	}
	get.AddCommand(getMembers)

	setMember := &cobra.Command{
		Use:   "member",
		Short: "Set a member property",
	}
	set.AddCommand(setMember)

	// Set Ip4 ~~~

	setMemberIp4 := &cobra.Command{
		Use:   "ip4 [network] [mac] [ip4]",
		Short: "Set a members IPv4 address",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			setMemberIp4(args[0], args[1], args[2])
		},
	}
	setMember.AddCommand(setMemberIp4)

	// Set Name ~~~

	setMemberName := &cobra.Command{
		Use:   "name [network] [mac] [name]",
		Short: "Set a members name",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			setMemberName(args[0], args[1], args[2])
		},
	}
	setMember.AddCommand(setMemberName)

	// Change ~~~

	change := &cobra.Command{
		Use:   "change [network] [old-mac] [new-mac]",
		Short: "Change a members mac",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			change(args[0], args[1], args[2])
		},
	}
	root.AddCommand(change)

	// Add ~~~

	addMember := &cobra.Command{
		Use:   "member [network] [mac]",
		Short: "Add a member",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			addMember(args[0], args[1])
		},
	}
	add.AddCommand(addMember)

	// Delete ~~~

	deleteMember := &cobra.Command{
		Use:   "member [network] [mac]",
		Short: "Delete a member",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			deleteMember(args[0], args[1])
		},
	}
	delete.AddCommand(deleteMember)

}

func setMemberIp4(network, mac, ip4 string) {

	withClient(func(cli nex.NexClient) error {

		_, err := cli.UpdateMembers(ctx, &nex.UpdateList{
			Net: network,
			List: []*nex.MemberUpdate{
				{
					Mac: mac,
					Ip4: &nex.Lease{Address: ip4},
				},
			},
		})
		if err != nil {
			grpcFatal(err)
		}

		return nil

	})
}

func setMemberName(network, mac, name string) {
	withClient(func(cli nex.NexClient) error {

		_, err := cli.UpdateMembers(ctx, &nex.UpdateList{
			Net: network,
			List: []*nex.MemberUpdate{
				{
					Mac:  mac,
					Name: &wrappers.StringValue{Value: name},
				},
			},
		})
		if err != nil {
			grpcFatal(err)
		}

		return nil

	})
}

func change(network, old, new string) {
	withClient(func(cli nex.NexClient) error {

		_, err := cli.ChangeMemberID(ctx, &nex.ChangeList{
			Net: network,
			List: []*nex.IDChange{{
				Old: old,
				New: new,
			}},
		})
		if err != nil {
			grpcFatal(err)
		}

		return nil

	})
}

func getMembers(net string) {
	withClient(func(cli nex.NexClient) error {

		resp, err := cli.GetMembers(ctx, &nex.GetMembersRequest{Network: net})
		if err != nil {
			grpcFatal(err)
		}

		fmt.Fprint(tw, "mac\tname\tip4\tclient\n")
		for _, m := range resp.Members {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
				m.Mac,
				m.Name,
				showLease(m.Ip4),
				m.ClientName,
			)
		}
		tw.Flush()
		return nil

	})
}

func addMember(net, mac string) {
	withClient(func(cli nex.NexClient) error {

		_, err := cli.AddMembers(ctx, &nex.MemberList{
			Net: net,
			List: []*nex.Member{
				{Mac: mac},
			},
		})
		if err != nil {
			grpcFatal(err)
		}

		return nil

	})
}

func deleteMember(net, mac string) {
	withClient(func(cli nex.NexClient) error {

		_, err := cli.DeleteMembers(ctx, &nex.DeleteMembersRequest{
			Network: net,
			List:    []string{mac},
		})
		if err != nil {
			grpcFatal(err)
		}

		return nil

	})
}

func showLease(lease *nex.Lease) string {

	if lease == nil {
		return ""
	}

	expires := ""
	if lease.Expires != nil {
		t := time.Until(time.Unix(lease.Expires.Seconds, int64(lease.Expires.Nanos)))
		expires = fmt.Sprintf(" (%d:%d:%d)",
			int(t.Hours()), int(t.Minutes())%60, int(t.Seconds())%60)
	}

	return fmt.Sprintf("%s%s", lease.Address, expires)

}
