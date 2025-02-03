package nginx

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	CertbotCmd = "certbot"
	SSLDir     = "/etc/letsencrypt/live"
)

func ObtainSSLCert(config NginxServiceConfig, baseConfig *BaseConfig) error {
	if err := checkCertbot(); err != nil {
		return err
	}

	domain := BuildServerName(config.SubDomain, config.Domain, config.Name)
	certPath := filepath.Join(SSLDir, domain)
	if _, err := os.Stat(certPath); err == nil {
		fmt.Printf("Certificate for %s already exists\n", domain)
		return nil
	}

	fmt.Printf("Obtaining SSL certificate for %s...\n", domain)

	args := []string{
		"certonly",
		"--standalone",
		"--non-interactive",
		"--agree-tos",
		"--email", baseConfig.Email,
		"-d", domain,
		"--redirect",
	}

	cmd := exec.Command(CertbotCmd, args...)
	if err := PipeCmdToStdout(cmd, "certbot"); err != nil {
		return err
	}

	fmt.Printf("[certbot] Successfully obtained SSL certificate for %s\n", domain)
	return nil
}

func SetupAutoRenewal() error {
	if err := checkCertbot(); err != nil {
		return err
	}

	cmd := exec.Command("crontab", "-l")
	currentCrontabBytes, err := cmd.Output()
	if err != nil && cmd.ProcessState.ExitCode() != 1 {
		return fmt.Errorf("Failed to read current crontab: %w", err)
	}
	currentCrontab := string(currentCrontabBytes)

	renewalJob := "0 0,12 * * * /usr/bin/certbot renew --quiet --deploy-hook \"nginx -s reload\""
	if currentCrontab != "" && strings.Contains(currentCrontab, "certbot renew") {
		fmt.Println("Certificate auto-renewal already configured")
		return nil
	}

	if currentCrontab != "" && !strings.HasSuffix(currentCrontab, "\n") {
		currentCrontab += "\n"
	}
	currentCrontab += renewalJob + "\n"

	cmd = exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(currentCrontab)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to update crontab: %w", err)
	}

	fmt.Println("Successfully configured automatic certificate renewal")
	return nil
}

func CheckCertificateStatus(domain string) error {
	if err := checkCertbot(); err != nil {
		return err
	}

	cmd := exec.Command(CertbotCmd, "certificates")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to check certificate status: %w", err)
	}

	if !strings.Contains(string(output), domain) {
		return fmt.Errorf("No certificate found for domain %s", domain)
	}

	fmt.Printf("Certificate status for %s:\n%s\n", domain, string(output))
	return nil
}

func ForceRenewCertificate(domain string) error {
	if err := checkCertbot(); err != nil {
		return err
	}

	cmd := exec.Command(CertbotCmd, "renew", "--force-renewal", "-d", domain)
	if err := PipeCmdToStdout(cmd, "certbot"); err != nil {
		return err
	}

	if err := Reload(); err != nil {
		return err
	}

	return nil
}

func RevokeCertificate(domain string) error {
	if err := checkCertbot(); err != nil {
		return err
	}

	cmd := exec.Command(CertbotCmd, "revoke", "--cert-name", domain, "--non-interactive")
	if err := PipeCmdToStdout(cmd, "certbot"); err != nil {
		return err
	}

	cmd = exec.Command(CertbotCmd, "delete", "--cert-name", domain, "--non-interactive")
	if err := PipeCmdToStdout(cmd, "certbot"); err != nil {
		return err
	}

	return nil
}

func checkCertbot() error {
	_, err := exec.LookPath(CertbotCmd)
	if err != nil {
		return fmt.Errorf("Certbot is not installed\nPlease install certbot to continue.")
	}
	return nil
}
