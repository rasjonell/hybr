package services

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rasjonell/hybr/internal/docker"

	"github.com/charmbracelet/bubbles/textinput"
)

var (
	registry = make(map[string]*serviceImpl)
	mu       sync.RWMutex
)

type HybrService interface {
	GetName() string
	GetURL() string
	GetPort() string
	GetStatus() string
	GetDescription() string
	GetTemplates() []string
	GetHybrProxy() string
	GetTailscaleProxy() string
	GetInstallDate() time.Time
	GetLastStartTime() time.Time
	GetComponents() []*docker.Component
	GetVariables() map[string][]*VariableDefinition
}

type serviceImpl struct {
	Name           string                           `json:"name"`
	IsRoot         bool                             `json:"isRoot"`
	HybrProxy      string                           `json:"hybrProxy"`
	Description    string                           `json:"description"`
	TailscaleProxy string                           `json:"tailscaleProxy"`
	Templates      []string                         `json:"templates"`
	Variables      map[string][]*VariableDefinition `json:"variables"`

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

func (s *serviceImpl) GetHybrProxy() string {
	return s.HybrProxy
}

func (s *serviceImpl) GetName() string {
	return s.Name
}

func (s *serviceImpl) GetDescription() string {
	return s.Description
}

func (s *serviceImpl) GetTailscaleProxy() string {
	return s.TailscaleProxy
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

//go:embed all:templates/**/*
var templatesFS embed.FS

func initWorkingDirectory() error {
	if err := os.MkdirAll(filepath.Join(GetHybrDirectory(), "services"), 0755); err != nil {
		return err
	}

	return nil
}

func clearAndCopyDefaults() {
	servicesPath := filepath.Join(GetHybrDirectory(), "services")

	entries, err := os.ReadDir(servicesPath)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			entryPath := filepath.Join(servicesPath, entry.Name())
			if err := os.RemoveAll(entryPath); err != nil {
				panic(err)
			}
		}
	}

	err = fs.WalkDir(templatesFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == "templates" {
			return nil
		}

		relativePath := strings.TrimPrefix(path, "templates")
		destPath := filepath.Join(servicesPath, relativePath)

		if d.IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return err
			}
			return nil
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		sourceFile, err := templatesFS.Open(path)
		if err != nil {
			return fmt.Errorf("Failed opening file %s: %w", path, err)
		}
		defer sourceFile.Close()

		_, err = io.Copy(destFile, sourceFile)
		if err != nil {
			return fmt.Errorf("Failed copying file %s to %s: %w", path, destPath, err)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func initializeServices() []*serviceImpl {
	var services []*serviceImpl
	servicesPath := filepath.Join(GetHybrDirectory(), "services")

	entries, err := os.ReadDir(servicesPath)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		serviceFilePath := filepath.Join(servicesPath, entry.Name(), "service.json")

		if err := ValidateServiceJSON(serviceFilePath); err != nil {
			if shouldContinue := ConfirmInvalidService(err); shouldContinue {
				continue
			} else {
				os.Exit(1)
			}
		}

		data, err := os.ReadFile(serviceFilePath)

		if os.IsNotExist(err) {
			fmt.Printf("Found service without service.json: %s", entry.Name())
			continue
		}

		if err != nil {
			panic(err)
		}

		var service serviceImpl
		if err := json.Unmarshal(data, &service); err != nil {
			panic(fmt.Errorf("Data: %s\nErr: %w\nPath: %s", string(data), err, serviceFilePath))
		}

		services = append(services, &service)
	}

	return services
}
