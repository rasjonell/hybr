package cmd

import (
	"fmt"
	"strings"

	"github.com/rasjonell/hybr/internal/services"
	"github.com/spf13/cobra"
)

func init() {
	servicesStatusCmd.Flags().StringVarP(
		&serviceCmdFlags.service,
		"service", "s", "",
		"Name of the service",
	)
	servicesCmd.AddCommand(servicesStatusCmd)
}

var servicesStatusCmd = &cobra.Command{
	Use:    "status",
	Short:  "Show service status",
	Long:   "Shows service information",
	PreRun: checkRemoteHost,
	Run:    showServicesStatus,
}

func showServicesStatus(cmd *cobra.Command, args []string) {
	reg := services.GetRegistry()
	reg.RegisterServiceEvents()

	result := ""

	if serviceCmdFlags.service == "" {
		result = installedServicesStatus()
	} else {
		result = singleInstalledServiceStatus(serviceCmdFlags.service)
	}

	fmt.Printf(result + "\n\n")
}

func installedServicesStatus() string {
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

func singleInstalledServiceStatus(serviceName string) string {
	is, err := getService(serviceName)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf(
		"[%s]  %s%s", is.GetName(),
		displayStatusIcon(is.GetStatus()), is.GetStatus(),
	)
}
