package cmd

import (
	"fmt"
	"github.com/tigerMeta/tiger_CLI/cmd/internal"
	"github.com/spf13/cobra"
	"os"
)

// installCmd represents the config command
var installtigerSyncFolderCmd = &cobra.Command{
	Use:   "folder {path} ",
	Short: "install tiger_sync from local folder",
	Long:  `Run this command to install tiger_sync from a locally cloned folder`,
	Run:   installtigerSyncFolder,
	Args:  cobra.ExactArgs(1),
}

func init() {
	installtigerSyncCmd.AddCommand(installtigerSyncFolderCmd)
}

func installtigerSyncFolder(cmd *cobra.Command, args []string) {
	if err := internal.InstalltigerSyncFolder(args[0]); err != nil {
		fmt.Printf("[-] Failed to install service: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("[+] Successfully installed service!\n")
	}
}
