package ssl

import "time"

// CertificateInfo represents SSL certificate information
type CertificateInfo struct {
	Domain      string    `json:"domain"`
	Issuer      string    `json:"issuer"`
	Subject     string    `json:"subject"`
	NotBefore   time.Time `json:"not_before"`
	NotAfter    time.Time `json:"not_after"`
	IsValid     bool      `json:"is_valid"`
	DaysLeft    int       `json:"days_left"`
	CertPath    string    `json:"cert_path"`
	KeyPath     string    `json:"key_path"`
	Fingerprint string    `json:"fingerprint,omitempty"`
}

// CertbotConfig represents Certbot configuration
type CertbotConfig struct {
	Email          string
	Domain         string
	Webroot        string
	Standalone     bool
	ForceRenewal   bool
	TestCert       bool
	AgreeToS       bool
	NonInteractive bool
}

// RenewalConfig represents auto-renewal configuration
type RenewalConfig struct {
	Enabled       bool
	CheckInterval time.Duration
	RenewBefore   time.Duration
	PostHook      string
}
