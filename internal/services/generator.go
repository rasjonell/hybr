package services

import (
	"fmt"
	"hybr/internal/docker"
	"hybr/internal/nginx"
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

func getWorkingDirectory() string {
	if path := os.Getenv("HYBR_DIR"); path != "" {
		return path
	}
	return HybrDir
}

func cleanWorkingDirectory() error {
	if err := os.RemoveAll(getWorkingDirectory()); err != nil {
		return err
	}

	return nil
}

func initWorkingDirectory() error {
	if err := os.MkdirAll(filepath.Join(getWorkingDirectory(), "services"), 0755); err != nil {
		return err
	}

	return nil
}

func buildTemplateData(vars []*VariableDefinition) map[string]string {
	data := make(map[string]string)

	for _, v := range vars {
		data[v.Key] = v.Value
	}

	return data
}

func InstallServices(selected []*SelectedServiceModel, bc *nginx.BaseConfig) (err error) {
	ir := GetRegistry()

	for _, service := range selected {
		serviceDir := filepath.Join(getWorkingDirectory(), "services", service.ServiceName)
		servicePath := filepath.Join(serviceDir, "templates")

		var port string
		port, err = installTemplates(service, servicePath, serviceDir)
		if err != nil {
			return err
		}

		installation := &ServiceInstallation{
			Port:        port,
			InstallDate: time.Now(),
			Status:      "installing",
			Name:        service.ServiceName,
			Components:  []docker.Component{},
			Variables:   make(map[string]string),
			URL:         nginx.BuildServerName(service.SubDomain, bc.Domain, service.ServiceName),
		}

		for _, vars := range service.Variables {
			for _, v := range vars {
				installation.Variables[v.Key] = v.Value
			}
		}

		comps, err := docker.DetectComponents(serviceDir)
		if err != nil {
			return fmt.Errorf("Failed to detect docker components: %w", err)
		}
		installation.Components = comps

		if err := runInstallCommand(service, installation, serviceDir); err != nil {
			installation.Status = "failed"
			ir.AddInstallation(installation)
			return err
		}

		if err := nginx.AddSevice(service.ServiceName, port, service.SubDomain); err != nil {
			installation.Status = "failed"
			ir.AddInstallation(installation)
			return err
		}

		installation.Status = "running"
		if err := ir.AddInstallation(installation); err != nil {
			return err
		}

		fmt.Printf("\n[%s] Is Installed and Nginx Is Configured\n", service.ServiceName)
		fmt.Printf("[%s] Is Running at %s\n", service.ServiceName, installation.URL)
	}

	return
}

func installTemplates(service *SelectedServiceModel, servicePath, serviceDir string) (string, error) {
	var port string
	err := filepath.Walk(servicePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".templ" {
			return nil
		}

		filename := filepath.Base(path)
		varDef, exists := service.Variables[filename]

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
			err = tmpl.Execute(out, buildTemplateData(varDef))
			if port == "" {
				port = findPort(varDef)
			}
		}

		return nil
	})

	return port, err
}

func runInstallCommand(service *SelectedServiceModel, installation *ServiceInstallation, serviceDir string) error {
	if service.InstallCommand == "" {
		return fmt.Errorf("Service %s doesn't have an `installCommand`", service.ServiceName)
	}

	cmd := exec.Command("sh", "-c", service.InstallCommand)
	cmd.Dir = serviceDir

	if err := nginx.PipeCmdToStdout(cmd, service.ServiceName); err != nil {
		return fmt.Errorf("Unable to install %s Service", service.ServiceName)
	}

	for i := range installation.Components {
		installation.Components[i].Status = "running"
	}

	return nil
}
