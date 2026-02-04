package service

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// RedisService manages the Redis service
type RedisService struct {
	config *ServiceConfig
}

// NewRedisService creates a new Redis service manager
func NewRedisService(config *ServiceConfig) *RedisService {
	return &RedisService{
		config: config,
	}
}

// GetName returns the service name
func (rs *RedisService) GetName() string {
	return "Redis"
}

// Start starts the Redis service
func (rs *RedisService) Start(ctx context.Context) error {
	// Check if already running
	if err := rs.HealthCheck(ctx); err == nil {
		return fmt.Errorf("redis is already running")
	}

	// Start via systemctl
	cmd := exec.CommandContext(ctx, "sudo", "systemctl", "start", "redis")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start redis: %w", err)
	}

	// Enable on boot
	cmd = exec.CommandContext(ctx, "sudo", "systemctl", "enable", "redis")
	_ = cmd.Run() // Don't fail if enable fails

	// Wait for service to be ready
	time.Sleep(2 * time.Second)

	return rs.HealthCheck(ctx)
}

// Stop stops the Redis service
func (rs *RedisService) Stop(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "sudo", "systemctl", "stop", "redis")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop redis: %w", err)
	}
	return nil
}

// Restart restarts the Redis service
func (rs *RedisService) Restart(ctx context.Context) error {
	if err := rs.Stop(ctx); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return rs.Start(ctx)
}

// HealthCheck checks if Redis is healthy
func (rs *RedisService) HealthCheck(ctx context.Context) error {
	// Check systemctl status
	cmd := exec.CommandContext(ctx, "systemctl", "is-active", "redis")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("redis is not active")
	}

	if strings.TrimSpace(string(output)) != "active" {
		return fmt.Errorf("redis status: %s", strings.TrimSpace(string(output)))
	}

	// Try to ping Redis
	cmd = exec.CommandContext(ctx, "redis-cli", "ping")
	output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	if strings.TrimSpace(string(output)) != "PONG" {
		return fmt.Errorf("unexpected redis ping response: %s", string(output))
	}

	return nil
}

// GetStatus returns the current status
func (rs *RedisService) GetStatus(ctx context.Context) ServiceStatus {
	status := ServiceStatus{
		Name:  rs.GetName(),
		Ports: []int{6379},
	}

	if err := rs.HealthCheck(ctx); err == nil {
		status.Running = true
		status.Healthy = true
		status.LastSeen = time.Now()
	} else {
		status.Running = false
		status.Healthy = false
		status.Errors = append(status.Errors, err.Error())
	}

	return status
}
