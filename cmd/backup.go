package cmd

import (
	"fmt"

	"github.com/atenamir/update-toolbox-v2/internal/patch"
	"github.com/spf13/cobra"
)

func newBackupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "backup <patch.tar.gz> <install-dir> <installed-version-file>",
		Short: "Back up installed files that a patch will replace or delete",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			manifest, err := patch.ReadVerifiedPatchManifest(args[0])
			if err != nil {
				return err
			}

			if err := patch.CheckInstalledVersion(manifest, args[2]); err != nil {
				return err
			}

			backupDir, err := patch.CreateBackup(manifest, args[1])
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), "Backup created:", backupDir)
			return err
		},
	}
}
