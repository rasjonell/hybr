package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	servicesComponentsCmd.Flags().StringVarP(
		&serviceCmdFlags.service,
		"service", "s", "",
		"Name of the service",
	)
	servicesComponentsCmd.MarkFlagRequired("service")
	servicesCmd.AddCommand(servicesComponentsCmd)
}

var servicesComponentsCmd = &cobra.Command{
	Use:    "components",
	Short:  "Show service components",
	Long:   "Shows docker components the servies is composed of",
	PreRun: checkRemoteHost,
	Run:    showServiceComponents,
}

func showServiceComponents(cmd *cobra.Command, args []string) {
	is, err := getService(serviceCmdFlags.service)
	if err != nil {
		fmt.Println(err)
		return
	}

	maxNameLength := len("Name")
	maxVersionLength := len("Version")
	maxStatusLength := len("Status")

	for _, c := range is.GetComponents() {
		if len(c.Name) > maxNameLength {
			maxNameLength = len(c.Name)
		}
		if len(c.Version) > maxVersionLength {
			maxVersionLength = len(c.Version)
		}
		statusDisplay := displayStatusIcon(c.Status) + " " + c.Status
		if len(statusDisplay) > maxStatusLength {
			maxStatusLength = len(statusDisplay)
		}
	}

	headerFormat := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds\n", maxNameLength, maxVersionLength, maxStatusLength)
	header := fmt.Sprintf(headerFormat, "Name", "Version", "Status")

	separator := strings.Repeat("-", maxNameLength+maxVersionLength+maxStatusLength+6) + "\n" // +6 for spaces

	dataRows := ""
	for _, c := range is.GetComponents() {
		statusDisplay := displayStatusIcon(c.Status) + c.Status
		rowFormat := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds\n", maxNameLength, maxVersionLength, maxStatusLength)
		dataRows += fmt.Sprintf(rowFormat, c.Name, c.Version, statusDisplay)
	}

	fmt.Println(separator + header + separator + dataRows)
}
