package cmd

import (
	"fmt"

	"github.com/atenamir/update-toolbox-v2/internal/patch"
	"github.com/spf13/cobra"
)

func newVerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify <patch.tar.gz>",
		Short: "Verify an update patch archive",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := patch.VerifyPatchArchive(args[0]); err != nil {
				return err
			}
         	fmt.Fprintln(cmd.OutOrStdout(), "Patch verification passed")
			return nil
		},
	}
}