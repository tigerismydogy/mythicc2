package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Interact with the tiger volumes",
	Long:  `Run this command to interact with the tiger volumes`,
	Run:   volumes,
}
var volumeName string
var sourceName string
var destinationName string

func init() {
	rootCmd.AddCommand(volumeCmd)
}

func volumes(cmd *cobra.Command, args []string) {

}
