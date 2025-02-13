package nginx

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type NginxConfig interface {
	GetEmail() string
	GetDomain() string
}

type NginxServiceConfig struct {
	SubDomain bool
	Name      string
	Port      string
	Domain    string
}

type ConfigImpl struct {
	Email  string
	Domain string
}

func (b *ConfigImpl) GetEmail() string {
	return b.Email
}

func (b *ConfigImpl) GetDomain() string {
	return b.Domain
}

const (
	ConfDir  = "/etc/nginx"
	ConfDDir = "/etc/nginx/conf.d"
)

var baseConfig *ConfigImpl
var forceNoSSL bool

func Init(bc NginxConfig, forceDefault, noSSL bool) error {
	var ok bool
	forceNoSSL = noSSL
	baseConfig, ok = bc.(*ConfigImpl)
	if !ok {
		return fmt.Errorf("Invalid Nginx Config type")
	}

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

func AddSevice(name, port string, subDomain bool) error {
	if port == "" {
		return fmt.Errorf("Failed to find PORT for service: %s", name)
	}

	config := NginxServiceConfig{
		Name:      name,
		Port:      port,
		SubDomain: subDomain,
		Domain:    baseConfig.Domain,
	}

	if err := generateServiceConfig(config, forceNoSSL); err != nil {
		return fmt.Errorf("Failed to generate service config: %w", err)
	}

	if !forceNoSSL {
		if err := ObtainSSLCert(config, baseConfig); err != nil {
			return err
		}
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
