package cmd

import "github.com/spf13/cobra"

// ==========rootCmd is the main command of Update Toolbox ===========
var rootCmd = &cobra.Command{
	Use:   "update-toolbox",
	Short: "Secure software update and recovery toolkit",
	Long: `Update Toolbox V2 is a CLI application
          for generating update packages, validating
          deployments, applying updates, and recovering
          from failed update operations.`,
	// --- Error handling settings ---
	SilenceErrors: true,
	SilenceUsage:  true,
}

// =========  Execute runs the root command and its subcommands ===========
func Execute() error {
	return rootCmd.Execute()
}
