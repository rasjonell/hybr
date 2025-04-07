package cmd

import (
	"fmt"
	"github.com/rasjonell/hybr/internal/services"

	"github.com/spf13/cobra"
)

func init() {
	servicesStopCmd.Flags().StringVarP(
		&serviceCmdFlags.service,
		"service", "s", "",
		"Name of the service",
	)
	servicesStopCmd.MarkFlagRequired("service")
	servicesCmd.AddCommand(servicesStopCmd)
}

var servicesStopCmd = &cobra.Command{
	Use:    "stop",
	Short:  "Stop the service",
	Long:   "Stop the service",
	PreRun: checkRemoteHost,
	Run:    stopService,
}

func stopService(cmd *cobra.Command, args []string) {
	is, err := getService(serviceCmdFlags.service)
	if err != nil {
		fmt.Println(err)
		return
	}

	if is.GetStatus() != "stopped" {
		services.Stop(is.GetName())
	} else {
		fmt.Printf("Service [%s] is already stopped\n", is.GetName())
	}
}
