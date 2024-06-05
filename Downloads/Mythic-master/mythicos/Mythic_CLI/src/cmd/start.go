package cmd

import (
	"github.com/tigerMeta/tiger_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var startCmd = &cobra.Command{
	Use:   "start [container names]",
	Short: "Start tiger containers",
	Long:  `Run this command to start all tiger containers. If you want to only start certain containers, specify their names.`,
	Run:   start,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func start(cmd *cobra.Command, args []string) {
	if err := internal.ServiceStart(args); err != nil {

	}
}
