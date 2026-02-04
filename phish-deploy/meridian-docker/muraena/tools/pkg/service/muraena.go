package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// MuraenaService manages the Muraena reverse proxy service
type MuraenaService struct {
	config     *ServiceConfig
	workDir    string
	binaryPath string
	configPath string
	logPath    string
}

// NewMuraenaService creates a new Muraena service manager
func NewMuraenaService(config *ServiceConfig) *MuraenaService {
	workDir := config.MuraenaDir
	binaryPath := filepath.Join(workDir, config.MuraenaBinary)
	configPath := filepath.Join(workDir, config.MuraenaConfig)
	logPath := filepath.Join(workDir, "muraena.log")

	return &MuraenaService{
		config:     config,
		workDir:    workDir,
		binaryPath: binaryPath,
		configPath: configPath,
		logPath:    logPath,
	}
}

// GetName returns the service name
func (ms *MuraenaService) GetName() string {
	return "Muraena"
}

// Start starts the Muraena service
func (ms *MuraenaService) Start(ctx context.Context) error {
	// Check if already running
	if ms.isRunning() {
		return fmt.Errorf("muraena is already running")
	}

	// Verify binary exists
	if _, err := os.Stat(ms.binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("muraena binary not found: %s", ms.binaryPath)
	}

	// Verify config exists
	if _, err := os.Stat(ms.configPath); os.IsNotExist(err) {
		return fmt.Errorf("muraena config not found: %s", ms.configPath)
	}

	// Create log directory if needed
	logDir := filepath.Dir(ms.logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Start in screen session
	screenCmd := fmt.Sprintf("cd %s && ./muraena.bin -config config.toml | tee %s",
		ms.workDir, ms.logPath)

	cmd := exec.CommandContext(ctx, "screen", "-dmS", "muraena", "bash", "-c", screenCmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start muraena: %w", err)
	}

	// Wait for service to start
	time.Sleep(3 * time.Second)

	// Verify it started
	if !ms.isRunning() {
		return fmt.Errorf("muraena failed to start (screen session not found)")
	}

	return nil
}

// Stop stops the Muraena service
func (ms *MuraenaService) Stop(ctx context.Context) error {
	if !ms.isRunning() {
		return nil // Already stopped
	}

	cmd := exec.CommandContext(ctx, "screen", "-S", "muraena", "-X", "quit")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop muraena: %w", err)
	}

	// Wait for graceful shutdown
	time.Sleep(2 * time.Second)

	return nil
}

// Restart restarts the Muraena service
func (ms *MuraenaService) Restart(ctx context.Context) error {
	if err := ms.Stop(ctx); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	return ms.Start(ctx)
}

// HealthCheck checks if Muraena is healthy
func (ms *MuraenaService) HealthCheck(ctx context.Context) error {
	// Check if screen session exists
	if !ms.isRunning() {
		return fmt.Errorf("muraena screen session not found")
	}

	// Check if ports are listening
	ports := []int{80, 443}
	for _, port := range ports {
		if !ms.isPortListening(port) {
			return fmt.Errorf("port %d not listening", port)
		}
	}

	// Check if log file exists and is being written to
	if info, err := os.Stat(ms.logPath); err == nil {
		// Check if log was modified recently (within last 5 minutes)
		if time.Since(info.ModTime()) > 5*time.Minute {
			return fmt.Errorf("log file not being updated (last modified: %s)", info.ModTime())
		}
	}

	return nil
}

// GetStatus returns the current status
func (ms *MuraenaService) GetStatus(ctx context.Context) ServiceStatus {
	status := ServiceStatus{
		Name:  ms.GetName(),
		Ports: []int{80, 443},
	}

	if err := ms.HealthCheck(ctx); err == nil {
		status.Running = true
		status.Healthy = true
		status.LastSeen = time.Now()

		// Try to get PID
		if pid := ms.getPID(); pid > 0 {
			status.PID = pid
		}
	} else {
		status.Running = ms.isRunning()
		status.Healthy = false
		status.Errors = append(status.Errors, err.Error())
	}

	return status
}

// isRunning checks if the screen session exists
func (ms *MuraenaService) isRunning() bool {
	cmd := exec.Command("screen", "-list")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), "muraena")
}

// isPortListening checks if a port is listening
func (ms *MuraenaService) isPortListening(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// getPID attempts to get the process ID
func (ms *MuraenaService) getPID() int {
	cmd := exec.Command("pgrep", "-f", "muraena.bin")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	var pid int
	fmt.Sscanf(string(output), "%d", &pid)
	return pid
}

// GetLogs returns the last N lines of logs
func (ms *MuraenaService) GetLogs(lines int) (string, error) {
	cmd := exec.Command("tail", "-n", fmt.Sprintf("%d", lines), ms.logPath)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}
	return string(output), nil
}

// FollowLogs tails the log file in real-time
func (ms *MuraenaService) FollowLogs(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "tail", "-f", ms.logPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
