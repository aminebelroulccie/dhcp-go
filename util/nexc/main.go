package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	nex "gitlab.com/mergetb/tech/nex/pkg"
	"gitlab.com/mergetb/tech/nex/util/nexc/options"
)



func main() {
	log.SetFlags(0)

	root := &cobra.Command{
		Use:   "nex",
		Short: "Nex dhcp/dns client",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}
	root.PersistentFlags().StringVarP(
		&options.Server, "server", "s", "localhost:6000", "nexd server to connect to")

	rootCmds(root)

	root.Execute()
}

func rootCmds(root *cobra.Command) {

	version := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			log.Print(nex.Version)
		},
	}
	root.AddCommand(version)

	get := &cobra.Command{
		Use:   "get",
		Short: "Get something",
	}
	root.AddCommand(get)

	set := &cobra.Command{
		Use:   "set",
		Short: "Set something",
	}
	root.AddCommand(set)

	add := &cobra.Command{
		Use:   "add",
		Short: "Add something",
	}
	root.AddCommand(add)

	delete := &cobra.Command{
		Use:   "delete",
		Short: "Delete something",
	}
	root.AddCommand(delete)

	autocomplete := &cobra.Command{
		Use:   "autocomplete",
		Short: "Generates bash completion scripts",
		Run: func(cmd *cobra.Command, args []string) {
			root.GenBashCompletion(os.Stdout)
		},
	}
	root.AddCommand(autocomplete)

	options.ApplyCmd(root)
	options.MemberCmds(root, get, set, add, delete)
	options.NetworkCmds(get, set, add, delete)

}
