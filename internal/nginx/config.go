package nginx

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed templates/*.templ
var templatesFS embed.FS

func generateMainConfig() error {
	configPath := filepath.Join(ConfDir, "nginx.conf")
	tmpl, err := parseTemplate("main", nil)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, tmpl, 0644)
}

func generateServiceConfig(config NginxServiceConfig, forceNoSSL bool) error {
	tmplName := "service-ssl"
	if forceNoSSL {
		tmplName = "service"
	}

	config.Domain = BuildServerName(config.SubDomain, config.Domain, config.Name)
	tmpl, err := parseTemplate(tmplName, config)
	if err != nil {
		return err
	}

	configPath := filepath.Join(ConfDir, "conf.d", fmt.Sprintf("%s.conf", config.Name))
	return os.WriteFile(configPath, tmpl, 0644)
}
