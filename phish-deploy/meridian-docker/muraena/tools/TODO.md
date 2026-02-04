# Muraena Shell-to-Go Refactoring TODO

## Project Overview
Refactor all Muraena-related shell scripts into a unified Go codebase for better maintainability, cross-platform support, and type safety.

---

## Phase 1: Foundation (Week 1)

### Setup & Infrastructure
- [ ] Initialize Go module in `muraena/tools/`
- [ ] Create project directory structure
- [ ] Set up Makefile for building
- [ ] Configure `.gitignore` for Go
- [ ] Set up CI/CD pipeline (GitHub Actions)

### Common Utilities
- [ ] Implement `pkg/common/logger.go` - Structured logging
- [ ] Implement `pkg/common/colors.go` - Terminal colors
- [ ] Implement `pkg/common/errors.go` - Error handling
- [ ] Implement `pkg/common/utils.go` - Utility functions
- [ ] Write unit tests for common package

### Remote Execution
- [ ] Implement `pkg/remote/ssh.go` - SSH client wrapper
- [ ] Implement `pkg/remote/executor.go` - Command execution
- [ ] Implement `pkg/remote/session.go` - Session management
- [ ] Implement `pkg/remote/transfer.go` - File transfer (SCP/SFTP)
- [ ] Write integration tests for remote package

---

## Phase 2: Service Management (Week 1-2)

### Core Service Manager
- [ ] Implement `pkg/service/manager.go` - Service orchestration
- [ ] Implement `pkg/service/redis.go` - Redis service control
- [ ] Implement `pkg/service/proxy.go` - Muraena proxy control
- [ ] Implement `pkg/service/necrobrowser.go` - NecroBrowser control
- [ ] Implement `pkg/service/health.go` - Health checks
- [ ] Write unit tests for service package

### CLI Tool: service-manager
- [ ] Implement `cmd/service-manager/main.go`
- [ ] Add `start` command - Start all services
- [ ] Add `stop` command - Stop all services
- [ ] Add `restart` command - Restart services
- [ ] Add `status` command - Show service status
- [ ] Add `logs` command - View service logs
- [ ] Add `health` command - Run health checks
- [ ] Write integration tests

### Shell Script Migration
- [ ] Migrate `scripts/service/start-services.sh` → Go
- [ ] Migrate `scripts/service/stop-services.sh` → Go
- [ ] Migrate `scripts/service-manager.sh` → Go
- [ ] Create backward-compatible shell wrappers
- [ ] Update documentation

---

## Phase 3: Configuration Management (Week 2)

### Configuration Core
- [ ] Implement `pkg/config/manager.go` - Config management
- [ ] Implement `pkg/config/generator.go` - TOML generation
- [ ] Implement `pkg/config/templates.go` - Template engine
- [ ] Implement `pkg/config/validator.go` - Config validation
- [ ] Implement `pkg/config/loader.go` - Config loading
- [ ] Write unit tests for config package

### Target Presets
- [ ] Implement `pkg/config/targets/westpac.go`
- [ ] Implement `pkg/config/targets/commbank.go`
- [ ] Implement `pkg/config/targets/anz.go`
- [ ] Implement `pkg/config/targets/nab.go`
- [ ] Implement `pkg/config/targets/registry.go` - Target registry
- [ ] Create YAML preset definitions

### CLI Tool: config-manager
- [ ] Implement `cmd/config-manager/main.go`
- [ ] Add `list` command - List available targets
- [ ] Add `status` command - Show current config
- [ ] Add `set` command - Set target configuration
- [ ] Add `generate` command - Generate config from template
- [ ] Add `validate` command - Validate configuration
- [ ] Add `backup` command - Backup current config
- [ ] Add `restore` command - Restore config from backup
- [ ] Write integration tests

### Shell Script Migration
- [ ] Migrate `scripts/config/change-target.sh` → Go
- [ ] Migrate `scripts/config/change-target-server.sh` → Go
- [ ] Migrate `scripts/config/setup-config.sh` → Go
- [ ] Migrate `scripts/config/update-target-westpac.sh` → Go
- [ ] Migrate `scripts/config/fix-meridian.sh` → Go
- [ ] Create backward-compatible shell wrappers
- [ ] Update documentation

---

## Phase 4: Credential Extraction (Week 2-3)

### Extraction Core
- [ ] Implement `pkg/extract/redis.go` - Redis data extraction
- [ ] Implement `pkg/extract/parser.go` - Data parsing
- [ ] Implement `pkg/extract/exporter.go` - Export formats
- [ ] Implement `pkg/extract/filter.go` - Data filtering
- [ ] Implement `pkg/extract/analyzer.go` - Data analysis
- [ ] Write unit tests for extract package

### Export Formats
- [ ] Implement CSV export
- [ ] Implement JSON export
- [ ] Implement XML export
- [ ] Implement HTML report generation
- [ ] Implement encrypted export (GPG)

### CLI Tool: credential-extractor
- [ ] Implement `cmd/credential-extractor/main.go`
- [ ] Add `list` command - List captured credentials
- [ ] Add `export` command - Export credentials
- [ ] Add `stats` command - Show statistics
- [ ] Add `search` command - Search credentials
- [ ] Add `victims` command - List victims
- [ ] Add `sessions` command - List sessions
- [ ] Write integration tests

### Shell Script Migration
- [ ] Migrate `scripts/extract/extract-credentials.sh` → Go
- [ ] Migrate `scripts/extract/check-credentials.sh` → Go
- [ ] Create backward-compatible shell wrappers
- [ ] Update documentation

---

## Phase 5: SSL Management (Week 3)

### SSL Core
- [ ] Implement `pkg/ssl/certbot.go` - Certbot integration
- [ ] Implement `pkg/ssl/validator.go` - Certificate validation
- [ ] Implement `pkg/ssl/renewal.go` - Auto-renewal
- [ ] Implement `pkg/ssl/generator.go` - Self-signed certs
- [ ] Write unit tests for ssl package

### CLI Tool: ssl-manager
- [ ] Implement `cmd/ssl-manager/main.go`
- [ ] Add `generate` command - Generate certificate
- [ ] Add `renew` command - Renew certificate
- [ ] Add `validate` command - Validate certificate
- [ ] Add `info` command - Show certificate info
- [ ] Add `auto-renew` command - Setup auto-renewal
- [ ] Write integration tests

### Shell Script Migration
- [ ] Migrate SSL-related functions from `deploy-meridian.sh`
- [ ] Migrate `scripts/test/test-ssl-autoconfig.sh` → Go
- [ ] Create backward-compatible shell wrappers
- [ ] Update documentation

---

## Phase 6: Deployment (Week 3-4)

### Deployment Core
- [ ] Implement `pkg/deploy/ec2.go` - EC2 deployment
- [ ] Implement `pkg/deploy/ssh.go` - SSH operations
- [ ] Implement `pkg/deploy/transfer.go` - File transfer
- [ ] Implement `pkg/deploy/campaign.go` - Campaign deployment
- [ ] Implement `pkg/deploy/verification.go` - Deployment verification
- [ ] Write unit tests for deploy package

### CLI Tool: deployer
- [ ] Implement `cmd/deployer/main.go`
- [ ] Add `deploy` command - Full deployment
- [ ] Add `transfer` command - Transfer files
- [ ] Add `verify` command - Verify deployment
- [ ] Add `rollback` command - Rollback deployment
- [ ] Add `campaign` command - Deploy campaign
- [ ] Write integration tests

### Shell Script Migration
- [ ] Migrate `scripts/deploy/deploy-to-ec2.sh` → Go
- [ ] Migrate `scripts/deploy/deploy-campaign.sh` → Go
- [ ] Migrate `scripts/deploy/finalize-deployment.sh` → Go
- [ ] Migrate `deploy-meridian.sh` (main deployment) → Go
- [ ] Create backward-compatible shell wrappers
- [ ] Update documentation

---

## Phase 7: Testing Framework (Week 4)

### Testing Core
- [ ] Implement `pkg/test/workflow.go` - Workflow testing
- [ ] Implement `pkg/test/integration.go` - Integration tests
- [ ] Implement `pkg/test/validator.go` - Validation tests
- [ ] Implement `pkg/test/reporter.go` - Test reporting
- [ ] Write unit tests for test package

### CLI Tool: test-runner
- [ ] Implement `cmd/test-runner/main.go`
- [ ] Add `run` command - Run test suite
- [ ] Add `workflow` command - Run workflow tests
- [ ] Add `integration` command - Run integration tests
- [ ] Add `report` command - Generate test report
- [ ] Write integration tests

### Shell Script Migration
- [ ] Migrate `scripts/test/test-complete-workflow.sh` → Go
- [ ] Migrate `scripts/test/final_verification.sh` → Go
- [ ] Migrate `scripts/test/simulate_capture.sh` → Go
- [ ] Migrate `scripts/test-suite.sh` → Go
- [ ] Create backward-compatible shell wrappers
- [ ] Update documentation

---

## Phase 8: Unified CLI (Week 4-5)

### Main CLI Tool
- [ ] Implement `cmd/muraena-cli/main.go` - Unified CLI
- [ ] Integrate all sub-commands
- [ ] Add `service` subcommand group
- [ ] Add `config` subcommand group
- [ ] Add `extract` subcommand group
- [ ] Add `ssl` subcommand group
- [ ] Add `deploy` subcommand group
- [ ] Add `test` subcommand group
- [ ] Add `version` command
- [ ] Add `help` command with examples
- [ ] Implement shell completion (bash, zsh, fish)
- [ ] Write comprehensive integration tests

### Configuration
- [ ] Implement global config file support
- [ ] Implement environment variable support
- [ ] Implement config file discovery
- [ ] Implement config validation
- [ ] Add config migration tool

---

## Phase 9: Documentation & Migration (Week 5)

### Documentation
- [ ] Write comprehensive README.md
- [ ] Write installation guide
- [ ] Write migration guide (shell → Go)
- [ ] Write API documentation
- [ ] Write examples and tutorials
- [ ] Create man pages
- [ ] Create video tutorials (optional)

### Migration Tools
- [ ] Create migration script (shell → Go)
- [ ] Create compatibility layer
- [ ] Create deprecation warnings
- [ ] Update all existing documentation
- [ ] Update deployment guides

### Backward Compatibility
- [ ] Create shell wrapper scripts
- [ ] Ensure API compatibility
- [ ] Add deprecation notices
- [ ] Create migration timeline

---

## Phase 10: Testing & Quality Assurance (Week 5-6)

### Testing
- [ ] Write unit tests (target: 80%+ coverage)
- [ ] Write integration tests
- [ ] Write end-to-end tests
- [ ] Perform load testing
- [ ] Perform security testing
- [ ] Perform cross-platform testing (Linux, macOS, Windows)

### Code Quality
- [ ] Run `go vet`
- [ ] Run `golint`
- [ ] Run `staticcheck`
- [ ] Run `gosec` (security scanner)
- [ ] Fix all linter warnings
- [ ] Optimize performance bottlenecks

### CI/CD
- [ ] Set up GitHub Actions
- [ ] Add automated testing
- [ ] Add automated builds
- [ ] Add automated releases
- [ ] Add code coverage reporting
- [ ] Add security scanning

---

## Phase 11: Release & Deployment (Week 6)

### Release Preparation
- [ ] Finalize version 1.0.0
- [ ] Create release notes
- [ ] Create changelog
- [ ] Tag release in Git
- [ ] Build binaries for all platforms
- [ ] Create installation packages (deb, rpm, brew)

### Deployment
- [ ] Deploy to production
- [ ] Monitor for issues
- [ ] Gather user feedback
- [ ] Fix critical bugs
- [ ] Plan next iteration

---

## Ongoing Maintenance

### Regular Tasks
- [ ] Monitor for security vulnerabilities
- [ ] Update dependencies
- [ ] Fix reported bugs
- [ ] Add requested features
- [ ] Improve documentation
- [ ] Optimize performance

---

## Success Metrics

### Code Quality
- [ ] 80%+ test coverage
- [ ] Zero critical security issues
- [ ] All linters passing
- [ ] Performance benchmarks met

### Functionality
- [ ] All shell scripts migrated
- [ ] Feature parity achieved
- [ ] Backward compatibility maintained
- [ ] Cross-platform support verified

### Documentation
- [ ] Complete API documentation
- [ ] Migration guide published
- [ ] Examples and tutorials available
- [ ] User feedback incorporated

---

## Notes

### Dependencies
- Go 1.21+
- Redis client library
- SSH library (golang.org/x/crypto/ssh)
- TOML parser
- Cobra (CLI framework)
- Viper (configuration)

### Challenges
- Maintaining backward compatibility
- Cross-platform SSH operations
- Screen session management in Go
- Certbot integration
- Remote command execution

### Opportunities
- Better error handling
- Type safety
- Cross-platform support
- Easier testing
- Better maintainability
- Single binary distribution

---

**Last Updated:** January 2026  
**Status:** Planning Phase  
**Target Completion:** 6 weeks from start
