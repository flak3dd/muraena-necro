package service

import (
	"context"
	"fmt"
	"time"
)

// Manager orchestrates multiple services
type Manager struct {
	config   *ServiceConfig
	services map[string]Service
}

// NewManager creates a new service manager
func NewManager(config *ServiceConfig) *Manager {
	if config == nil {
		config = DefaultServiceConfig()
	}

	m := &Manager{
		config:   config,
		services: make(map[string]Service),
	}

	// Register services
	m.services["redis"] = NewRedisService(config)
	m.services["muraena"] = NewMuraenaService(config)
	m.services["necrobrowser"] = NewNecroBrowserService(config)

	return m
}

// StartAll starts all services in order
func (m *Manager) StartAll(ctx context.Context) error {
	// Start services in dependency order
	serviceOrder := []string{"redis", "muraena", "necrobrowser"}

	for _, name := range serviceOrder {
		service, ok := m.services[name]
		if !ok {
			continue // Skip if service not registered yet
		}

		if err := service.Start(ctx); err != nil {
			return fmt.Errorf("failed to start %s: %w", name, err)
		}

		// Wait between service starts
		time.Sleep(2 * time.Second)
	}

	return m.VerifyAll(ctx)
}

// StopAll stops all services in reverse order
func (m *Manager) StopAll(ctx context.Context) error {
	// Stop in reverse order
	serviceOrder := []string{"necrobrowser", "muraena", "redis"}

	var errors []error
	for _, name := range serviceOrder {
		service, ok := m.services[name]
		if !ok {
			continue
		}

		if err := service.Stop(ctx); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors stopping services: %v", errors)
	}

	return nil
}

// RestartAll restarts all services
func (m *Manager) RestartAll(ctx context.Context) error {
	if err := m.StopAll(ctx); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	return m.StartAll(ctx)
}

// VerifyAll verifies all services are running
func (m *Manager) VerifyAll(ctx context.Context) error {
	for name, service := range m.services {
		if err := service.HealthCheck(ctx); err != nil {
			return fmt.Errorf("%s health check failed: %w", name, err)
		}
	}
	return nil
}

// GetStatus returns status of all services
func (m *Manager) GetStatus(ctx context.Context) map[string]ServiceStatus {
	status := make(map[string]ServiceStatus)

	for name, service := range m.services {
		status[name] = service.GetStatus(ctx)
	}

	return status
}

// GetService returns a specific service by name
func (m *Manager) GetService(name string) (Service, error) {
	service, ok := m.services[name]
	if !ok {
		return nil, fmt.Errorf("service not found: %s", name)
	}
	return service, nil
}
