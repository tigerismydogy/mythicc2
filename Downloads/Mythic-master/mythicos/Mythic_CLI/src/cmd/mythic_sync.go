package cmd

import (
	"github.com/spf13/cobra"
)

// installCmd represents the config command
var tigerSyncCmd = &cobra.Command{
	Use:   "tiger_sync",
	Short: "Install/Uninstall tiger_sync",
	Long:  `Run this command's subcommands to install/uninstall tiger_sync `,
	Run:   tigerSync,
}

func init() {
	rootCmd.AddCommand(tigerSyncCmd)
}

func tigerSync(cmd *cobra.Command, args []string) {
	cmd.Help()
}
