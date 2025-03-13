package services

import (
	"fmt"
	"github.com/rasjonell/hybr/internal/orchestration"
	"log"
	"strings"
)

func Stop(serviceName string) {
	orchestration.SendWarningNotification(
		fmt.Sprintf("Stoping %s Service", strings.ToTitle(serviceName)),
	)

	service, err := initServiceAction(serviceName, "stopping", "stopped")
	if err != nil {
		orchestration.SendErrorNotification(
			fmt.Sprintf("%s", err.Error()),
		)
		return
	}

	if err := StopService(serviceName); err != nil {
		orchestration.SendErrorNotification(
			fmt.Sprintf("%s", err.Error()),
		)
		return
	}

	service.Status = "stopped"
	for _, c := range service.Components {
		c.Status = "stopped"
	}

	if err := installationRegistry.save(); err != nil {
		orchestration.SendErrorNotification(
			fmt.Sprintf("%s", err.Error()),
		)
		return
	}

	orchestration.SendSuccessNotification(
		fmt.Sprintf("%s Service Stopped", strings.ToTitle(serviceName)),
	)
}

func Start(serviceName string) {
	orchestration.SendInfoNotification(
		fmt.Sprintf("Starting %s Service", strings.ToTitle(serviceName)),
	)

	service, err := initServiceAction(serviceName, "starting")
	if err != nil {
		orchestration.SendErrorNotification(
			fmt.Sprintf("%s", err.Error()),
		)
		return
	}

	if err := StartService(serviceName); err != nil {
		orchestration.SendErrorNotification(
			fmt.Sprintf("%s", err.Error()),
		)
		return
	}

	service.Status = "running"
	for _, c := range service.Components {
		c.Status = "running"
	}

	if err := installationRegistry.save(); err != nil {
		orchestration.SendErrorNotification(
			fmt.Sprintf("%s", err.Error()),
		)
		return
	}

	orchestration.SendSuccessNotification(
		fmt.Sprintf("%s Service Started", strings.ToTitle(serviceName)),
	)
}

func Restart(serviceName string) {
	orchestration.SendInfoNotification(
		fmt.Sprintf("Restarting %s Service", strings.ToTitle(serviceName)),
	)

	service, err := initServiceAction(serviceName, "restarting")
	if err != nil {
		orchestration.SendErrorNotification(
			fmt.Sprintf("%s", err.Error()),
		)
		return
	}

	if err := RestartService(serviceName); err != nil {
		orchestration.SendErrorNotification(
			fmt.Sprintf("%s", err.Error()),
		)
		return
	}

	service.Status = "running"
	for _, c := range service.Components {
		c.Status = "running"
	}

	if err := installationRegistry.save(); err != nil {
		orchestration.SendErrorNotification(
			fmt.Sprintf("%s", err.Error()),
		)
		return
	}

	orchestration.SendSuccessNotification(
		fmt.Sprintf("%s Service Restarted", strings.ToTitle(serviceName)),
	)
}

func UpdateVars(serviceName string, updatedVars map[string][]*VariableDefinition) {
	orchestration.SendInfoNotification(
		fmt.Sprintf("Upadting %s Service Variables", strings.ToTitle(serviceName)),
	)

	service, err := initServiceAction(serviceName, "restarting")
	if err != nil {
		log.Println(err)
		return
	}

	for fileName, vars := range updatedVars {
		for _, updatedVar := range vars {
			for _, existingVar := range service.Variables[fileName] {
				if existingVar.Name == updatedVar.Name {
					existingVar.Value = updatedVar.Value
					break
				}
			}
		}
	}

	if err := reinstallTemplates(service, serviceName); err != nil {
		orchestration.SendErrorNotification(
			fmt.Sprintf("%s", err.Error()),
		)
		return
	}

	if err := RestartService(serviceName); err != nil {
		orchestration.SendErrorNotification(
			fmt.Sprintf("%s", err.Error()),
		)
		return
	}

	service.Status = "running"
	for _, c := range service.Components {
		c.Status = "running"
	}

	if err := installationRegistry.save(); err != nil {
		orchestration.SendErrorNotification(
			fmt.Sprintf("%s", err.Error()),
		)
		return
	}
}

func initServiceAction(serviceName, initStatus string, exitStatuses ...string) (*serviceImpl, error) {
	exitStatus := initStatus
	if len(exitStatuses) != 0 {
		exitStatus = exitStatuses[0]
	}

	installation, exists := installationRegistry.GetInstallation(serviceName)
	if !exists {
		return nil, fmt.Errorf("Service doesn't exist")
	}

	if installation.GetStatus() == exitStatus {
		return nil, fmt.Errorf("Exit status triggered")
	}

	service, ok := installation.(*serviceImpl)
	if !ok {
		return nil, fmt.Errorf("Invalid service type")
	}

	service.Status = initStatus
	for _, c := range service.Components {
		c.Status = initStatus
	}

	if err := installationRegistry.save(); err != nil {
		return nil, fmt.Errorf("Failed to save installtion: %v", err)
	}

	return service, nil
}
