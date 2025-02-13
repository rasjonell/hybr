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
		data[v.Name] = v.Value
	}

	return data
}

func InstallServices(selected []HybrService, bc nginx.NginxConfig) (err error) {
	ir := GetRegistry()

	for _, service := range selected {
		serviceDir := filepath.Join(getWorkingDirectory(), "services", service.GetName())

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
			URL:         nginx.BuildServerName(service.IsSubDomain(), bc.GetDomain(), service.GetName()),
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

		if err := nginx.AddSevice(service.GetName(), port, service.IsSubDomain()); err != nil {
			installation.Status = "failed"
			ir.AddInstallation(installation)
			return err
		}

		installation.Status = "running"
		if err := ir.AddInstallation(installation); err != nil {
			return err
		}

		fmt.Printf("\n[%s] Is Installed and Nginx Is Configured\n", service.GetName())
		fmt.Printf("[%s] Is Running at %s\n", service.GetName(), installation.URL)
	}

	return
}

func RestartService(serviceName string) error {
	serviceDir := filepath.Join(getWorkingDirectory(), "services", serviceName)
	cmd := exec.Command("sh", "-c", "docker compose down -v")
	cmd.Dir = serviceDir
	if err := nginx.PipeCmdToStdout(cmd, "docker"); err != nil {
		return err
	}

	cmd = exec.Command("sh", "-c", "docker compose up -d")
	cmd.Dir = serviceDir
	if err := nginx.PipeCmdToStdout(cmd, "docker"); err != nil {
		return err
	}

	return nil
}

func installTemplates(service HybrService, serviceName string) (string, error) {
	serviceDir := filepath.Join(getWorkingDirectory(), "services", serviceName)
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
			err = tmpl.Execute(out, buildTemplateData(varDef))
			if port == "" {
				port = findPort(varDef)
			}
		}

		return nil
	})

	return port, err
}

func reinstallTemplates(service HybrService, serviceName string) error {
	serviceDir := filepath.Join(getWorkingDirectory(), "services", serviceName)
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
			err = tmpl.Execute(out, buildTemplateData(varDef))
		}

		return nil
	})
}

func runInstallCommand(serviceName string, installation *serviceImpl) error {
	serviceDir := filepath.Join(getWorkingDirectory(), "services", serviceName)
	cmd := exec.Command("sh", "-c", "docker compose up -d")
	cmd.Dir = serviceDir

	if err := nginx.PipeCmdToStdout(cmd, serviceName); err != nil {
		return fmt.Errorf("Unable to install %s Service", serviceName)
	}

	for i := range installation.Components {
		installation.Components[i].Status = "running"
	}

	return nil
}
