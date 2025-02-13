package services

import (
	"embed"
	"encoding/json"
	"hybr/internal/docker"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

var (
	registry = make(map[string]*serviceImpl)
	mu       sync.RWMutex
)

type HybrService interface {
	GetName() string
	GetDescription() string
	IsSubDomain() bool
	GetTemplates() []string
	GetVariables() map[string][]*VariableDefinition
	GetStatus() string
	GetPort() string
	GetURL() string
	GetComponents() []*docker.Component
	GetInstallDate() time.Time
	GetLastStartTime() time.Time
}

type serviceImpl struct {
	Name        string                           `json:"name"`
	Description string                           `json:"description"`
	SubDomain   bool                             `json:"subDomain"`
	Templates   []string                         `json:"templates"`
	Variables   map[string][]*VariableDefinition `json:"variables"`

	Status        string
	Port          string
	URL           string
	InstallDate   time.Time
	LastStartTime time.Time
	Components    []*docker.Component
}

type VariableDefinition struct {
	Name        string `json:"name"`
	Default     string `json:"default"`
	Description string `json:"description"`

	Value    string
	Template string
	Input    textinput.Model `json:"-"`
}

func (s *serviceImpl) GetName() string {
	return s.Name
}

func (s *serviceImpl) GetDescription() string {
	return s.Description
}

func (s *serviceImpl) IsSubDomain() bool {
	return s.SubDomain
}

func (s *serviceImpl) GetTemplates() []string {
	return s.Templates
}

func (s *serviceImpl) GetVariables() map[string][]*VariableDefinition {
	return s.Variables
}

func (s *serviceImpl) GetStatus() string {
	return s.Status
}

func (s *serviceImpl) GetPort() string {
	return s.Port
}

func (s *serviceImpl) GetURL() string {
	return s.URL
}

func (s *serviceImpl) GetComponents() []*docker.Component {
	return s.Components
}

func (s *serviceImpl) GetInstallDate() time.Time {
	return s.InstallDate
}

func (s *serviceImpl) GetLastStartTime() time.Time {
	return s.LastStartTime
}

type Template struct {
	SourcePath string `json:"sourcePath"`
	TargetName string `json:"targetName"`
}

func register(s *serviceImpl) {
	mu.Lock()
	defer mu.Unlock()

	registry[s.GetName()] = s
}

//go:embed templates/services.json
var defaultJsonData []byte

//go:embed templates/**/*
var templatesFS embed.FS

func initializeServices() []*serviceImpl {
	var err error = nil
	var services []*serviceImpl

	servicesPath := filepath.Join(getWorkingDirectory(), "services")
	destPath := filepath.Join(getWorkingDirectory(), "services.json")

	_, err = os.Stat(destPath)
	if err == nil {
		data, err := os.ReadFile(destPath)
		if err != nil {
			panic("Unable To Read services.json")
		}

		if err := json.Unmarshal(data, &services); err != nil {
			panic(err)
		}
	}

	if os.IsNotExist(err) {
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
