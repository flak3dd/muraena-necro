# Muraena Tools - Go Refactoring Project

## Overview

This directory contains the Go-based refactoring of all Muraena-related shell scripts. The goal is to create a unified, cross-platform, maintainable codebase that replaces the legacy bash scripts with type-safe Go implementations.

## Project Status

**Status:** Initial Setup Phase  
**Version:** 0.1.0-alpha  
**Target Completion:** 6 weeks from start date

## Architecture

```
muraena/tools/
├── cmd/                          # Command-line applications
│   ├── service-manager/          # Service management CLI
│   ├── config-manager/           # Configuration management CLI
│   ├── credential-extractor/     # Credential extraction CLI
│   ├── ssl-manager/              # SSL certificate management CLI
│   ├── deployer/                 # Deployment automation CLI
│   └── muraena-cli/              # Unified CLI (combines all above)
├── pkg/                          # Reusable packages
│   ├── common/                   # Common utilities
│   ├── service/                  # Service management
│   ├── config/                   # Configuration management
│   ├── extract/                  # Credential extraction
│   ├── ssl/                      # SSL management
│   ├── deploy/                   # Deployment automation
│   ├── test/                     # Testing framework
│   └── remote/                   # Remote execution (SSH)
├── internal/                     # Internal packages
│   ├── models/                   # Data models
│   └── constants/                # Constants
├── configs/                      # Configuration templates
│   ├── templates/                # TOML templates
│   └── presets/                  # Target presets
├── go.mod                        # Go module definition
├── go.sum                        # Dependency checksums
├── Makefile                      # Build automation
├── README.md                     # This file
└── TODO.md                       # Detailed task list
```

## Features

### Completed
- ✅ Project structure created
- ✅ Go module initialized
- ✅ Basic logger implementation
- ✅ Service manager CLI skeleton
- ✅ Comprehensive refactoring plan
- ✅ Detailed TODO list

### In Progress
- 🔄 Service management package
- 🔄 Configuration management package
- 🔄 Remote execution package

### Planned
- ⏳ Credential extraction package
- ⏳ SSL management package
- ⏳ Deployment automation package
- ⏳ Testing framework
- ⏳ Unified CLI tool

## Shell Scripts Being Replaced

### Service Management
- `scripts/service/start-services.sh` → `service-manager start`
- `scripts/service/stop-services.sh` → `service-manager stop`
- `scripts/service-manager.sh` → `service-manager`

### Configuration Management
- `scripts/config/change-target.sh` → `config-manager set`
- `scripts/config/setup-config.sh` → `config-manager init`
- `scripts/config/update-target-westpac.sh` → `config-manager set westpac`

### Credential Extraction
- `scripts/extract/extract-credentials.sh` → `credential-extractor export`
- `scripts/extract/check-credentials.sh` → `credential-extractor list`

### Deployment
- `scripts/deploy/deploy-to-ec2.sh` → `deployer ec2`
- `scripts/deploy/deploy-campaign.sh` → `deployer campaign`
- `deploy-meridian.sh` → `deployer full`

### Testing
- `scripts/test/test-complete-workflow.sh` → `test-runner workflow`
- `scripts/test/final_verification.sh` → `test-runner verify`

## Installation

### Prerequisites
- Go 1.21 or higher
- Access to the Muraena infrastructure
- SSH keys configured

### Build from Source

```bash
# Navigate to tools directory
cd muraena/tools

# Download dependencies
go mod download

# Build all tools
make build

# Or build specific tool
make build-service-manager
make build-config-manager
make build-credential-extractor

# Install to system
make install
```

### Using Pre-built Binaries

```bash
# Download latest release
wget https://github.com/muraenateam/muraena/releases/latest/download/muraena-tools-linux-amd64.tar.gz

# Extract
tar -xzf muraena-tools-linux-amd64.tar.gz

# Move to PATH
sudo mv muraena-tools/* /usr/local/bin/
```

## Usage

### Service Manager

```bash
# Start all services
service-manager start

# Stop all services
service-manager stop

# Restart services
service-manager restart

# Check status
service-manager status

# View logs
service-manager logs muraena

# Run health checks
service-manager health
```

### Configuration Manager

```bash
# List available targets
config-manager list

# Show current configuration
config-manager status

# Set target (e.g., Westpac)
config-manager set westpac sect00.com

# Generate configuration
config-manager generate --target westpac --domain sect00.com

# Validate configuration
config-manager validate

# Backup configuration
config-manager backup

# Restore configuration
config-manager restore backup-20260115.tar.gz
```

### Credential Extractor

```bash
# List all captured credentials
credential-extractor list

# Export to CSV
credential-extractor export --format csv --output creds.csv

# Export to JSON
credential-extractor export --format json --output creds.json

# Show statistics
credential-extractor stats

# Search credentials
credential-extractor search --username "victim@email.com"

# List victims
credential-extractor victims

# List sessions
credential-extractor sessions
```

### SSL Manager

```bash
# Generate SSL certificate
ssl-manager generate --domain sect00.com

# Renew certificate
ssl-manager renew --domain sect00.com

# Validate certificate
ssl-manager validate --domain sect00.com

# Show certificate info
ssl-manager info --domain sect00.com

# Setup auto-renewal
ssl-manager auto-renew --enable
```

### Deployer

```bash
# Full deployment
deployer full --domain sect00.com --target westpac

# Deploy to EC2
deployer ec2 --instance ec2-54-81-35-64.compute-1.amazonaws.com

# Deploy campaign
deployer campaign --name q1_2026 --target westpac

# Verify deployment
deployer verify

# Rollback deployment
deployer rollback
```

### Unified CLI (Future)

```bash
# All commands will be available under single CLI
muraena service start
muraena config set westpac sect00.com
muraena extract list
muraena ssl generate sect00.com
muraena deploy full --target westpac
```

## Development

### Project Structure

```go
// Example: Service Manager
package service

type ServiceManager struct {
    logger  *common.Logger
    redis   *RedisService
    muraena *MuraenaService
    necro   *NecroBrowserService
}

func (sm *ServiceManager) StartAll(ctx context.Context) error {
    // Start Redis
    if err := sm.redis.Start(ctx); err != nil {
        return err
    }
    
    // Start Muraena
    if err := sm.muraena.Start(ctx); err != nil {
        return err
    }
    
    // Start NecroBrowser
    if err := sm.necro.Start(ctx); err != nil {
        return err
    }
    
    return sm.VerifyAll(ctx)
}
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./pkg/service/...

# Run integration tests
make test-integration
```

### Code Quality

```bash
# Run linters
make lint

# Format code
make fmt

# Run security scanner
make security-scan

# Generate documentation
make docs
```

## Migration Guide

### For Existing Users

1. **Install Go tools** alongside existing shell scripts
2. **Test Go equivalents** in non-production environment
3. **Gradually migrate** workflows to Go tools
4. **Deprecate shell scripts** after verification
5. **Remove shell scripts** after full migration

### Backward Compatibility

Shell wrapper scripts will be provided for backward compatibility:

```bash
# Old way (still works)
./scripts/service/start-services.sh

# New way (recommended)
service-manager start

# Wrapper script (transitional)
./scripts/service/start-services.sh  # Calls service-manager internally
```

## Benefits of Go Refactoring

### Technical Benefits
- ✅ **Type Safety** - Compile-time error detection
- ✅ **Cross-Platform** - Single binary for Linux, macOS, Windows
- ✅ **Better Error Handling** - Structured error propagation
- ✅ **Easier Testing** - Unit and integration tests
- ✅ **Better Maintainability** - Clear code structure
- ✅ **Performance** - Compiled binary vs interpreted shell

### Operational Benefits
- ✅ **Single Binary** - No dependency hell
- ✅ **Consistent CLI** - Unified command structure
- ✅ **Better Logging** - Structured, colored output
- ✅ **Progress Indicators** - Real-time feedback
- ✅ **Configuration Management** - Type-safe configs
- ✅ **Remote Execution** - Built-in SSH support

## Dependencies

### Core Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `github.com/fatih/color` - Terminal colors
- `github.com/go-redis/redis/v8` - Redis client
- `github.com/pelletier/go-toml/v2` - TOML parser
- `golang.org/x/crypto` - SSH and crypto

### Development Dependencies
- `github.com/stretchr/testify` - Testing framework
- `github.com/golangci/golangci-lint` - Linter
- `github.com/securego/gosec` - Security scanner

## Contributing

### Development Workflow

1. **Fork the repository**
2. **Create feature branch** (`git checkout -b feature/amazing-feature`)
3. **Write tests** for new functionality
4. **Implement feature** with proper error handling
5. **Run tests** (`make test`)
6. **Run linters** (`make lint`)
7. **Commit changes** (`git commit -m 'Add amazing feature'`)
8. **Push to branch** (`git push origin feature/amazing-feature`)
9. **Open Pull Request**

### Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Write meaningful commit messages
- Add comments for exported functions
- Keep functions small and focused

## Roadmap

### Phase 1: Foundation (Week 1)
- [x] Project setup
- [x] Common utilities
- [ ] Remote execution package
- [ ] Service management core

### Phase 2: Service Management (Week 1-2)
- [ ] Complete service manager
- [ ] Health checks
- [ ] Log management
- [ ] Shell script migration

### Phase 3: Configuration (Week 2)
- [ ] Config generator
- [ ] Target presets
- [ ] Validation
- [ ] Shell script migration

### Phase 4: Extraction (Week 2-3)
- [ ] Redis extraction
- [ ] Export formats
- [ ] Data analysis
- [ ] Shell script migration

### Phase 5: SSL & Deployment (Week 3-4)
- [ ] SSL management
- [ ] Deployment automation
- [ ] Campaign management
- [ ] Shell script migration

### Phase 6: Testing & QA (Week 5-6)
- [ ] Comprehensive testing
- [ ] Documentation
- [ ] Performance optimization
- [ ] Release preparation

## License

BSD 3-Clause License (same as Muraena)

## Support

- **Documentation:** See `REFACTORING_PLAN.md` for detailed plan
- **Issues:** GitHub Issues
- **Discussions:** GitHub Discussions

## Acknowledgments

- Original Muraena Team for the shell scripts
- Go community for excellent libraries
- Contributors to this refactoring effort

---

**Last Updated:** January 2026  
**Status:** Active Development  
**Maintainer:** Muraena Team
