package cmd

import (
	"fmt"
	"os"

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
			manifestPath := args[0]

			// 1. Load the release manifest.
			manifests, err := patch.LoadManifest(manifestPath)
			if err != nil {
				return err
			}

			// 2. Validate the manifest.
			if err := patch.ValidateManifest(manifests); err != nil {
				return err
			}

			// 3. Get source and target versions.
			fromTag := manifests[0].Version
			toTag := manifests[1].Version

			// 4. Analyze Git changes.
			changedFiles, err := patch.GetChangedFiles(
				repoPath,
				fromTag,
				toTag,
			)
			if err != nil {
				return err
			}

			// 5. Build the patch manifest.
			patchManifest := patch.BuildPatchManifest(
				fromTag,
				toTag,
				changedFiles,
			)

			// 6. Create the TAR.GZ archive.
			if _, err := os.Stat(outputPath); err == nil {
				return fmt.Errorf("output file already exists: %s", outputPath)
			} else if !os.IsNotExist(err) {
				return err
			}

			if err := patch.CreatePatchArchive(
				patchManifest,
				repoPath,
				outputPath,
			); err != nil {
				return err
			}

			// 7. Print success.
			fmt.Fprintln(cmd.OutOrStdout(), "Patch created:", outputPath)
			fmt.Fprintln(cmd.OutOrStdout(), "Changed files:", len(changedFiles))

			return nil
		},
	}

	createCmd.Flags().StringVar(
		&repoPath,
		"repo",
		"",
		"Path to the Git repository",
	)

	createCmd.Flags().StringVar(
		&outputPath,
		"output",
		"",
		"Output path for the patch archive",
	)

	_ = createCmd.MarkFlagRequired("repo")
	_ = createCmd.MarkFlagRequired("output")

	return createCmd
}
