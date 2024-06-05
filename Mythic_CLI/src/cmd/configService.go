package cmd

import (
	"fmt"
	"github.com/tigerMeta/tiger_CLI/cmd/config"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

// configServiceCmd
var configServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Get configurations for remote services",
	Long: `Get configuration variables to use with a remote service - 
a service that runs on a host other than the host where tiger is running`,
	Run: configService,
}

func init() {
	configCmd.AddCommand(configServiceCmd)
}

func configService(cmd *cobra.Command, args []string) {
	// initialize tabwriter
	writer := new(tabwriter.Writer)
	// Set minwidth, tabwidth, padding, padchar, and flags
	writer.Init(os.Stdout, 8, 8, 1, '\t', 0)

	defer writer.Flush()

	fmt.Println("[+] Getting configuration values:")
	fmt.Fprintf(writer, "\n %s\t%s", "Setting", "Value")
	fmt.Fprintf(writer, "\n %s\t%s", "–––––––", "–––––––")

	configuration := config.GetConfigStrings([]string{
		"tiger_SERVER_HOST",
		"tiger_SERVER_PORT",
		"tiger_SERVER_GRPC_PORT",
		"RABBITMQ_HOST",
		"RABBITMQ_PASSWORD",
		"RABBITMQ_PORT",
	})
	keys := make([]string, 0, len(configuration))
	for k := range configuration {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(writer, "\n %s\t%s", strings.ToUpper(key), configuration[key])
	}
	tigerServerStatus := config.GetConfigStrings([]string{"tiger_SERVER_BIND_LOCALHOST_ONLY"})
	if val, ok := tigerServerStatus["tiger_SERVER_BIND_LOCALHOST_ONLY"]; ok {
		if val == "true" {
			fmt.Fprintf(writer, "\t\t")
			fmt.Fprintf(writer, "tiger_SERVER_BIND_LOCALHOST_ONLY is set to true - set this to false and restart tiger")
		}
	}
	rabbitmqStatus := config.GetConfigStrings([]string{"RABBITMQ_BIND_LOCALHOST_ONLY"})
	if val, ok := rabbitmqStatus["RABBITMQ_BIND_LOCALHOST_ONLY"]; ok {
		if val == "true" {
			fmt.Fprintf(writer, "\t\t")
			fmt.Fprintf(writer, "RABBITMQ_BIND_LOCALHOST_ONLY is set to true - set this to false and restart tiger")
		}
	}
	fmt.Fprintln(writer, "")
}
