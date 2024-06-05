package cmd

import (
	"github.com/tigerMeta/tiger_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

// installCmd represents the config command
var uninstalltigerSyncCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall tiger_sync and remove it from disk",
	Long:  `Run this command to uninstall tiger_sync and remove it from disk`,
	Run:   uninstalltigerSyncGitHub,
}

func init() {
	tigerSyncCmd.AddCommand(uninstalltigerSyncCmd)
}

func uninstalltigerSyncGitHub(cmd *cobra.Command, args []string) {
	internal.UninstalltigerSync()
}
