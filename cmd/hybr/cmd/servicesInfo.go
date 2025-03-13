package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	servicesInfoCmd.Flags().StringVarP(
		&serviceCmdFlags.service,
		"service", "s", "",
		"Name of the service",
	)
	servicesInfoCmd.MarkFlagRequired("service")
	servicesCmd.AddCommand(servicesInfoCmd)
}

var servicesInfoCmd = &cobra.Command{
	Use:    "info",
	Short:  "Show service information",
	Long:   "Shows service configuration",
	PreRun: checkRemoteHost,
	Run:    showServiceInfo,
}

func showServiceInfo(cmd *cobra.Command, args []string) {
	is, err := getService(serviceCmdFlags.service)
	if err != nil {
		fmt.Println(err)
		return
	}

	keys := []string{
		"Name",
		"Status",
		"Global URL",
		"Local URL",
		"Install Date",
		"Last Start Date",
	}
	values := []string{
		is.GetName(),
		is.GetStatus(),
		is.GetURL(),
		"localhost:" + is.GetPort(),
		is.GetInstallDate().Format(time.RFC850),
		is.GetLastStartTime().Format(time.RFC850),
	}

	maxNameLength := len("Last Start Date")

	lines := make([]string, 0)

	for i, key := range keys {
		value := values[i]
		padding := maxNameLength - len(key)
		name := fmt.Sprintf("[%s]%s", key, strings.Repeat(" ", padding))
		result := fmt.Sprintf(
			"%s  %s", name, value,
		)
		if key == "Status" {
			result = fmt.Sprintf(
				"%s  %s%s", name, displayStatusIcon(value), value,
			)
		}

		lines = append(lines, result)
	}

	fmt.Printf(strings.Join(lines, "\n") + "\n\n")
}
