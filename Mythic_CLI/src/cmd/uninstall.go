package cmd

import (
	"github.com/tigerMeta/tiger_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

// installCmd represents the config command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall [container name]",
	Short: "uninstall services locally and remove them from disk",
	Long:  `Run this command to uninstall a local tiger service and remove its contents from disk`,
	Run:   uninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

func uninstall(cmd *cobra.Command, args []string) {
	internal.UninstallService(args)
}
