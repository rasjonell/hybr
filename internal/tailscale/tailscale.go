package tailscale

import (
	"encoding/json"
	"fmt"
	"github.com/rasjonell/hybr/internal/system"
	"os/exec"
	"strings"
)

var magicDns string

func RunOnRemote(remoteHost string, remoteCmd string) (err error) {
	fmt.Printf("Running [%s] on %s...\n\n", remoteCmd, remoteHost)

	args := []string{
		"ssh",
		remoteHost,
		remoteCmd,
	}
	cmd := exec.Command("tailscale", args...)
	if err = system.PipeCmdToStdout(cmd, remoteHost); err != nil {
		return err
	}

	return nil
}

func Start(authKey string) (err error) {
	args := []string{
		"-p", "Enter your password to start tailscale:",
		"tailscale",
		"up",
		"--ssh",
	}
	if authKey != "" {
		args = append(args, "--auth-key", authKey)
	}
	cmd := exec.Command("sudo", args...)

	if err = system.PipeCmdToStdout(cmd, "tailscale"); err != nil {
		return err
	}

	magicDns, err = retrieveMagicDNS()
	if err != nil {
		return err
	}

	return nil
}

func Stop() error {
	cmd := exec.Command("sudo", "-p", "Enter your password to stop tailscale:", "tailscale", "down")
	if err := system.PipeCmdToStdout(cmd, "tailscale"); err != nil {
		return err
	}

	return nil
}

func AddServeTunnel(hybrProxy, port, tailscaleProxy string) (string, error) {
	path := "/" + strings.TrimPrefix(tailscaleProxy, "/")
	hybrPath := strings.TrimPrefix(hybrProxy, "/")

	cmd := exec.Command(
		"sudo",
		"-p", "Enter your password to add tailscale serve tunnel:",
		"tailscale", "serve",
		"--bg", "--set-path", path,
		fmt.Sprintf("localhost:%s/%s", port, hybrPath),
	)

	if err := system.PipeCmdToStdout(cmd, "tailscale"); err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s%s", magicDns, path), nil
}

func GetDNSName() string {
	return magicDns
}

func retrieveMagicDNS() (string, error) {
	cmd := exec.Command("tailscale", "status", "--json")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Failed to run tailscale status: %v", err)
	}

	var status struct {
		Self struct {
			DNSName string `json:"DNSName"`
		} `json:"Self"`
	}

	if err := json.Unmarshal(output, &status); err != nil {
		return "", fmt.Errorf("Failed to parse tailscale status output: %v", err)
	}

	if status.Self.DNSName == "" {
		return "", fmt.Errorf("Failed to find DNS Name in tailscale statuss")
	}

	return strings.TrimSuffix(status.Self.DNSName, "."), nil
}
