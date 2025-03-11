package services

import (
	"fmt"
	"hybr/internal/docker"
	"hybr/internal/system"
	"hybr/internal/tailscale"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const (
	HybrDir = "/var/lib/hybr"
)

func GetHybrDirectory() string {
	if path := os.Getenv("HYBR_DIR"); path != "" {
		return path
	}
	return HybrDir
}

func cleanWorkingDirectory() error {
	if err := os.RemoveAll(GetHybrDirectory()); err != nil {
		return err
	}

	return nil
}

func initWorkingDirectory() error {
	if err := os.MkdirAll(filepath.Join(GetHybrDirectory(), "services"), 0755); err != nil {
		return err
	}

	return nil
}

func buildTemplateData(vars []*VariableDefinition, service HybrService) map[string]any {
	dns := tailscale.GetDNSName()
	data := make(map[string]any)

	for _, v := range vars {
		data[v.Name] = v.Value
	}

	data["Extras"] = map[string]string{
		"TS_DNS_NAME":  dns,
		"SERVICE_NAME": service.GetName(),
	}

	return data
}

func InstallServices(selected []HybrService) (err error) {
	ir := GetRegistry()

	for _, service := range selected {
		serviceDir := filepath.Join(GetHybrDirectory(), "services", service.GetName())

		var port string
		port, err = installTemplates(service, service.GetName())
		if err != nil {
			return err
		}

		installation := &serviceImpl{
			Port:        port,
			InstallDate: time.Now(),
			Status:      "installing",
			Name:        service.GetName(),
			Components:  []*docker.Component{},
			Variables:   make(map[string][]*VariableDefinition),
		}

		for fileName, vars := range service.GetVariables() {
			installation.Variables[strings.TrimSuffix(fileName, ".templ")] = vars
		}

		comps, err := docker.DetectComponents(serviceDir)
		if err != nil {
			return fmt.Errorf("Failed to detect docker components: %w", err)
		}
		installation.Components = comps

		if err := runInstallCommand(service.GetName(), installation); err != nil {
			installation.Status = "failed"
			ir.AddInstallation(installation)
			return err
		}

		magicDNSURL, err := tailscale.AddServeTunnel(
			service.GetIsRoot(),
			installation.GetName(),
			installation.GetPort(),
			service.GetTailscaleProxy(),
		)

		if err != nil {
			fmt.Printf("tailscale ERROR: %v\n", err)
			os.Exit(1)
		}

		installation.Status = "running"
		installation.URL = magicDNSURL
		if err := ir.AddInstallation(installation); err != nil {
			return err
		}

		fmt.Printf("[%s] Is Running at %s\n", service.GetName(), magicDNSURL)
	}

	return
}

func RestartService(serviceName string) error {
	if err := StopService(serviceName); err != nil {
		return err
	}

	if err := StartService(serviceName); err != nil {
		return err
	}

	return nil
}

func StartService(serviceName string) error {
	serviceDir := filepath.Join(GetHybrDirectory(), "services", serviceName)
	cmd := exec.Command("sh", "-c", "docker compose up -d")
	cmd.Dir = serviceDir
	if err := system.PipeCmdToStdout(cmd, "docker"); err != nil {
		return err
	}

	return nil
}

func StopService(serviceName string) error {
	serviceDir := filepath.Join(GetHybrDirectory(), "services", serviceName)
	cmd := exec.Command("sh", "-c", "docker compose down -v")
	cmd.Dir = serviceDir
	if err := system.PipeCmdToStdout(cmd, "docker"); err != nil {
		return err
	}
	return nil
}

func installTemplates(service HybrService, serviceName string) (string, error) {
	serviceDir := filepath.Join(GetHybrDirectory(), "services", serviceName)
	servicePath := filepath.Join(serviceDir, "templates")
	var port string
	err := filepath.Walk(servicePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".templ" {
			return nil
		}

		filename := filepath.Base(path)
		varDef, exists := service.GetVariables()[filename]

		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return fmt.Errorf("Unable to parse template %s: %w", filename, err)
		}

		outputPath := filepath.Join(serviceDir, strings.TrimSuffix(filename, ".templ"))
		out, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer out.Close()

		if !exists {
			err = tmpl.Execute(out, nil)
		} else {
			err = tmpl.Execute(out, buildTemplateData(varDef, service))
			if port == "" {
				port = findPort(varDef)
			}
		}

		return nil
	})

	return port, err
}

func reinstallTemplates(service HybrService, serviceName string) error {
	serviceDir := filepath.Join(GetHybrDirectory(), "services", serviceName)
	servicePath := filepath.Join(serviceDir, "templates")
	return filepath.Walk(servicePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".templ" {
			return nil
		}

		filename := strings.TrimSuffix(filepath.Base(path), ".templ")
		varDef, exists := service.GetVariables()[filename]

		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return fmt.Errorf("Unable to parse template %s: %w", filename, err)
		}

		outputPath := filepath.Join(serviceDir, filename)
		out, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("Unable to create file: %w", err)
		}
		defer out.Close()

		if !exists {
			err = tmpl.Execute(out, nil)
		} else {
			err = tmpl.Execute(out, buildTemplateData(varDef, service))
		}

		return nil
	})
}

func runInstallCommand(serviceName string, installation *serviceImpl) error {
	serviceDir := filepath.Join(GetHybrDirectory(), "services", serviceName)
	cmd := exec.Command("sh", "-c", "docker compose up -d")
	cmd.Dir = serviceDir

	if err := system.PipeCmdToStdout(cmd, serviceName); err != nil {
		return fmt.Errorf("Unable to install %s Service", serviceName)
	}

	for i := range installation.Components {
		installation.Components[i].Status = "running"
	}

	return nil
}
