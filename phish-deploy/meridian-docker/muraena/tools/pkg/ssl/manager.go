package ssl

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Manager manages SSL certificates
type Manager struct {
	certbotPath string
}

// NewManager creates a new SSL manager
func NewManager() *Manager {
	return &Manager{
		certbotPath: "certbot",
	}
}

// GenerateCertificate generates a new SSL certificate using Certbot
func (m *Manager) GenerateCertificate(ctx context.Context, config CertbotConfig) error {
	args := []string{"certonly"}

	if config.Standalone {
		args = append(args, "--standalone")
	} else if config.Webroot != "" {
		args = append(args, "--webroot", "-w", config.Webroot)
	}

	args = append(args, "-d", config.Domain)

	if config.Email != "" {
		args = append(args, "--email", config.Email)
	}

	if config.AgreeToS {
		args = append(args, "--agree-tos")
	}

	if config.NonInteractive {
		args = append(args, "--non-interactive")
	}

	if config.ForceRenewal {
		args = append(args, "--force-renewal")
	}

	if config.TestCert {
		args = append(args, "--test-cert")
	}

	cmd := exec.CommandContext(ctx, "sudo", append([]string{m.certbotPath}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("certbot failed: %w", err)
	}

	return nil
}

// RenewCertificate renews an existing certificate
func (m *Manager) RenewCertificate(ctx context.Context, domain string) error {
	cmd := exec.CommandContext(ctx, "sudo", m.certbotPath, "renew", "--cert-name", domain)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("certificate renewal failed: %w", err)
	}

	return nil
}

// GetCertificateInfo retrieves information about a certificate
func (m *Manager) GetCertificateInfo(domain string) (*CertificateInfo, error) {
	certPath := filepath.Join("/etc/letsencrypt/live", domain, "fullchain.pem")
	keyPath := filepath.Join("/etc/letsencrypt/live", domain, "privkey.pem")

	// Read certificate file
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate: %w", err)
	}

	// Parse PEM block
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}

	// Parse certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Calculate days left
	now := time.Now()
	daysLeft := int(cert.NotAfter.Sub(now).Hours() / 24)
	isValid := now.After(cert.NotBefore) && now.Before(cert.NotAfter)

	info := &CertificateInfo{
		Domain:    domain,
		Issuer:    cert.Issuer.String(),
		Subject:   cert.Subject.String(),
		NotBefore: cert.NotBefore,
		NotAfter:  cert.NotAfter,
		IsValid:   isValid,
		DaysLeft:  daysLeft,
		CertPath:  certPath,
		KeyPath:   keyPath,
	}

	return info, nil
}

// ValidateCertificate validates a certificate
func (m *Manager) ValidateCertificate(domain string) error {
	info, err := m.GetCertificateInfo(domain)
	if err != nil {
		return err
	}

	if !info.IsValid {
		return fmt.Errorf("certificate is not valid (expired or not yet valid)")
	}

	if info.DaysLeft < 0 {
		return fmt.Errorf("certificate has expired")
	}

	if info.DaysLeft < 30 {
		return fmt.Errorf("certificate expires soon (%d days left)", info.DaysLeft)
	}

	return nil
}

// SetupAutoRenewal sets up automatic certificate renewal
func (m *Manager) SetupAutoRenewal(config RenewalConfig) error {
	// Create systemd timer or cron job for auto-renewal
	cronEntry := "0 0,12 * * * root certbot renew --quiet"

	if config.PostHook != "" {
		cronEntry += " --post-hook '" + config.PostHook + "'"
	}

	// Write to cron.d
	cronFile := "/etc/cron.d/certbot-renewal"
	if err := os.WriteFile(cronFile, []byte(cronEntry+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to create cron job: %w", err)
	}

	return nil
}

// ListCertificates lists all installed certificates
func (m *Manager) ListCertificates() ([]string, error) {
	liveDir := "/etc/letsencrypt/live"

	entries, err := os.ReadDir(liveDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificates directory: %w", err)
	}

	var domains []string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "README" {
			domains = append(domains, entry.Name())
		}
	}

	return domains, nil
}

// DeleteCertificate deletes a certificate
func (m *Manager) DeleteCertificate(ctx context.Context, domain string) error {
	cmd := exec.CommandContext(ctx, "sudo", m.certbotPath, "delete", "--cert-name", domain)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete certificate: %w", err)
	}

	return nil
}
