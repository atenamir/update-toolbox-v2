package cmd

import "github.com/spf13/cobra"

// ===================== newRootCmd creates the root command=======================
func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "update-toolbox",
		Short: "Secure software update and recovery toolkit",
		Long: `Update Toolbox V2 is a CLI application
for generating update packages, validating
deployments, applying updates, and recovering
from failed update operations.`,

		//--- Error handling settings ---
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newPatchCmd())
	return rootCmd
}

// ====================== Execute builds and runs the CLI =======================
func Execute() error {
	return newRootCmd().Execute()
}
