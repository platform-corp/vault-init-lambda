.DEFAULT_GOAL := all

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

# Application name
APPNAME = vault-init

# Directories
SRC_DIR = .
DIST_DIR = .

# Build targets
BINARY_NAME = $(DIST_DIR)/bootstrap

# Compressed output file
ZIP_FILE = $(DIST_DIR)/$(APPNAME).zip

# Default target OS and architecture
DEFAULT_GOOS = linux
DEFAULT_GOARCH = amd64

# User-defined or environment-specified values
GOOS ?= $(DEFAULT_GOOS)
GOARCH ?= $(DEFAULT_GOARCH)

# Build the application
build:
	env GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) $(SRC_DIR)/...

# Clean up build files
clean:
	$(GOCLEAN)
	rm -rf $(DIST_DIR)/$(APPNAME).zip
	rm -rf $(DIST_DIR)/bootstrap

# Compress the application binary into a zip file
compress: build
	mkdir -p $(DIST_DIR)
	zip $(ZIP_FILE) ${BINARY_NAME}

# Run tests
test:
	$(GOTEST) -v $(SRC_DIR)/...

# Install dependencies
deps: update
	$(GOGET) ./...

# Update dependencies
update:
	$(GOMOD) tidy

all: 
	@$(MAKE) deps || (echo "Error cleaning build files"; exit 1)
	@$(MAKE) build || (echo "Error building the application"; exit 1)
#	@$(MAKE) compress || (echo "Error compressing the application"; exit 1)

# Help target
help:
	@echo "Available targets:"
	@echo "  build      : Build the application"
	@echo "  clean      : Clean up build files"
	@echo "  compress   : Compress the application binary into a zip file"
	@echo "  test       : Run tests"
	@echo "  deps       : Install dependencies"
	@echo "  update     : Update dependencies"
	@echo "  help       : Show this help message"
