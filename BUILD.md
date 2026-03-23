# Building and Installing azcopy

This guide covers how to build and install azcopy from source code, including all dependencies and platform-specific considerations.

## Prerequisites

### Go Environment
- **Go 1.21 or later** is required
- Ensure `GOPATH` and `GOROOT` are properly configured
- Verify installation: `go version`

### System Requirements

#### Linux
- GCC compiler toolchain
- libc development headers
- Git for source code management

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install build-essential git

# RHEL/CentOS/Fedora
sudo yum groupinstall "Development Tools"
sudo yum install git
```

#### macOS
- Xcode Command Line Tools
- Homebrew (recommended for dependency management)

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Install Homebrew (if not already installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

#### Windows
- Go for Windows
- Git for Windows
- Visual Studio Build Tools or MinGW-w64

## Building from Source

### 1. Clone the Repository

```bash
git clone https://github.com/Azure/azure-storage-azcopy.git
cd azure-storage-azcopy
```

### 2. Install Dependencies

The project uses Go modules for dependency management. Dependencies will be automatically downloaded during the build process.

```bash
# Download and verify dependencies
go mod download
go mod verify
```

### 3. Build azcopy

#### Standard Build

```bash
# Build for current platform
go build -o azcopy

# Or use the Makefile (if available)
make build
```

#### Cross-Platform Build

```bash
# Build for Linux (from any platform)
GOOS=linux GOARCH=amd64 go build -o azcopy-linux-amd64

# Build for macOS (from any platform)  
GOOS=darwin GOARCH=amd64 go build -o azcopy-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o azcopy-darwin-arm64

# Build for Windows (from any platform)
GOOS=windows GOARCH=amd64 go build -o azcopy-windows-amd64.exe
```

#### Build with Version Information

```bash
# Build with version and build info
go build -ldflags "-X 'github.com/Azure/azure-storage-azcopy/v10/common.azcopyVersion=10.x.x'" -o azcopy
```

#### Optimized Release Build

```bash
# Build optimized binary (smaller size, no debug symbols)
go build -ldflags "-s -w" -o azcopy

# Build with static linking (recommended for distribution)
CGO_ENABLED=0 go build -ldflags "-s -w -extldflags '-static'" -o azcopy
```

## Dependencies

### Runtime Dependencies

#### Process Locking (Unix/Linux/macOS)
- **github.com/gofrs/flock**: File-based process locking
- Automatically included in Go module dependencies
- Not required on Windows (no-op implementation)

#### Core Dependencies
- **Azure SDK for Go**: Azure service integration
- **Various Azure storage libraries**: Blob, File, Queue services
- **HTTP/networking libraries**: For data transfer operations

### Development Dependencies
- **Testing frameworks**: For running unit and integration tests
- **Mock libraries**: For testing Azure service interactions
- **Benchmarking tools**: For performance testing

## Installation

### Option 1: Install from Built Binary

```bash
# Copy to system binary directory
sudo cp azcopy /usr/local/bin/

# Make executable (if needed)
chmod +x /usr/local/bin/azcopy

# Verify installation
azcopy --version
```

### Option 2: Install from Source (Go Install)

```bash
# Install directly from source (latest main branch)
go install github.com/Azure/azure-storage-azcopy/v10@latest

# Install specific version/tag
go install github.com/Azure/azure-storage-azcopy/v10@v10.x.x
```

### Option 3: Package Managers

#### macOS (Homebrew)
```bash
brew install azcopy
```

#### Linux (Package Managers)
```bash
# Various distributions have different package names
# Check your distribution's documentation
```

#### Windows (Chocolatey/Scoop)
```bash
# Chocolatey
choco install azcopy

# Scoop
scoop install azcopy
```

## Platform-Specific Considerations

### Unix/Linux/macOS Features
- **Process-level locking**: Full flock-based inter-process synchronization
- **Signal handling**: Graceful shutdown on SIGTERM/SIGINT
- **File permissions**: Proper handling of Unix file attributes
- **Symbolic links**: Native symlink support

### Windows Features
- **Thread-level synchronization**: In-process locking only
- **Windows file attributes**: ACL and security descriptor handling
- **UNC paths**: Network path support
- **Windows-specific error handling**

## Development Setup

### Setting up Development Environment

```bash
# Clone and setup
git clone https://github.com/Azure/azure-storage-azcopy.git
cd azure-storage-azcopy

# Install dependencies
go mod download

# Install development tools (optional)
go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/lint/golint@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Running Tests

```bash
# Run unit tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test package
go test ./common
go test ./ste
go test ./cmd

# Run integration tests (requires Azure credentials)
go test -tags=integration ./...
```

### Code Quality Checks

```bash
# Format code
go fmt ./...

# Import organization
goimports -w .

# Lint code
golangci-lint run

# Vet code for common issues
go vet ./...
```

## Build Configuration

### Environment Variables

```bash
# Enable CGO (required for some features)
export CGO_ENABLED=1

# Set build mode
export GO111MODULE=on

# Cross-compilation settings
export GOOS=linux
export GOARCH=amd64

# Optimization settings
export GODEBUG=netdns=go  # Use Go's DNS resolver
```

### Build Tags

```bash
# Build without certain features
go build -tags="no_process_lock" -o azcopy

# Build for specific environments
go build -tags="development" -o azcopy-dev
go build -tags="production" -o azcopy-prod
```

## Performance Optimization

### Build Optimizations

```bash
# Minimal binary size
go build -ldflags "-s -w" -o azcopy

# Static linking (no external dependencies)
CGO_ENABLED=0 go build -ldflags "-s -w -extldflags '-static'" -o azcopy

# Profile-guided optimization (PGO) - Go 1.21+
go build -pgo=auto -o azcopy
```

### Runtime Optimizations

```bash
# Set Go runtime parameters
export GOMAXPROCS=4
export GOGC=100
export GODEBUG=gctrace=1
```

## Verification

### Test Installation

```bash
# Basic functionality test
azcopy --version
azcopy --help

# Test authentication (requires Azure credentials)
azcopy login

# Test basic operations (optional)
azcopy list "https://myaccount.blob.core.windows.net/mycontainer"
```

### Validate Features

```bash
# Test process locking (Unix/Linux/macOS only)
# Run multiple azcopy instances simultaneously
azcopy copy source1 dest1 &
azcopy copy source2 dest2 &
wait

# Check for any lock-related errors in logs
```

## Troubleshooting

### Common Build Issues

#### Go Version Compatibility
```bash
# Error: "go: module requires Go 1.21"
# Solution: Update Go version
go version
# If outdated, install latest Go from https://golang.org/dl/
```

#### Missing Dependencies
```bash
# Error: "package not found"
# Solution: Clean and rebuild dependency cache
go clean -modcache
go mod download
go mod tidy
```

#### CGO Issues
```bash
# Error: CGO-related build failures
# Solution: Ensure proper C compiler setup
export CGO_ENABLED=1
# Install build tools as described in prerequisites
```

### Platform-Specific Issues

#### Linux Build Issues
```bash
# Static linking issues
# Solution: Use musl instead of glibc
docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:alpine go build -v
```

#### macOS Build Issues
```bash
# Xcode command line tools missing
# Solution: Install/update Xcode tools
xcode-select --install
sudo xcode-select --switch /Applications/Xcode.app/Contents/Developer
```

#### Windows Build Issues
```bash
# Build tools not found
# Solution: Install Visual Studio Build Tools or use Windows Subsystem for Linux (WSL)
```

### Runtime Issues

#### Permission Errors (Unix/Linux)
```bash
# Error: "permission denied"
# Solution: Check file permissions and ownership
ls -la azcopy
chmod +x azcopy
```

#### Lock File Issues
```bash
# Error: "failed to acquire lock"
# Solution: Check lock directory permissions
ls -la ~/.azcopy/.locks/  # or job plan directory
rm -f ~/.azcopy/.locks/*.lock  # Remove stale locks if needed
```

## Contributing

### Development Workflow

1. **Fork and Clone**
2. **Create Feature Branch**: `git checkout -b feature/my-feature`
3. **Make Changes**
4. **Test Thoroughly**: `go test ./...`
5. **Check Code Quality**: `golangci-lint run`
6. **Commit**: `git commit -m "Add my feature"`
7. **Push**: `git push origin feature/my-feature`
8. **Create Pull Request**

### Code Standards

- Follow Go conventions and best practices
- Add tests for new functionality
- Update documentation for API changes
- Ensure backward compatibility when possible
- Use meaningful commit messages

## Release Process

### Creating Releases

```bash
# Tag a release
git tag -a v10.x.x -m "Release v10.x.x"
git push origin v10.x.x

# Build release binaries
make release  # or run build script

# Create release artifacts
tar czf azcopy-v10.x.x-linux-amd64.tar.gz azcopy
zip azcopy-v10.x.x-windows-amd64.zip azcopy.exe
```

### Distribution

- Upload binaries to GitHub Releases
- Update package manager repositories
- Update documentation and changelog
- Announce release in appropriate channels

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Azure Storage Documentation](https://docs.microsoft.com/en-us/azure/storage/)
- [Project Repository](https://github.com/Azure/azure-storage-azcopy)
- [Issue Tracker](https://github.com/Azure/azure-storage-azcopy/issues)
- [Contributing Guidelines](CONTRIBUTING.md)
- [Process Locking Documentation](docs/ProcessLocking.md)