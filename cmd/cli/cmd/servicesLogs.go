package cmd

import (
	"fmt"
	"hybr/internal/nginx"
	"hybr/internal/services"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	servicesLogsCmd.Flags().StringVarP(
		&serviceCmdFlags.service,
		"service", "s", "",
		"Name of the service",
	)
	servicesLogsCmd.MarkFlagRequired("service")
	servicesCmd.AddCommand(servicesLogsCmd)
}

var servicesLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Service logs",
	Long:  "Shows docker compose logs for the service",
	Run:   logService,
}

func logService(_ *cobra.Command, args []string) {
	is, err := getService(serviceCmdFlags.service)
	if err != nil {
		fmt.Println(err)
		return
	}

	if is.GetStatus() != "running" {
		fmt.Printf("Service [%s] is not running\n", is.GetName())
		fmt.Printf("Check the service status with\n\thybr services -s %s info\n", is.GetName())
		return
	}

	cmd := exec.Command("docker", "compose", "logs",
		"-f",
		"--no-color",
		"--tail", "10",
	)
	cmd.Dir = filepath.Join(services.GetHybrDirectory(), "services", is.GetName())

	nginx.PipeCmdToStdout(cmd, is.GetName())
}
