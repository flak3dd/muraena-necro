package service

import (
	"context"
	"time"
)

// Service represents a manageable service
type Service interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Restart(ctx context.Context) error
	HealthCheck(ctx context.Context) error
	GetStatus(ctx context.Context) ServiceStatus
	GetName() string
}

// ServiceStatus represents the current status of a service
type ServiceStatus struct {
	Name     string        `json:"name"`
	Running  bool          `json:"running"`
	PID      int           `json:"pid,omitempty"`
	Uptime   time.Duration `json:"uptime,omitempty"`
	Ports    []int         `json:"ports,omitempty"`
	Errors   []string      `json:"errors,omitempty"`
	Healthy  bool          `json:"healthy"`
	LastSeen time.Time     `json:"last_seen,omitempty"`
}

// ServiceConfig holds configuration for services
type ServiceConfig struct {
	RedisAddr          string
	RedisPassword      string
	MuraenaDir         string
	MuraenaBinary      string
	MuraenaConfig      string
	NecroBrowserDir    string
	NecroBrowserBinary string
	NecroBrowserConfig string
	LogDir             string
}

// DefaultServiceConfig returns default configuration
func DefaultServiceConfig() *ServiceConfig {
	return &ServiceConfig{
		RedisAddr:          "localhost:6379",
		RedisPassword:      "",
		MuraenaDir:         "/home/ubuntu/muraena",
		MuraenaBinary:      "muraena.bin",
		MuraenaConfig:      "config.toml",
		NecroBrowserDir:    "/home/ubuntu/necrobrowser",
		NecroBrowserBinary: "necrobrowser.js",
		NecroBrowserConfig: "config.toml",
		LogDir:             "/home/ubuntu/logs",
	}
}
