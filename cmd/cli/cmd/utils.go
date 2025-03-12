package cmd

import (
	"fmt"
	"hybr/internal/services"
	"hybr/internal/tailscale"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

func checkRemoteHost(ccmd *cobra.Command, args []string) {
	if serviceCmdFlags.remote == "" {
		checkHybrInstallation(ccmd)
		return
	}

	child := ccmd
	cmds := []string{getFullCmd(child)}
	for child.HasParent() == true {
		cmds = append(cmds, getFullCmd(child.Parent()))
		child = child.Parent()
	}
	slices.Reverse(cmds)

	if err := tailscale.RunOnRemote(
		serviceCmdFlags.remote,
		strings.Join(cmds, " "),
	); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

func getFullCmd(cmd *cobra.Command) string {
	finalCmd := fmt.Sprintf("%s ", cmd.Name())
	cmd.Flags().Visit(func(f *pflag.Flag) {
		// Omit remote host from remote execution command
		if f.Shorthand != "r" {
			finalCmd += fmt.Sprintf("-%s %s ", f.Shorthand, f.Value)
		}
	})

	return strings.TrimSpace(finalCmd)
}

func checkHybrInstallation(cmd *cobra.Command) {
	_, err := os.Stat(filepath.Join(services.HybrDir, "installations.json"))
	if os.IsNotExist(err) {
		fmt.Println("\nLooks like hybr was not initialized on this system")
		fmt.Println("Please either initialize services by running:")
		fmt.Println("\thybr init")
		fmt.Printf("Or provide a remote hybr host using -r or --remote-host\n\n")
		cmd.Usage()
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
