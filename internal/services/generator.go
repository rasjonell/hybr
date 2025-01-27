package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	HybrDir = ".hybr"
)

func getWorkingDirectory() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, HybrDir)
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

func InstallServices(selected []*SelectedServiceModel) (err error) {
	for _, service := range selected {
		servicePath := filepath.Join(getWorkingDirectory(), "services", service.ServiceName, "templates")

		err = filepath.Walk(servicePath, func(path string, info os.FileInfo, err error) error {
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

			outputPath := filepath.Join(filepath.Dir(path), "..", strings.TrimSuffix(filename, ".templ"))
			out, err := os.Create(outputPath)
			if err != nil {
				return err
			}
			defer out.Close()

			if !exists {
				err = tmpl.Execute(out, nil)
			} else {
				err = tmpl.Execute(out, buildTemplateData(varDef))
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return
}
