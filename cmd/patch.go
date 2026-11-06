package cmd

import "github.com/spf13/cobra"

// ================ patchCmd groups commands for managing update patches =================
func newPatchCmd() *cobra.Command {
	patchCmd := &cobra.Command{
		Use:   "patch",
		Short: "manage update patches",
		Long: `Manage update packages used by Update Toolbox.
Available operations include creating and verifying update packages.`,
	}
	patchCmd.AddCommand(newCreateCmd())

	return patchCmd
}
