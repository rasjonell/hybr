package services

import (
	"log"
)

func UpdateVars(serviceName string, updatedVars map[string][]*VariableDefinition) {
	installation, exists := installationRegistry.GetInstallation(serviceName)
	if !exists {
		log.Printf("Service doesn't exist")
		return
	}
	service, ok := installation.(*serviceImpl)
	if !ok {
		log.Printf("Invalid service type")
		return
	}

	service.Status = "restarting"
	for _, c := range service.Components {
		c.Status = "restarting"
	}

	if err := installationRegistry.save(); err != nil {
		log.Printf("Failed to save installtion: %v", err)
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
