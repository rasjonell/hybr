package services

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	HybrDir = ".hybr"
)

func getWorkingDirectory() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, HybrDir)
}

func initWorkingDirectory() error {
	path := getWorkingDirectory()

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	return nil
}

func installServices(selected map[string]SelectedServiceModel) error {
	basePath := filepath.Join(getWorkingDirectory(), "services")

	for _, service := range selected {
		servicePath := filepath.Join(basePath, service.Name)
		if err := os.MkdirAll(servicePath, 0755); err != nil {
			return err
		}

		for _, tpl := range service.Templates {
			if err := processTemplate(servicePath, tpl, service.Variables); err != nil {
				return err
			}
		}
	}

	return nil
}

func processTemplate(basePath string, tpl Template, vars []VariableDefinition) error {
	_, _, _ = basePath, tpl, vars
	// TODO implement template processing
	return fmt.Errorf("Not Implemented")
}
