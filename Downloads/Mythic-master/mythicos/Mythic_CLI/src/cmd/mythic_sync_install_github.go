package cmd

import (
	"fmt"
	"github.com/tigerMeta/tiger_CLI/cmd/internal"
	"github.com/spf13/cobra"
	"os"
)

// installCmd represents the config command
var installtigerSyncGitHubCmd = &cobra.Command{
	Use:   "github [url] [branch] [-f]",
	Short: "install services from GitHub or other Git-based repositories",
	Long:  `Run this command to install services from Git-based repositories by doing a git clone`,
	Run:   installtigerSyncGitHub,
	Args:  cobra.RangeArgs(0, 2),
}

func init() {
	installtigerSyncCmd.AddCommand(installtigerSyncGitHubCmd)
	installtigerSyncGitHubCmd.Flags().BoolVarP(
		&force,
		"force",
		"f",
		false,
		`Force installing from GitHub and don't prompt to overwrite files if an older version is already installed`,
	)
	installtigerSyncGitHubCmd.Flags().StringVarP(
		&branch,
		"branch",
		"b",
		"",
		`Install a specific branch from GitHub instead of the main/master branch`,
	)
}

func installtigerSyncGitHub(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		if err := internal.InstalltigerSync("https://github.com/GhostManager/tiger_sync", ""); err != nil {
			fmt.Printf("[-] Failed to install service: %v\n", err)
			os.Exit(1)
		}
	} else {
		if len(args) == 2 {
			branch = args[1]
		}
		if err := internal.InstalltigerSync(args[0], branch); err != nil {
			fmt.Printf("[-] Failed to install service: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("[+] Successfully installed service!\n")
		}
	}

}
