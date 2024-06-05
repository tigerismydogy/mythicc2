package cmd

import (
	"fmt"
	"github.com/tigerMeta/tiger_CLI/cmd/config"
	"github.com/tigerMeta/tiger_CLI/cmd/internal"
	"github.com/tigerMeta/tiger_CLI/cmd/manager"
	"os"
)

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tiger-cli",
	Short: fmt.Sprintf("A command line interface for managing tiger. Version: %s", config.Version),
	Long: `tiger CLI is a command line interface for managing the tiger application and associated containers and services.
Commands are grouped by their use and all support '-h' for help.
For a list of available services to install, check out: https://tigermeta.github.io/overview/`,
}

var force bool
var branch string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Create or parse the Docker ``.env`` file
	config.Initialize()
	manager.Initialize()
	internal.Initialize()
}
