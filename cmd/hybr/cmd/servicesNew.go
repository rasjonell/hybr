package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

func init() {
	servicesCmd.AddCommand(servicesNewCmd)
}

var servicesNewCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Creates a new hybr service directory",
	Long: `Creates a new directory with the given [name]
Adds default service definiton.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires one argument: [name]")
		}

		return nil
	},
	Run: newService,
}

func newService(cmd *cobra.Command, args []string) {
	newServiceName := args[0]
	newServicePath := filepath.Join(".", newServiceName)
	templatesPath := filepath.Join(newServicePath, "templates")

	if _, err := os.Stat(newServicePath); err == nil {
		fmt.Printf("Error: Directory '%s' already exists.\nTry a different name?\n\n", newServicePath)
		os.Exit(1)
	} else if !os.IsNotExist(err) {
		fmt.Printf("Error checking directory '%s': %v\n", newServicePath, err)
		os.Exit(1)
	}

	if err := os.MkdirAll(templatesPath, 0755); err != nil {
		fmt.Printf("Failed to create a directory %s: %v\n", templatesPath, err)
		os.Exit(1)
	}

	if err := createServiceJson(newServicePath, newServiceName); err != nil {
		fmt.Printf("Failed to create a service.json: %v\n", err)
		os.Exit(1)
	}

	if err := createEnvTempl(templatesPath); err != nil {
		fmt.Printf("Failed to create a .env.templ: %v\n", err)
		os.Exit(1)
	}

	if err := createDockerComposeTempl(templatesPath); err != nil {
		fmt.Printf("Failed to create a docker-compose.yml.templ: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nService Template For %s Created!\n\n", newServiceName)
	fmt.Println("Next Steps:")
	fmt.Printf(
		"\n1.  `cd %s`\n",
		newServiceName,
	)

	fmt.Println("\n2.  Modify your service")
	fmt.Println("    * Modify `service.json` for service definition and")
	fmt.Println("    * Modify `templates/.env.templ` and `templates/docker-compose.yml.templ`")
	fmt.Println("      to customize your service.")

	fmt.Println("\n3.  After finishing up, you can:")
	fmt.Println("    * Push to GitHub and install with `hybr install <github_url>`")
	fmt.Printf("    * Install locally using `hybr install ./%s`\n", newServiceName)
	fmt.Println(
		"\nFor a complete custom service creation guide, visit: " +
			"https://hybr.dev/docs/service",
	)
}

func createServiceJson(newServicePath, newServiceName string) error {
	templateContent := `{
  "name": "{{ .newServiceName }}",
  "description": "{{ .newServiceName }} - Description of your service",
  "hybrProxy": "/{{ .newServiceName }}",
  "tailscaleProxy": "/",
  "variables": {
    ".env.templ": [
      {
        "name": "SERVICE_VARIABLE",
        "default": "default-value",
        "description": "SERVICE_VARIABLE Description"
      }
    ]
  },
  "templates": [
    "docker-compose.yml.templ",
    ".env.templ"
  ]
}
`
	filePath := filepath.Join(newServicePath, "service.json")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.New("serviceJSON").Parse(templateContent)
	if err != nil {
		return err
	}

	data := map[string]string{
		"newServiceName": newServiceName,
	}

	return tmpl.Execute(file, data)
}

func createEnvTempl(templatesPath string) error {
	content := `# Your service.json "variables" definition will be populated here
SERVICE_VARIABLE={{ .SERVICE_VARIABLE }}

# Also we provide some extra variables that you can use here
EXTRA_VARIABLE={{ .Extras.TS_DNS_NAME }}`

	filePath := filepath.Join(templatesPath, ".env.templ")
	return os.WriteFile(filePath, []byte(content), 0644)
}

func createDockerComposeTempl(templatesPath string) error {
	content := `services:
  service_1:
    image: service_1_image
    restart: always
    volumes:
      - service_1:/var/lib/service_1
    env_file:
      - .env
    environment:
      - SERVICE_VARIABLE=${SERVICE_VARIABLE}

  service_2:
    image: service_2
    restart: always

volumes:
  service_1:
  service_2:`

	filePath := filepath.Join(templatesPath, "docker-compose.yml.templ")
	return os.WriteFile(filePath, []byte(content), 0644)
}
