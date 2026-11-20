package cmd

import (
	"fmt"

	"github.com/atenamir/update-toolbox-v2/internal/patch"
	"github.com/spf13/cobra"
)

func newPreflightCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "preflight <patch.tar.gz> <installed-version-file>",
		Short: "Check that a patch matches the installed version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			manifest, err := patch.ReadVerifiedPatchManifest(args[0])
			if err != nil {
				return err
			}

			if err := patch.CheckInstalledVersion(manifest, args[1]); err != nil {
				return err
			}

			_, err = fmt.Fprintf(
				cmd.OutOrStdout(),
				"Installed version matches patch source: %s\n",
				manifest.FromVersion,
			)
			return err
		},
	}
}
