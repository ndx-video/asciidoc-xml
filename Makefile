# Makefile for AsciiDoc XML Converter
# Builds CLI and web binaries for multiple platforms

# Version (can be overridden via VERSION variable)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Build directory
BUILD_DIR = bin

# Target platforms
PLATFORMS = linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64

# Distribution directory
DIST_DIR = dist

# Go build flags
LDFLAGS = -s -w
ifneq ($(VERSION),dev)
	LDFLAGS += -X main.version=$(VERSION)
endif

# Default target
.DEFAULT_GOAL := build-cli

# Build CLI for current platform
.PHONY: cli
cli:
	@echo "Building adc for current platform..."
	@mkdir -p $(BUILD_DIR)/$(shell go env GOOS)-$(shell go env GOARCH)
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(shell go env GOOS)-$(shell go env GOARCH)/adc$(shell if [ "$(shell go env GOOS)" = "windows" ]; then echo ".exe"; fi) ./cli

# Build web server for current platform
.PHONY: web
web:
	@echo "Building asciidoc-xml-web for current platform..."
	@mkdir -p $(BUILD_DIR)/$(shell go env GOOS)-$(shell go env GOARCH)
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(shell go env GOOS)-$(shell go env GOARCH)/asciidoc-xml-web$(shell if [ "$(shell go env GOOS)" = "windows" ]; then echo ".exe"; fi) ./web

# Build CLI for all target platforms
.PHONY: build-cli
build-cli: clean-cli
	@echo "Building adc for all target platforms..."
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'-' -f1); \
		arch=$$(echo $$platform | cut -d'-' -f2); \
		echo "Building for $$os/$$arch..."; \
		mkdir -p $(BUILD_DIR)/$$platform; \
		GOOS=$$os GOARCH=$$arch go build -ldflags "$(LDFLAGS)" \
			-o $(BUILD_DIR)/$$platform/adc$$(if [ "$$os" = "windows" ]; then echo .exe; fi) \
			./cli; \
	done
	@echo "CLI build complete. Binaries in $(BUILD_DIR)/"

# Build web server for all target platforms
.PHONY: build-web
build-web: clean-web
	@echo "Building asciidoc-xml-web for all target platforms..."
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'-' -f1); \
		arch=$$(echo $$platform | cut -d'-' -f2); \
		echo "Building for $$os/$$arch..."; \
		mkdir -p $(BUILD_DIR)/$$platform; \
		GOOS=$$os GOARCH=$$arch go build -ldflags "$(LDFLAGS)" \
			-o $(BUILD_DIR)/$$platform/asciidoc-xml-web$$(if [ "$$os" = "windows" ]; then echo .exe; fi) \
			./web; \
	done
	@echo "Web server build complete. Binaries in $(BUILD_DIR)/"

# Build both CLI and web for all platforms
.PHONY: build-all
build-all: build-cli build-web

# Create CLI-only distribution packages
.PHONY: dist-cli
dist-cli: build-cli
	@echo "Creating CLI-only distribution packages..."
	@rm -rf $(DIST_DIR)
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		if [ ! -f $(BUILD_DIR)/$$platform/adc$$(if [ "$$(echo $$platform | cut -d'-' -f1)" = "windows" ]; then echo .exe; fi) ]; then \
			echo "Warning: Binary not found for $$platform, skipping..."; \
			continue; \
		fi; \
		package_name="asciidoc-xml-cli-$(VERSION)-$$platform"; \
		package_dir="$(DIST_DIR)/$$package_name"; \
		mkdir -p $$package_dir/bin $$package_dir/examples; \
		cp $(BUILD_DIR)/$$platform/adc$$(if [ "$$(echo $$platform | cut -d'-' -f1)" = "windows" ]; then echo .exe; fi) $$package_dir/bin/; \
		cp LICENSE $$package_dir/ 2>/dev/null || true; \
		cp README.md $$package_dir/ 2>/dev/null || true; \
		cp -r examples/* $$package_dir/examples/ 2>/dev/null || true; \
		if [ "$$(echo $$platform | cut -d'-' -f1)" = "windows" ]; then \
			cd $(DIST_DIR) && zip -r $$package_name.zip $$package_name > /dev/null 2>&1 || (cd $$package_name && zip -r ../$$package_name.zip . > /dev/null 2>&1); \
			echo "Created $$package_name.zip"; \
		else \
			cd $(DIST_DIR) && tar -czf $$package_name.tar.gz $$package_name; \
			echo "Created $$package_name.tar.gz"; \
		fi; \
		rm -rf $$package_dir; \
	done
	@echo "CLI distribution packages created in $(DIST_DIR)/"

# Create full distribution packages (CLI + web)
.PHONY: dist-full
dist-full: build-all
	@echo "Creating full distribution packages..."
	@rm -rf $(DIST_DIR)
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		if [ ! -f $(BUILD_DIR)/$$platform/adc$$(if [ "$$(echo $$platform | cut -d'-' -f1)" = "windows" ]; then echo .exe; fi) ]; then \
			echo "Warning: CLI binary not found for $$platform, skipping..."; \
			continue; \
		fi; \
		if [ ! -f $(BUILD_DIR)/$$platform/asciidoc-xml-web$$(if [ "$$(echo $$platform | cut -d'-' -f1)" = "windows" ]; then echo .exe; fi) ]; then \
			echo "Warning: Web binary not found for $$platform, skipping..."; \
			continue; \
		fi; \
		package_name="asciidoc-xml-full-$(VERSION)-$$platform"; \
		package_dir="$(DIST_DIR)/$$package_name"; \
		mkdir -p $$package_dir/bin $$package_dir/examples $$package_dir/xslt; \
		cp $(BUILD_DIR)/$$platform/adc$$(if [ "$$(echo $$platform | cut -d'-' -f1)" = "windows" ]; then echo .exe; fi) $$package_dir/bin/; \
		cp $(BUILD_DIR)/$$platform/asciidoc-xml-web$$(if [ "$$(echo $$platform | cut -d'-' -f1)" = "windows" ]; then echo .exe; fi) $$package_dir/bin/; \
		cp LICENSE $$package_dir/ 2>/dev/null || true; \
		cp README.md $$package_dir/ 2>/dev/null || true; \
		cp -r examples/* $$package_dir/examples/ 2>/dev/null || true; \
		cp -r xslt/* $$package_dir/xslt/ 2>/dev/null || true; \
		if [ "$$(echo $$platform | cut -d'-' -f1)" = "windows" ]; then \
			cd $(DIST_DIR) && zip -r $$package_name.zip $$package_name > /dev/null 2>&1 || (cd $$package_name && zip -r ../$$package_name.zip . > /dev/null 2>&1); \
			echo "Created $$package_name.zip"; \
		else \
			cd $(DIST_DIR) && tar -czf $$package_name.tar.gz $$package_name; \
			echo "Created $$package_name.tar.gz"; \
		fi; \
		rm -rf $$package_dir; \
	done
	@echo "Full distribution packages created in $(DIST_DIR)/"

# Clean build artifacts
.PHONY: clean-cli
clean-cli:
	@echo "Cleaning CLI build artifacts..."
	@find $(BUILD_DIR) -name "adc*" -type f -delete 2>/dev/null || true

.PHONY: clean-web
clean-web:
	@echo "Cleaning web build artifacts..."
	@find $(BUILD_DIR) -name "asciidoc-xml-web*" -type f -delete 2>/dev/null || true

.PHONY: clean
clean: clean-cli clean-web
	@echo "Cleaning all build artifacts..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)

# Test
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./... -v

# Install CLI to local system (current platform only)
.PHONY: install-cli
install-cli: cli
	@echo "Installing adc to local system..."
	@cp $(BUILD_DIR)/$(shell go env GOOS)-$(shell go env GOARCH)/adc$(shell if [ "$(shell go env GOOS)" = "windows" ]; then echo ".exe"; fi) $(shell go env GOPATH)/bin/ 2>/dev/null || \
		sudo cp $(BUILD_DIR)/$(shell go env GOOS)-$(shell go env GOARCH)/adc$(shell if [ "$(shell go env GOOS)" = "windows" ]; then echo ".exe"; fi) /usr/local/bin/
	@echo "Installation complete"

# Help
.PHONY: help
help:
	@echo "AsciiDoc XML Converter - Build System"
	@echo ""
	@echo "Targets:"
	@echo "  cli              Build adc CLI for current platform"
	@echo "  web              Build web server for current platform"
	@echo "  build-cli        Build adc CLI for all target platforms"
	@echo "  build-web        Build web server for all target platforms"
	@echo "  build-all        Build both CLI and web for all platforms"
	@echo "  dist-cli         Create CLI-only distribution packages"
	@echo "  dist-full        Create full distribution packages (CLI + web)"
	@echo "  clean            Remove all build artifacts"
	@echo "  clean-cli        Remove CLI build artifacts"
	@echo "  clean-web        Remove web build artifacts"
	@echo "  test             Run all tests"
	@echo "  install-cli      Install CLI to local system (current platform)"
	@echo "  help             Show this help message"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION          Version string for distribution packages (default: git describe or 'dev')"
	@echo ""
	@echo "Examples:"
	@echo "  make build-cli           # Build CLI for all platforms"
	@echo "  make dist-cli            # Create CLI distribution packages"
	@echo "  make VERSION=1.0.0 dist-full  # Create full packages with version 1.0.0"

