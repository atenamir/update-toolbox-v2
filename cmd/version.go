package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// =============== Version identifies the current build ======================
var Version = "dev"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show current version of Update Toolbox",
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "Update Toolbox V2")
			fmt.Fprintln(cmd.OutOrStdout(), "Version:", Version)
		},
	}
}
