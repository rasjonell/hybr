package nginx

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

func parseTemplate(name string, data any) ([]byte, error) {
	content, err := templatesFS.ReadFile("templates/" + name + ".conf.templ")
	if err != nil {
		return nil, fmt.Errorf("Failed to read template %s: %w", name, err)
	}

	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("Failed to execure template %s: %w", name, err)
	}

	return buf.Bytes(), nil
}

func PipeCmdToStdout(cmd *exec.Cmd, label string) error {
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Error creating stdout pipe: %v", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Error creating stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Failed to start command for %s: %v", label, err)
	}

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			fmt.Printf("[%s] %s\n", label, scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			fmt.Fprintf(os.Stderr, "[%s] %s\n", label, scanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed for %s: %v", label, err)
	}

	return nil
}

func BuildServerName(subDomain bool, domain, name string) string {
	if subDomain {
		return name + "." + domain
	}

	return domain
}
