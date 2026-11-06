package cmd

import (
	"fmt"
	"github.com/atenamir/update-toolbox-v2/internal/patch"
	"github.com/spf13/cobra"
)

// ============== newCreateCmd creates the patch create command ==================
func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create a new update patch",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			manifests, err := patch.LoadManifest(path)
			if err != nil {
				return err
			}
			if err := patch.ValidateManifest(manifests); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Manifest loaded successfully")
			fmt.Fprintln(cmd.OutOrStdout(), "Release count:", len(manifests))

			return nil
		},
	}
}
