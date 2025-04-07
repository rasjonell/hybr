package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rasjonell/hybr/internal/docker"
)

type InstallationRegistry struct {
	stateFile     string
	mu            sync.RWMutex
	installations map[string]*serviceImpl
}

var (
	once                 sync.Once
	installationRegistry *InstallationRegistry
)

func InitRegistry(forceResetTemplates bool) {
	initWorkingDirectory()
	if forceResetTemplates {
		clearAndCopyDefaults()
	}

	services := initializeServices()

	for _, service := range services {
		register(service)
	}
}

func GetRegistry() *InstallationRegistry {
	once.Do(func() {
		installationRegistry = &InstallationRegistry{
			installations: make(map[string]*serviceImpl),
			stateFile:     filepath.Join(GetHybrDirectory(), "installations.json"),
		}
		installationRegistry.load()
	})

	return installationRegistry
}

func (r *InstallationRegistry) load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.stateFile)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("Failed to read installation registry: %w", err)
	}

	return json.Unmarshal(data, &r.installations)
}

func (r *InstallationRegistry) save() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := json.MarshalIndent(r.installations, "", " ")
	if err != nil {
		return fmt.Errorf("Failed to marshal installation registry: %w", err)
	}

	return os.WriteFile(r.stateFile, data, 0644)
}

func (r *InstallationRegistry) RegisterServiceEvents() {
	for name := range installationRegistry.installations {
		RegisterEventSources(name)
	}
}

func (r *InstallationRegistry) AddInstallation(service *serviceImpl) error {
	r.mu.Lock()
	if service.GetStatus() == "running" {
		service.LastStartTime = time.Now()
	}
	r.installations[service.Name] = service
	RegisterEventSources(service.Name)
	r.mu.Unlock()

	return r.save()
}

func (r *InstallationRegistry) GetInstallation(name string) (HybrService, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	service, exists := r.installations[name]
	return service, exists
}

func (r *InstallationRegistry) UpdateStatus(name, status string) error {
	r.mu.Lock()
	if service, exists := r.installations[name]; exists {
		service.Status = status
		if status == "running" {
			service.LastStartTime = time.Now()
		}
	}
	r.mu.Unlock()

	return r.save()
}

func (r *InstallationRegistry) UpdateComponent(name string, component *docker.Component) error {
	r.mu.Lock()
	if service, exists := r.installations[name]; exists {
		for _, comp := range service.Components {
			if comp.Name == component.Name {
				comp = component
				break
			}
		}
	}
	r.mu.Unlock()

	return r.save()
}

func (r *InstallationRegistry) ListInstallations() []HybrService {
	r.mu.RLock()
	defer r.mu.RUnlock()

	services := make([]HybrService, 0, len(r.installations))
	for _, service := range r.installations {
		services = append(services, service)
	}
	return services
}

func (r *InstallationRegistry) RemoveInstalltion(name string) error {
	r.mu.Lock()
	delete(r.installations, name)
	r.mu.Unlock()

	return r.save()
}

func ListInstalledServiceNames() []string {
	r := GetRegistry()
	r.mu.RLock()
	defer r.mu.RUnlock()

	services := make([]string, 0, len(r.installations))
	for _, service := range r.installations {
		services = append(services, service.Name)
	}
	return services
}
