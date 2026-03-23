# azcopy Thread-Safe Edition

This is a thread-safe and process-safe version of Microsoft's azcopy tool, enhanced with robust inter-process locking capabilities for Unix/Linux/macOS systems.

## 🚀 Key Features

### Multi-Process Safety
- **Process-level locking** using `flock` on Unix/Linux/macOS
- **Thread-level synchronization** fallback for Windows  
- **Concurrent operation support** - run multiple azcopy instances safely
- **Automatic cleanup** on process termination

### Protected Operations
- Job part plan file creation and mapping
- Job management and lifecycle operations
- OAuth token caching and credential management
- Memory mapped file operations

## 🏗️ Automated Builds

This repository includes GitHub Actions workflows that automatically build azcopy for multiple platforms:

### Supported Platforms
- **Linux AMD64** (`azcopy-linux-amd64`)
- **macOS ARM64** (`azcopy-macos-arm64`) - Apple Silicon
- **Windows AMD64** (`azcopy-windows-amd64.exe`)

### Build Features
- ✅ **Cross-platform compilation** with optimized static binaries
- ✅ **Automated testing** on all platforms
- ✅ **Build artifacts** stored for 30 days
- ✅ **Release automation** for tagged versions
- ✅ **Build summaries** with detailed information

### Download Pre-built Binaries

You can download the latest builds from the [Actions tab](../../actions) or [Releases page](../../releases).

## 📦 Installation

### Quick Install (Pre-built)
```bash
# Download latest artifacts from GitHub Actions
# Or install from releases for tagged versions
```

### Build from Source
```bash
git clone https://github.com/VirtuallyScott/azure-storage-azcopy-threadsafe.git
cd azure-storage-azcopy-threadsafe
go build -o azcopy
```

See [BUILD.md](BUILD.md) for detailed build instructions and [INSTALL.md](INSTALL.md) for quick installation guide.

## 🔒 Process Locking

The enhanced locking system provides:

- **File-based locking** using `github.com/gofrs/flock`
- **Resource isolation** prevents concurrent access conflicts  
- **Graceful fallback** if locking fails
- **Platform-aware implementation** (Unix vs Windows)

### Usage Example
```bash
# These operations are now safe to run concurrently:
azcopy copy source1/ dest1/ &
azcopy copy source2/ dest2/ &
azcopy sync source3/ dest3/ &
wait
```

For detailed information, see [docs/ProcessLocking.md](docs/ProcessLocking.md).

## 🧪 Testing

Run the test suite:
```bash
go test ./...
```

## 📋 Requirements

- **Go 1.21+** for building from source
- **Unix/Linux/macOS** for full process-level locking
- **Windows** supported with thread-level synchronization

## 📖 Documentation

- [BUILD.md](BUILD.md) - Comprehensive build guide
- [INSTALL.md](INSTALL.md) - Quick installation instructions  
- [docs/ProcessLocking.md](docs/ProcessLocking.md) - Process locking details

## 🤝 Contributing

This is a fork of Microsoft's azcopy with thread-safety enhancements. For the original project:
- [Original azcopy repository](https://github.com/Azure/azure-storage-azcopy)
- [Original documentation](https://docs.microsoft.com/en-us/azure/storage/common/storage-use-azcopy-v10)

## 📄 License

This project maintains the same license as the original azcopy project.

---

**Note**: This thread-safe version is designed for environments where multiple azcopy processes need to run concurrently without conflicts. For single-process usage, the original azcopy provides the same functionality.

:eight_spoked_asterisk: [Transfer data with AzCopy and Google GCP buckets](https://docs.microsoft.com/en-us/azure/storage/common/storage-use-azcopy-google-cloud)

:eight_spoked_asterisk: [Use data transfer tools in Azure Stack Hub Storage](https://docs.microsoft.com/en-us/azure-stack/user/azure-stack-storage-transfer)

:eight_spoked_asterisk: [Configure, optimize, and troubleshoot AzCopy](https://docs.microsoft.com/azure/storage/common/storage-use-azcopy-configure)

:eight_spoked_asterisk: [AzCopy WiKi](https://github.com/Azure/azure-storage-azcopy/wiki)

## Supported Operations

The general format of the AzCopy commands is: `azcopy [command] [arguments] --[flag-name]=[flag-value]`

* `bench` - Runs a performance benchmark by uploading or downloading test data to or from a specified destination

* `copy` - Copies source data to a destination location. The supported directions and forms of authorization are:
    Source | Destination
    --- | ---
    Local | Azure Blob (Microsoft Entra ID or SAS)
    Local | Azure Files SMB (Microsoft Entra ID or share/directory SAS)
    Local | Azure Files NFS (Microsoft Entra ID or share/directory SAS)
    Local | Azure Data Lake Storage (Microsoft Entra ID, SAS, or Shared Key)
    Azure Blob (Microsoft Entra ID or SAS) | Local
    Azure Files SMB (Microsoft Entra ID or share/directory SAS) | Local
    Azure Files NFS (Microsoft Entra ID or share/directory SAS) | Local
    Azure Data Lake Storage (Microsoft Entra ID, SAS, or Shared Key) | Local
    Azure Blob (Microsoft Entra ID, SAS, or public) | Azure Blob (Microsoft Entra ID or SAS)
    Azure Blob (Microsoft Entra ID, SAS, or public) | Azure Files SMB (Microsoft Entra ID or SAS)
    Azure Blob (Microsoft Entra ID or SAS) | Azure Data Lake Storage (Microsoft Entra ID or SAS)
    Azure Data Lake Storage (Microsoft Entra ID or SAS) | Azure Blob (Microsoft Entra ID or SAS)
    Azure Data Lake Storage (Microsoft Entra ID or SAS) | Azure Data Lake Storage (Microsoft Entra ID or SAS)
    Azure Files SMB (Microsoft Entra ID or SAS) | Azure Blob (Microsoft Entra ID or SAS)
    Azure Files SMB (Microsoft Entra ID or SAS) | Azure Files SMB (Microsoft Entra ID or SAS)
    Azure Files NFS (Microsoft Entra ID or SAS) | Azure Files NFS (Microsoft Entra ID or SAS)
    AWS S3 (Access Key) | Azure Block Blob (Microsoft Entra ID or SAS)
    Google Cloud Storage (Service Account Key) | Azure Block Blob (Microsoft Entra ID or SAS)

* `sync` - Replicate source to the destination location. The supported directions and forms of authorization are:
    Source | Destination
    --- | ---
    Local | Azure Blob (Microsoft Entra ID or SAS)
    Local | Azure File (Microsoft Entra ID or SAS)
    Azure Blob (Microsoft Entra ID or SAS) | Local
    Azure File (Microsoft Entra ID or SAS) | Local
    Azure Blob (Microsoft Entra ID or SAS) | Azure Blob (Microsoft Entra ID or SAS)
    Azure Blob (Microsoft Entra ID or SAS) | Azure File (Microsoft Entra ID or SAS)
    Azure Data Lake Storage (Microsoft Entra ID or SAS) | Azure Data Lake Storage (Microsoft Entra ID or SAS)
    Azure File (Microsoft Entra ID or SAS) | Azure Blob (Microsoft Entra ID or SAS)
    Azure File (SAS or public) | Azure File (SAS)

* `login` - Log in to Azure Active Directory (AD) to access Azure Storage resources.* `logout` - Log out to terminate access to Azure Storage resources.

* `list` - List the entities in a given resource

* `doc` - Generates documentation for the tool in Markdown format

* `env` - Shows the environment variables that you can use to configure the behavior of AzCopy.

* `help` - Help about any command

* `jobs` - Sub-commands related to managing jobs

* `load` - Sub-commands related to transferring data in specific formats

* `make` - Create a container or file share.

* `remove` - Delete blobs or files from an Azure storage account

## Find help from your command prompt

For convenience, consider adding the AzCopy directory location to your system path for ease of use. That way you can type `azcopy` from any directory on your system.

To see a list of commands, type `azcopy -h` and then press the ENTER key.

To learn about a specific command, just include the name of the command (For example: `azcopy list -h`).

![AzCopy command help example](readme-command-prompt.png)

If you choose not to add AzCopy to your path, you'll have to change directories to the location of your AzCopy executable and type `azcopy` or `.\azcopy` in Windows PowerShell command prompts.

## Frequently asked questions

### What is the difference between `sync` and `copy`?

* The `copy` command is a simple transferring operation. It scans/enumerates the source and attempts to transfer every single file/blob present on the source to the destination.
  The supported source/destination pairs are listed in the help message of the tool.

* On the other hand, `sync` scans/enumerates both the source, and the destination to find the incremental change.
  It makes sure that whatever is present in the source will be replicated to the destination. For `sync`,

* If your goal is to simply move some files, then `copy` is definitely the right command, since it offers much better performance.
  If the use case is to incrementally transfer data (files present only on source) then `sync` is the better choice, since only the modified/missing files will be transferred.
  Since `sync` enumerates both source and destination to find the incremental change, it is relatively slower as compared to `copy`

### Will `copy` overwrite my files?

By default, AzCopy will overwrite the files at the destination if they already exist. To avoid this behavior, please use the flag `--overwrite=false`.

### Will `sync` overwrite my files?

By default, AzCopy `sync` use last-modified-time to determine whether to transfer the same file present at both the source, and the destination.
i.e, If the source file is newer compared to the destination file, we overwrite the destination
You can change this default behaviour and overwrite files at the destination by using the flag `--mirror-mode=true`

### Will 'sync' delete files in the destination if they no longer exist in the source location?

By default, the 'sync' command doesn't delete files in the destination unless you use an optional flag with the command.
To learn more, see [Synchronize files](https://docs.microsoft.com/en-us/azure/storage/common/storage-use-azcopy-blobs-synchronize).

## How to contribute to AzCopy v10

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
