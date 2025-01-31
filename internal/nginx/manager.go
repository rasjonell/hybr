package nginx

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type NginxServiceConfig struct {
	Name   string
	Port   string
	Domain string
}

type BaseConfig struct {
	Email  string
	Domain string
}

const (
	ConfDir  = "/etc/nginx"
	ConfDDir = "/etc/nginx/conf.d"
)

var baseConfig *BaseConfig

func Init(bc *BaseConfig, forceDefault bool) error {
	baseConfig = bc

	dirs := []string{ConfDir, ConfDDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("Failed to create directory %s: %w", dir, err)
		}
	}

	if _, err := os.Stat("/etc/nginx/nginx.conf"); os.IsNotExist(err) || forceDefault {
		if err := generateMainConfig(); err != nil {
			return fmt.Errorf("Failed to generate main config: %w", err)
		}
	}

	if err := VerifyConfig(); err != nil {
		return fmt.Errorf("Invalid Nginx configutaration: %w", err)
	}

	if err := StartService(); err != nil {
		return fmt.Errorf("Failed to start Nginx service: %w", err)
	}

	return nil
}

func AddSevice(name, port string) error {
	if port == "" {
		return fmt.Errorf("Failed to find PORT for service: %s", name)
	}

	config := NginxServiceConfig{
		Name:   name,
		Port:   port,
		Domain: baseConfig.Domain,
	}

	if err := generateServiceConfig(config); err != nil {
		return fmt.Errorf("Failed to generate service config: %w", err)
	}

	return Reload()
}

func RemoveService(name string) error {
	configPath := filepath.Join(ConfDir, fmt.Sprintf("%s.conf", name))
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Failed to remove service config: %w", err)
	}

	return nil
}

func VerifyConfig() error {
	if err := exec.Command("nginx", "-t").Run(); err != nil {
		return err
	}

	return nil
}

func StartService() error {
	if err := exec.Command("pgrep", "nginx").Run(); err != nil {
		if err := exec.Command("nginx").Run(); err != nil {
			return err
		}
	}

	return nil
}

func Reload() error {
	if err := exec.Command("nginx", "-s", "reload").Run(); err != nil {
		return fmt.Errorf("Failed to reload Nginx service: %w", err)
	}

	return nil
}
