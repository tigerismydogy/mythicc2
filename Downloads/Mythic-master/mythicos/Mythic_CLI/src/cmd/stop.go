package cmd

import (
	"github.com/tigerMeta/tiger_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all of tiger",
	Long: `Run this command stop all tiger containers. Use subcommands to
adjust specific containers to stop.`,
	Run: stop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func stop(cmd *cobra.Command, args []string) {
	if err := internal.ServiceStop(args); err != nil {

	}
}
