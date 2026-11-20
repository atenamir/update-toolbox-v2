package cmd

import (
	"fmt"

	"github.com/atenamir/update-toolbox-v2/internal/patch"
	"github.com/spf13/cobra"
)

func newApplyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "apply <patch.tar.gz> <lab-install-dir>",
		Short: "Apply a patch to a temporary copy of a lab installation",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			updatedDir, backupDir, err := patch.ApplyPatchToCopy(args[0], args[1])
			if err != nil {
				if backupDir != "" {
					return fmt.Errorf("apply failed (backup: %s): %w", backupDir, err)
				}
				return err
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(),
				"Updated installation: %s\nBackup: %s\n", updatedDir, backupDir)
			return err
		},
	}
}
