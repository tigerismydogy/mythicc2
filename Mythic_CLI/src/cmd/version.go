package cmd

import (
	"fmt"
	"github.com/tigerMeta/tiger_CLI/cmd/config"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print information about the tiger-cli and tiger versions",
	Long:  `Run this command to print versioning information about tiger and tiger-cli `,
	Run:   tigerVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func tigerVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("[*] tiger-cli version:    %s\n", config.Version)
	if fileContents, err := os.ReadFile("VERSION"); err != nil {
		fmt.Printf("[!] Failed to get tiger version: %v\n", err)
	} else {
		fmt.Printf("[*] tiger Server version: v%s\n", string(fileContents))
	}
	if fileContents, err := os.ReadFile(filepath.Join(".", "tigerReactUI", "src", "index.js")); err != nil {
		fmt.Printf("[!] Failed to get tigerReactUI version: %v\n", err)
	} else {
		fileLines := strings.Split(string(fileContents), "\n")
		for _, line := range fileLines {
			if strings.Contains(line, "tigerUIVersion") {
				uiVersionPieces := strings.Split(line, "=")
				uiVersion := uiVersionPieces[1]
				fmt.Printf("[*] React UI Version:      v%s\n", uiVersion[2:len(uiVersion)-2])
				return
			}
		}
	}
}
