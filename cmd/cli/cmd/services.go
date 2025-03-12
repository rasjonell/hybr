package cmd

import (
	"fmt"
	"hybr/internal/services"
	"strings"

	"github.com/spf13/cobra"
)

type serviceCmdFlag struct {
	service string
	remote  string
}

var serviceCmdFlags serviceCmdFlag

func init() {
	servicesCmd.Flags().StringVarP(
		&serviceCmdFlags.service,
		"service", "s", "",
		"Name of the service",
	)

	servicesCmd.PersistentFlags().StringVarP(
		&serviceCmdFlags.remote,
		"remote-host", "r", "",
		"Hostname of a remote hybr node",
	)

	rootCmd.AddCommand(servicesCmd)
}

var servicesCmd = &cobra.Command{
	Use:    "services",
	Short:  "Show services info",
	Long:   "Show installed services with their status",
	PreRun: checkRemoteHost,
	Run:    listServices,
}

func listServices(cmd *cobra.Command, args []string) {
	reg := services.GetRegistry()
	reg.RegisterServiceEvents()

	result := ""

	if serviceCmdFlags.service == "" {
		result = listAllServices()
	} else {
		result = listSingleService(serviceCmdFlags.service)
	}

	fmt.Printf(result + "\n\n")
}

func listAllServices() string {
	reg := services.GetRegistry()
	installedServices := reg.ListInstallations()

	maxNameLength := 0
	for _, is := range installedServices {
		if len(is.GetName()) > maxNameLength {
			maxNameLength = len(is.GetName())
		}
	}

	lines := make([]string, len(installedServices))

	for i, is := range installedServices {
		padding := maxNameLength - len(is.GetName())
		name := fmt.Sprintf("[%s]%s", is.GetName(), strings.Repeat(" ", padding))
		lines[i] = fmt.Sprintf(
			"%s  %s%s", name, displayStatusIcon(is.GetStatus()), is.GetStatus(),
		)
	}

	return fmt.Sprintf(strings.Join(lines, "\n"))
}

func listSingleService(serviceName string) string {
	is, err := getService(serviceName)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf(
		"[%s]  %s%s", is.GetName(),
		displayStatusIcon(is.GetStatus()), is.GetStatus(),
	)
}
