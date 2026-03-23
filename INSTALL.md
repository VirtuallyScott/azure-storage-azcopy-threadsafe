# Quick Installation Guide

## Install Pre-built Binary

### Linux/macOS
```bash
# Download and install azcopy
curl -sSL https://aka.ms/downloadazcopy-v10-linux | tar -xzf -
sudo mv azcopy_linux_amd64_*/azcopy /usr/local/bin/
```

### macOS (Homebrew)
```bash
brew install azcopy
```

### Windows (PowerShell)
```powershell
# Download and extract
Invoke-WebRequest -Uri "https://aka.ms/downloadazcopy-v10-windows" -OutFile "azcopy.zip"
Expand-Archive -Path "azcopy.zip" -DestinationPath "."
Move-Item ".\azcopy_windows_amd64_*\azcopy.exe" "$env:USERPROFILE\bin\"
```

## Build from Source

### Prerequisites
- Go 1.21 or later
- Git

### Quick Build
```bash
git clone https://github.com/Azure/azure-storage-azcopy.git
cd azure-storage-azcopy
go build -o azcopy
sudo mv azcopy /usr/local/bin/
```

### Verify Installation
```bash
azcopy --version
azcopy login
```

## Process Locking Features (New)

For Unix/Linux/macOS systems, azcopy now includes robust process-level locking that prevents conflicts when multiple azcopy instances run simultaneously.

**Key Benefits:**
- Safe concurrent operation
- Automatic lock cleanup
- No configuration required

**Requirements:**
- Unix/Linux/macOS (Windows uses thread-level sync only)
- `github.com/gofrs/flock` dependency (auto-installed)

For detailed information, see [BUILD.md](BUILD.md) and [docs/ProcessLocking.md](docs/ProcessLocking.md).