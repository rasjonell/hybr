package cmd

import (
	"fmt"
	"hybr/internal/services"
	"os"
	"os/exec"
	"strconv"
)

func checkRootPermissions(msgs ...string) {
	msg := "You need root privileges to run this program\nPlease run with sudo\n"
	if len(msgs) != 0 {
		msg = msgs[0] + "\n"
	}
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	i, err := strconv.Atoi(string(output[:len(output)-1]))

	if err != nil {
		panic(err)
	}

	if i != 0 {
		fmt.Printf(msg + "\n\tsudo hybr [command]\n\n")
		os.Exit(1)
	}
}

func getService(serviceName string) (services.HybrService, error) {
	reg := services.GetRegistry()
	is, exists := reg.GetInstallation(serviceName)
	if !exists {
		return nil, fmt.Errorf(`No service installed with the given name: %s
To see the list of services run:
      hybr services
`, serviceCmdFlags.service)
	}

	return is, nil
}

func displayStatusIcon(status string) string {
	statusIcon := "✅"
	if status != "running" {
		statusIcon = "🛑"
	}

	return statusIcon
}
