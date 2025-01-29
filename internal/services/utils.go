package services

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func GetRegisteredServices() []Service {
	mu.RLock()
	defer mu.RUnlock()

	services := make([]Service, 0, len(registry))
	for _, s := range registry {
		services = append(services, s)
	}
	return services
}

func GetInstalledServices() []string {
	mu.RLock()
	defer mu.RUnlock()

	servicesDir := filepath.Join(getWorkingDirectory(), "services")

	var installedServices []string
	if _, err := os.Stat(servicesDir); os.IsNotExist(err) {
		fmt.Printf("Services directory not found")
		return installedServices
	}

	entries, err := os.ReadDir(servicesDir)
	if err != nil {
		fmt.Printf("Unable to read installed services %v\n", err)
		return installedServices
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		installedServices = append(installedServices, entry.Name())
	}

	return installedServices
}

func pipeCmdToStdout(cmd *exec.Cmd, label string) error {
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Error creating stdout pipe: %v", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Error creating stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Failed to start command for %s: %v", label, err)
	}

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			fmt.Printf("[%s] %s\n", label, scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			fmt.Fprintf(os.Stderr, "[%s] %s\n", label, scanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed for %s: %v", label, err)
	}

	return nil
}
