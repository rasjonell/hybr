package docker

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Component struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Version string `json:"version"`
}

func DetectComponents(serviceDir string) ([]*Component, error) {
	composer := filepath.Join(serviceDir, "docker-compose.yml")
	data, err := os.ReadFile(composer)
	if err != nil {
		return nil, err
	}

	var compose struct {
		Service map[string]struct {
			Image string `yaml:"image"`
		} `yaml:"services"`
	}

	if err := yaml.Unmarshal(data, &compose); err != nil {
		return nil, err
	}

	components := []*Component{}
	for name, service := range compose.Service {
		components = append(components, &Component{
			Name:    name,
			Status:  "installing",
			Version: parseImageVersion(service.Image),
		})
	}

	return components, nil
}

func parseImageVersion(image string) string {
	parts := strings.Split(image, ":")
	if len(parts) > 1 {
		return parts[1]
	}

	return "latest"
}
