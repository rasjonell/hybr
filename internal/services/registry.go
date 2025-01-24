package services

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

var (
	registry = make(map[string]Service)
	mu       sync.RWMutex
)

type Service struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Variables   []Variable `json:"variables"`
	Templates   []Template `json:"templates"`
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
	*Service
	Variables []VariableDefinition
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

func getServices(forceResetTemplates bool) []Service {
	var err error = nil
	var services []Service
	destPath := filepath.Join(getWorkingDirectory(), "services.json")

	if forceResetTemplates {
	}

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
	}

	return services
}

func InitRegistry(forceResetTemplates bool) {
	initWorkingDirectory()
	services := getServices(forceResetTemplates)

	for _, service := range services {
		register(service)
	}
}
