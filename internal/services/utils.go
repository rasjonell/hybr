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

func GetInstallableServicePaths() (directoryNames []string, serviceNames []string, err error) {
	servicesPath := filepath.Join(GetHybrDirectory(), "services")
	serviceNames, err = GetInstallableServices()
	if err != nil {
		return nil, nil, err
	}

	for _, name := range serviceNames {
		directoryNames = append(directoryNames, filepath.Join(servicesPath, name))
	}

	return directoryNames, serviceNames, nil
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

func GetInstallableServices() ([]string, error) {
	servicesPath := filepath.Join(GetHybrDirectory(), "services")

	files, err := os.ReadDir(servicesPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read directory %s: %w", servicesPath, err)
	}

	var installableServices []string
	for _, file := range files {
		if file.IsDir() {
			installableServices = append(installableServices, file.Name())
		}
	}

	return installableServices, nil
}
