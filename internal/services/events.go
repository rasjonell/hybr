package services

import (
	"bufio"
	"fmt"
	"hybr/internal/orchestration"
	"os/exec"
	"path/filepath"
	"time"
)

type ServiceLogMonitor struct {
	ServiceName string
	EventType   orchestration.EventType
}

type ServiceStatusMonitor struct {
	ServiceName string
	EventType   orchestration.EventType
}

type ServiceComponentStatusMonitor struct {
	ServiceName string
	EventType   orchestration.EventType
}

func GetServiceLogEvent(name string) orchestration.EventType {
	return orchestration.EventType(fmt.Sprintf("%s_service_log_event", name))
}

func GetServiceStatusEvent(name string) orchestration.EventType {
	return orchestration.EventType(fmt.Sprintf("%s_service_status_event", name))
}

func GetServiceComponentStatusEvent(name string) orchestration.EventType {
	return orchestration.EventType(fmt.Sprintf("%s_service_component_status_event", name))
}

func RegisterEventSources(name string) {
	logEventType := GetServiceLogEvent(name)
	statusEventType := GetServiceStatusEvent(name)
	componentStatusEventType := GetServiceComponentStatusEvent(name)

	subManager := orchestration.GetSubscriptionManager()

	subManager.RegisterEventSource(logEventType, &ServiceLogMonitor{
		ServiceName: name,
		EventType:   logEventType,
	})

	subManager.RegisterEventSource(statusEventType, &ServiceStatusMonitor{
		ServiceName: name,
		EventType:   statusEventType,
	})

	subManager.RegisterEventSource(componentStatusEventType, &ServiceComponentStatusMonitor{
		ServiceName: name,
		EventType:   componentStatusEventType,
	})
}

func (m *ServiceLogMonitor) Start(doneChan <-chan struct{}, logChan chan<- *orchestration.EventChannelData) {
	cmd := exec.Command("docker", "compose", "logs",
		"-f",
		"--no-color",
		"--tail", "10",
	)
	cmd.Dir = filepath.Join(GetHybrDirectory(), "services", m.ServiceName)

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

func (m *ServiceStatusMonitor) Start(doneChan <-chan struct{}, eventChan chan<- *orchestration.EventChannelData) {
	for {
		select {
		case <-doneChan:
			return
		default:
			ir := GetRegistry()
			fmt.Println("Service " + m.ServiceName + " Status is " + ir.installations[m.ServiceName].Status)
			eventChan <- orchestration.ToEventData(m.EventType, ir.installations[m.ServiceName].Status)
			time.Sleep(1 * time.Second)
		}
	}
}

func (m *ServiceComponentStatusMonitor) Start(doneChan <-chan struct{}, eventChan chan<- *orchestration.EventChannelData) {
	for {
		select {
		case <-doneChan:
			return
		default:
			ir := GetRegistry()
			for _, comp := range ir.installations[m.ServiceName].Components {
				eventChan <- orchestration.ToEventData(m.EventType, comp.Status, map[string]string{"ComponentName": comp.Name})
			}
			time.Sleep(1 * time.Second)
		}
	}
}
