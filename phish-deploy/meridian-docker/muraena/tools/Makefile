.PHONY: all build clean test help build-all

all: build-all

build: build-all

build-all: build-service-manager build-config-manager build-credential-extractor build-ssl-manager build-deployer

build-service-manager:
	@echo "Building service-manager..."
	go build -o bin/service-manager cmd/service-manager/main.go

build-config-manager:
	@echo "Building config-manager..."
	go build -o bin/config-manager cmd/config-manager/main.go

build-credential-extractor:
	@echo "Building credential-extractor..."
	go build -o bin/credential-extractor cmd/credential-extractor/main.go

build-ssl-manager:
	@echo "Building ssl-manager..."
	go build -o bin/ssl-manager cmd/ssl-manager/main.go

build-deployer:
	@echo "Building deployer..."
	go build -o bin/deployer cmd/deployer/main.go

build-linux: build-linux-service-manager build-linux-config-manager build-linux-credential-extractor build-linux-ssl-manager build-linux-deployer

build-linux-service-manager:
	@echo "Building service-manager for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/service-manager cmd/service-manager/main.go

build-linux-config-manager:
	@echo "Building config-manager for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/config-manager cmd/config-manager/main.go

build-linux-credential-extractor:
	@echo "Building credential-extractor for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/credential-extractor cmd/credential-extractor/main.go

build-linux-ssl-manager:
	@echo "Building ssl-manager for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/ssl-manager cmd/ssl-manager/main.go

build-linux-deployer:
	@echo "Building deployer for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/deployer cmd/deployer/main.go

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

test:
	@echo "Running tests..."
	go test -v ./...

help:
	@echo "Available targets:"
	@echo "  all                          - Build all binaries (default)"
	@echo "  build                        - Build all binaries"
	@echo "  build-all                    - Build all binaries"
	@echo "  build-service-manager        - Build service-manager binary"
	@echo "  build-config-manager         - Build config-manager binary"
	@echo "  build-credential-extractor   - Build credential-extractor binary"
	@echo "  build-ssl-manager            - Build ssl-manager binary"
	@echo "  build-deployer               - Build deployer binary"
	@echo "  build-linux                  - Build all binaries for Linux"
	@echo "  clean                        - Remove build artifacts"
	@echo "  test                         - Run tests"
	@echo "  help                         - Show this help message"
