package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// NecroBrowserService manages the NecroBrowser service
type NecroBrowserService struct {
	config     *ServiceConfig
	workDir    string
	binaryPath string
	configPath string
	logPath    string
	apiPort    int
}

// NewNecroBrowserService creates a new NecroBrowser service manager
func NewNecroBrowserService(config *ServiceConfig) *NecroBrowserService {
	workDir := config.NecroBrowserDir
	binaryPath := filepath.Join(workDir, config.NecroBrowserBinary)
	configPath := filepath.Join(workDir, config.NecroBrowserConfig)
	logDir := filepath.Join(workDir, "logs")
	logPath := filepath.Join(logDir, "necrobrowser_startup.log")

	return &NecroBrowserService{
		config:     config,
		workDir:    workDir,
		binaryPath: binaryPath,
		configPath: configPath,
		logPath:    logPath,
		apiPort:    3000,
	}
}

// GetName returns the service name
func (ns *NecroBrowserService) GetName() string {
	return "NecroBrowser"
}

// Start starts the NecroBrowser service
func (ns *NecroBrowserService) Start(ctx context.Context) error {
	// Check if already running
	if ns.isRunning() {
		return fmt.Errorf("necrobrowser is already running")
	}

	// Verify binary exists
	if _, err := os.Stat(ns.binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("necrobrowser binary not found: %s", ns.binaryPath)
	}

	// Verify config exists
	if _, err := os.Stat(ns.configPath); os.IsNotExist(err) {
		return fmt.Errorf("necrobrowser config not found: %s", ns.configPath)
	}

	// Create log directory if needed
	logDir := filepath.Dir(ns.logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Start in screen session
	// Use npm start or node necrobrowser.js depending on setup
	screenCmd := fmt.Sprintf("cd %s && npm start | tee %s", ns.workDir, ns.logPath)

	cmd := exec.CommandContext(ctx, "screen", "-dmS", "necro", "bash", "-c", screenCmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start necrobrowser: %w", err)
	}

	// Wait for service to start
	time.Sleep(3 * time.Second)

	// Verify it started
	if !ns.isRunning() {
		return fmt.Errorf("necrobrowser failed to start (screen session not found)")
	}

	// Wait for API to be ready
	if err := ns.waitForAPI(ctx, 30*time.Second); err != nil {
		return fmt.Errorf("necrobrowser API not ready: %w", err)
	}

	return nil
}

// Stop stops the NecroBrowser service
func (ns *NecroBrowserService) Stop(ctx context.Context) error {
	if !ns.isRunning() {
		return nil // Already stopped
	}

	cmd := exec.CommandContext(ctx, "screen", "-S", "necro", "-X", "quit")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop necrobrowser: %w", err)
	}

	// Wait for graceful shutdown
	time.Sleep(2 * time.Second)

	return nil
}

// Restart restarts the NecroBrowser service
func (ns *NecroBrowserService) Restart(ctx context.Context) error {
	if err := ns.Stop(ctx); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	return ns.Start(ctx)
}

// HealthCheck checks if NecroBrowser is healthy
func (ns *NecroBrowserService) HealthCheck(ctx context.Context) error {
	// Check if screen session exists
	if !ns.isRunning() {
		return fmt.Errorf("necrobrowser screen session not found")
	}

	// Check if API port is listening
	if !ns.isPortListening(ns.apiPort) {
		return fmt.Errorf("API port %d not listening", ns.apiPort)
	}

	// Try to hit the health endpoint
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://localhost:%d/health", ns.apiPort))
	if err != nil {
		return fmt.Errorf("API health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return nil
}

// GetStatus returns the current status
func (ns *NecroBrowserService) GetStatus(ctx context.Context) ServiceStatus {
	status := ServiceStatus{
		Name:  ns.GetName(),
		Ports: []int{ns.apiPort},
	}

	if err := ns.HealthCheck(ctx); err == nil {
		status.Running = true
		status.Healthy = true
		status.LastSeen = time.Now()

		// Try to get PID
		if pid := ns.getPID(); pid > 0 {
			status.PID = pid
		}
	} else {
		status.Running = ns.isRunning()
		status.Healthy = false
		status.Errors = append(status.Errors, err.Error())
	}

	return status
}

// isRunning checks if the screen session exists
func (ns *NecroBrowserService) isRunning() bool {
	cmd := exec.Command("screen", "-list")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), "necro")
}

// isPortListening checks if a port is listening
func (ns *NecroBrowserService) isPortListening(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// waitForAPI waits for the API to become available
func (ns *NecroBrowserService) waitForAPI(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 2 * time.Second}

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		resp, err := client.Get(fmt.Sprintf("http://localhost:%d/health", ns.apiPort))
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}

		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("timeout waiting for API to be ready")
}

// getPID attempts to get the process ID
func (ns *NecroBrowserService) getPID() int {
	cmd := exec.Command("pgrep", "-f", "necrobrowser")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	var pid int
	fmt.Sscanf(string(output), "%d", &pid)
	return pid
}

// GetLogs returns the last N lines of logs
func (ns *NecroBrowserService) GetLogs(lines int) (string, error) {
	cmd := exec.Command("tail", "-n", fmt.Sprintf("%d", lines), ns.logPath)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}
	return string(output), nil
}

// FollowLogs tails the log file in real-time
func (ns *NecroBrowserService) FollowLogs(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "tail", "-f", ns.logPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
