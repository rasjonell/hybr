package services

import (
	"fmt"
	"log"
)

func Stop(serviceName string) {
	service, err := initServiceAction(serviceName, "stopping", "stopped")
	if err != nil {
		log.Println(err)
		return
	}

	if err := StopService(serviceName); err != nil {
		log.Printf("%v\n", err)
		return
	}

	service.Status = "stopped"
	for _, c := range service.Components {
		c.Status = "stopped"
	}

	if err := installationRegistry.save(); err != nil {
		log.Printf("Failed to save installtion: %v", err)
		return
	}
}

func Restart(serviceName string) {
	service, err := initServiceAction(serviceName, "restarting")
	if err != nil {
		log.Println(err)
		return
	}

	if err := RestartService(serviceName); err != nil {
		log.Printf("%v\n", err)
		return
	}

	service.Status = "running"
	for _, c := range service.Components {
		c.Status = "running"
	}

	if err := installationRegistry.save(); err != nil {
		log.Printf("Failed to save installtion: %v", err)
		return
	}
}

func UpdateVars(serviceName string, updatedVars map[string][]*VariableDefinition) {
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
		log.Printf("%v\n", err)
		return
	}

	if err := RestartService(serviceName); err != nil {
		log.Printf("%v\n", err)
		return
	}

	service.Status = "running"
	for _, c := range service.Components {
		c.Status = "running"
	}

	if err := installationRegistry.save(); err != nil {
		log.Printf("Failed to save installtion: %v", err)
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
