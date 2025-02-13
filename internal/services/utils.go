package services

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetRegisteredServices() []HybrService {
	mu.RLock()
	defer mu.RUnlock()

	services := make([]HybrService, 0, len(registry))
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

func findPort(varDef []*VariableDefinition) string {
	var port string
	for _, def := range varDef {
		if def.Name == "PORT" {
			port = def.Value
		}
	}

	return port
}
