package cmd

import (
	"fmt"
	"hybr/internal/services"

	"github.com/spf13/cobra"
)

func init() {
	servicesStartCmd.Flags().StringVarP(
		&serviceCmdFlags.service,
		"service", "s", "",
		"Name of the service",
	)
	servicesStartCmd.MarkFlagRequired("service")
	servicesCmd.AddCommand(servicesStartCmd)
}

var servicesStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the service",
	Long: `Starts the service if it is NOT running.
Restarts the services otherwise`,
	PreRun: checkRemoteHost,
	Run:    startService,
}

func startService(cmd *cobra.Command, args []string) {
	checkRootPermissions()
	is, err := getService(serviceCmdFlags.service)
	if err != nil {
		fmt.Println(err)
		return
	}

	if is.GetStatus() == "running" {
		services.Restart(is.GetName())
	} else {
		services.Start(is.GetName())
	}
}
