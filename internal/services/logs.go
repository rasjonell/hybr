package services

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
)

func FollowLogs(doneChan <-chan struct{}, logChan chan<- string, serviceName string) {
	defer close(logChan)

	cmd := exec.Command("docker", "compose", "logs",
		"-f",
		"--no-color",
		"--tail", "10",
	)
	cmd.Dir = filepath.Join(getWorkingDirectory(), "services", serviceName)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Failed creating stdout pipe: %v\n", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		select {
		case <-doneChan:
			cmd.Process.Kill()
			return

		default:
			logChan <- scanner.Text()
		}
	}
}
