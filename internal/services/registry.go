package services

import (
	"embed"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var (
	registry = make(map[string]Service)
	mu       sync.RWMutex
)

type Service struct {
	Name           string                `json:"name"`
	Description    string                `json:"description"`
	InstallCommand string                `json:"installCommand"`
	Templates      []string              `json:"templates"`
	Variables      map[string][]Variable `json:"variables"`
}

type Variable struct {
	Name        string `json:"name"`
	Default     string `json:"default"`
	Description string `json:"description"`
}

type Template struct {
	SourcePath string `json:"sourcePath"`
	TargetName string `json:"targetName"`
}

type SelectedServiceModel struct {
	ServiceName    string
	InstallCommand string
	Variables      map[string][]*VariableDefinition
}

type VariableDefinition struct {
	Key   string
	Value string
}

func register(s Service) {
	mu.Lock()
	defer mu.Unlock()

	registry[s.Name] = s
}

//go:embed templates/services.json
var defaultJsonData []byte

//go:embed templates/**/*
var templatesFS embed.FS

func getServices(forceResetTemplates bool) []Service {
	var err error = nil
	var services []Service

	servicesPath := filepath.Join(getWorkingDirectory(), "services")
	destPath := filepath.Join(getWorkingDirectory(), "services.json")

	_, err = os.Stat(destPath)
	if !forceResetTemplates && err == nil {
		data, err := os.ReadFile(destPath)
		if err != nil {
			panic("Unable To Read services.json")
		}

		if err := json.Unmarshal(data, &services); err != nil {
			panic(err)
		}
	}

	if forceResetTemplates || os.IsNotExist(err) {
		if err := os.WriteFile(destPath, defaultJsonData, 0644); err != nil {
			panic(err)
		}

		if err := json.Unmarshal(defaultJsonData, &services); err != nil {
			panic(err)
		}

		for _, s := range services {
			tPath := filepath.Join(servicesPath, s.Name, "templates")
			if err := os.MkdirAll(tPath, 0755); err != nil {
				panic(err)
			}

			for _, templName := range s.Templates {
				sourceFile, err := templatesFS.Open(filepath.Join("templates", s.Name, templName))
				if err != nil {
					panic(err)
				}
				defer sourceFile.Close()

				destFile, err := os.Create(filepath.Join(tPath, templName))
				if err != nil {
					panic(err)
				}
				defer destFile.Close()

				_, err = io.Copy(destFile, sourceFile)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	return services
}

func InitRegistry(forceResetTemplates bool) {
	//TODO: if forceResetTemplates, empty $workingDirectory/services
	initWorkingDirectory()
	services := getServices(forceResetTemplates)

	for _, service := range services {
		register(service)
	}
}
