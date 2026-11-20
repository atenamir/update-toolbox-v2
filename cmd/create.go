package cmd

import (
	"fmt"

	"github.com/atenamir/update-toolbox-v2/internal/patch"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	var repoPath string
	var outputPath string

	createCmd := &cobra.Command{
		Use:   "create <package.json>",
		Short: "Create an update patch archive",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			manifest, err := patch.CreatePatch(args[0], repoPath, outputPath)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Patch created:", outputPath)
			fmt.Fprintln(cmd.OutOrStdout(), "Changed files:", len(manifest.Files))
			return nil
		},
	}

	createCmd.Flags().StringVar(&repoPath, "repo", "", "Path to the Git repository")
	createCmd.Flags().StringVar(&outputPath, "output", "", "Output path for the patch archive")

	_ = createCmd.MarkFlagRequired("repo")
	_ = createCmd.MarkFlagRequired("output")

	return createCmd
}
