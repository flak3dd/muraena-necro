package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/muraenateam/muraena/tools/pkg/common"
	"github.com/muraenateam/muraena/tools/pkg/ssl"
	"github.com/spf13/cobra"
)

var (
	verbose    bool
	logFile    string
	domain     string
	email      string
	standalone bool
	testCert   bool
	postHook   string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ssl-manager",
		Short: "Muraena SSL Certificate Manager",
		Long: `Manage SSL certificates for Muraena phishing infrastructure.
Integrates with Let's Encrypt via Certbot for automatic certificate generation.`,
	}

	// Generate command
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate SSL certificate",
		Long:  "Generate a new SSL certificate using Let's Encrypt",
		RunE:  runGenerate,
	}
	generateCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain name")
	generateCmd.Flags().StringVarP(&email, "email", "e", "", "Email address for Let's Encrypt")
	generateCmd.Flags().BoolVar(&standalone, "standalone", true, "Use standalone mode")
	generateCmd.Flags().BoolVar(&testCert, "test", false, "Generate test certificate")
	generateCmd.MarkFlagRequired("domain")
	generateCmd.MarkFlagRequired("email")

	// Renew command
	renewCmd := &cobra.Command{
		Use:   "renew",
		Short: "Renew SSL certificate",
		Long:  "Renew an existing SSL certificate",
		RunE:  runRenew,
	}
	renewCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain name")
	renewCmd.MarkFlagRequired("domain")

	// Info command
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show certificate information",
		Long:  "Display detailed information about a certificate",
		RunE:  runInfo,
	}
	infoCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain name")
	infoCmd.MarkFlagRequired("domain")

	// Validate command
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate certificate",
		Long:  "Validate a certificate and check expiry",
		RunE:  runValidate,
	}
	validateCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain name")
	validateCmd.MarkFlagRequired("domain")

	// List command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List certificates",
		Long:  "List all installed certificates",
		RunE:  runList,
	}

	// Delete command
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete certificate",
		Long:  "Delete a certificate",
		RunE:  runDelete,
	}
	deleteCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain name")
	deleteCmd.MarkFlagRequired("domain")

	// Auto-renew command
	autoRenewCmd := &cobra.Command{
		Use:   "auto-renew",
		Short: "Setup auto-renewal",
		Long:  "Setup automatic certificate renewal",
		RunE:  runAutoRenew,
	}
	autoRenewCmd.Flags().StringVar(&postHook, "post-hook", "", "Command to run after renewal")

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&logFile, "log", "l", "", "Log file path")

	// Add commands
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(renewCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(autoRenewCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runGenerate(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	logger.Info(fmt.Sprintf("Generating SSL certificate for: %s", domain))

	manager := ssl.NewManager()
	ctx := context.Background()

	config := ssl.CertbotConfig{
		Email:          email,
		Domain:         domain,
		Standalone:     standalone,
		TestCert:       testCert,
		AgreeToS:       true,
		NonInteractive: true,
	}

	if err := manager.GenerateCertificate(ctx, config); err != nil {
		return fmt.Errorf("failed to generate certificate: %w", err)
	}

	logger.Success(fmt.Sprintf("Certificate generated for: %s", domain))
	logger.Info("Certificate location: /etc/letsencrypt/live/" + domain + "/")

	return nil
}

func runRenew(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	logger.Info(fmt.Sprintf("Renewing certificate for: %s", domain))

	manager := ssl.NewManager()
	ctx := context.Background()

	if err := manager.RenewCertificate(ctx, domain); err != nil {
		return fmt.Errorf("failed to renew certificate: %w", err)
	}

	logger.Success(fmt.Sprintf("Certificate renewed for: %s", domain))

	return nil
}

func runInfo(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	logger.Info(fmt.Sprintf("Retrieving certificate information for: %s", domain))

	manager := ssl.NewManager()

	info, err := manager.GetCertificateInfo(domain)
	if err != nil {
		return fmt.Errorf("failed to get certificate info: %w", err)
	}

	fmt.Println("\n🔒 SSL Certificate Information:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Domain:       %s\n", info.Domain)
	fmt.Printf("Issuer:       %s\n", info.Issuer)
	fmt.Printf("Subject:      %s\n", info.Subject)
	fmt.Printf("Valid From:   %s\n", info.NotBefore.Format(time.RFC3339))
	fmt.Printf("Valid Until:  %s\n", info.NotAfter.Format(time.RFC3339))
	fmt.Printf("Days Left:    %d\n", info.DaysLeft)

	if info.IsValid {
		fmt.Printf("Status:       ✓ Valid\n")
	} else {
		fmt.Printf("Status:       ✗ Invalid/Expired\n")
	}

	fmt.Printf("\nCertificate:  %s\n", info.CertPath)
	fmt.Printf("Private Key:  %s\n", info.KeyPath)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}

func runValidate(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	logger.Info(fmt.Sprintf("Validating certificate for: %s", domain))

	manager := ssl.NewManager()

	if err := manager.ValidateCertificate(domain); err != nil {
		logger.Error(fmt.Sprintf("Validation failed: %v", err))
		return err
	}

	logger.Success(fmt.Sprintf("Certificate is valid for: %s", domain))

	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	logger.Info("Listing installed certificates...")

	manager := ssl.NewManager()

	domains, err := manager.ListCertificates()
	if err != nil {
		return fmt.Errorf("failed to list certificates: %w", err)
	}

	fmt.Println("\n📜 Installed Certificates:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, domain := range domains {
		info, err := manager.GetCertificateInfo(domain)
		if err != nil {
			fmt.Printf("[%d] %s (error reading certificate)\n", i+1, domain)
			continue
		}

		status := "✓ Valid"
		if !info.IsValid {
			status = "✗ Invalid"
		} else if info.DaysLeft < 30 {
			status = "⚠ Expires soon"
		}

		fmt.Printf("[%d] %s (%s, %d days left)\n", i+1, domain, status, info.DaysLeft)
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\nTotal: %d certificates\n", len(domains))

	return nil
}

func runDelete(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	logger.Warning(fmt.Sprintf("Deleting certificate for: %s", domain))

	fmt.Print("Are you sure? (yes/no): ")
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "yes" {
		logger.Info("Operation cancelled")
		return nil
	}

	manager := ssl.NewManager()
	ctx := context.Background()

	if err := manager.DeleteCertificate(ctx, domain); err != nil {
		return fmt.Errorf("failed to delete certificate: %w", err)
	}

	logger.Success(fmt.Sprintf("Certificate deleted for: %s", domain))

	return nil
}

func runAutoRenew(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	logger.Info("Setting up automatic certificate renewal...")

	manager := ssl.NewManager()

	config := ssl.RenewalConfig{
		Enabled:       true,
		CheckInterval: 12 * time.Hour,
		RenewBefore:   30 * 24 * time.Hour,
		PostHook:      postHook,
	}

	if err := manager.SetupAutoRenewal(config); err != nil {
		return fmt.Errorf("failed to setup auto-renewal: %w", err)
	}

	logger.Success("Auto-renewal configured")
	logger.Info("Certificates will be checked twice daily")
	logger.Info("Renewal will occur 30 days before expiry")

	return nil
}
