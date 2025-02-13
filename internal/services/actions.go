package services

import (
	"fmt"
)

func UpdateVars(serviceName string, updatedVars map[string][]*VariableDefinition) error {
	installation, exists := installationRegistry.GetInstallation(serviceName)
	if !exists {
		return fmt.Errorf("service %s not found", serviceName)
	}
	service, ok := installation.(*serviceImpl)
	if !ok {
		return fmt.Errorf("Invalid service type")
	}

	service.Status = "restarting"
	for _, c := range service.Components {
		c.Status = "restarting"
	}

	if err := installationRegistry.save(); err != nil {
		return fmt.Errorf("Failed to save installtion: %w", err)
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
		return err
	}

	if err := RestartService(serviceName); err != nil {
		return err
	}

	service.Status = "running"
	for _, c := range service.Components {
		c.Status = "running"
	}

	if err := installationRegistry.save(); err != nil {
		return fmt.Errorf("Failed to save installtion: %w", err)
	}

	return nil
}
