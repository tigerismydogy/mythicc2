package cmd

import (
	"github.com/tigerMeta/tiger_CLI/cmd/internal"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test tiger service connections",
	Long:  `Run this command to test tiger connections to RabbitMQ and the tiger UI`,
	Run:   test,
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func test(cmd *cobra.Command, args []string) {
	internal.TesttigerRabbitmqConnection()
	internal.TesttigerConnection()
}
