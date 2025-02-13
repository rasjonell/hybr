package services

import (
	"encoding/json"
	"fmt"
	"hybr/internal/docker"
	"os"
	"path/filepath"
	"sync"
	"time"
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
	if forceResetTemplates {
		cleanWorkingDirectory()
	}

	initWorkingDirectory()
	services := getServices()

	for _, service := range services {
		register(service)
	}
}

func GetRegistry() *InstallationRegistry {
	once.Do(func() {
		installationRegistry = &InstallationRegistry{
			installations: make(map[string]*serviceImpl),
			stateFile:     filepath.Join(getWorkingDirectory(), "installations.json"),
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

func (r *InstallationRegistry) AddInstallation(service *serviceImpl) error {
	r.mu.Lock()
	if service.GetStatus() == "running" {
		service.LastStartTime = time.Now()
	}
	r.installations[service.Name] = service
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

func (r *InstallationRegistry) UpdateComponent(name string, component docker.Component) error {
	r.mu.Lock()
	if service, exists := r.installations[name]; exists {
		for i, comp := range service.Components {
			if comp.Name == name {
				service.Components[i] = comp
				break
			}
		}
	}
	r.mu.Unlock()

	return r.save()
}

func (r *InstallationRegistry) ListInstallations() []*serviceImpl {
	r.mu.RLock()
	defer r.mu.RUnlock()

	services := make([]*serviceImpl, 0, len(r.installations))
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
