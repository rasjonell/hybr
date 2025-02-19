package services

import (
	"bufio"
	"fmt"
	"hybr/internal/orchestration"
	"os/exec"
	"path/filepath"
)

type ServiceLogMonitor struct {
	ServiceName string
	EventType   orchestration.EventType
}

func GetServiceLogEvent(name string) orchestration.EventType {
	return orchestration.EventType(fmt.Sprintf("%s_service_log_event", name))
}

func RegisterServiceLogMonitor(name string) {
	logEventType := GetServiceLogEvent(name)

	subManager := orchestration.GetSubscriptionManager()
	subManager.RegisterEventSource(logEventType, &ServiceLogMonitor{
		ServiceName: name,
		EventType:   logEventType,
	})
}

func (m *ServiceLogMonitor) Start(doneChan <-chan struct{}, logChan chan<- *orchestration.EventChannelData) {
	cmd := exec.Command("docker", "compose", "logs",
		"-f",
		"--no-color",
		"--tail", "10",
	)
	cmd.Dir = filepath.Join(getWorkingDirectory(), "services", m.ServiceName)

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
			logChan <- orchestration.ToEventData(m.EventType, scanner.Text())
		}
	}
}
