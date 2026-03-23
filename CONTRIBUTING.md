# Contributing to azcopy Thread-Safe Edition

We welcome contributions to make azcopy more robust and thread-safe! This guide will help you contribute effectively.

## Getting Started

### Prerequisites
- Go 1.21 or later
- Git with git-flow extension
- Understanding of concurrent programming concepts

### Repository Setup
```bash
# Fork the repository on GitHub
git clone git@github.com:YourUsername/azure-storage-azcopy-threadsafe.git
cd azure-storage-azcopy-threadsafe

# Set up git flow (if not already done)
git flow init -d
```

## Development Workflow

We use **Git Flow** for organized development. Please familiarize yourself with [our Git Flow documentation](docs/GitFlow.md).

### For New Features

1. **Start Feature Branch**:
   ```bash
   git flow feature start your-feature-name
   ```

2. **Develop Your Feature**:
   - Write clean, well-documented code
   - Add tests for your changes
   - Ensure thread-safety for concurrent operations
   - Test on multiple platforms (Unix/Linux/macOS/Windows)

3. **Test Thoroughly**:
   ```bash
   # Run unit tests
   go test ./...
   
   # Test concurrent operations (Unix/Linux/macOS)
   ./azcopy copy source1/ dest1/ &
   ./azcopy copy source2/ dest2/ &
   wait
   
   # Check for lock issues
   ls -la ~/.azcopy/.locks/
   ```

4. **Finish Feature**:
   ```bash
   git flow feature finish your-feature-name
   git push origin develop
   ```

### For Bug Fixes

#### Non-Critical Bugs
Follow the same process as features above.

#### Critical Bugs (Hotfixes)
```bash
git flow hotfix start hotfix-name
# Make minimal changes to fix the issue
# Test thoroughly
git flow hotfix finish hotfix-name
git push origin main && git push origin develop && git push --tags
```

## Code Guidelines

### Thread Safety Requirements
- All new code must be thread-safe and process-safe
- Use proper locking mechanisms (our `ProcessLockManager`)
- Avoid race conditions in file operations
- Test concurrent scenarios

### Code Style
```go
// Use meaningful names
func (plm *ProcessLockManager) AcquireLock(resourceName string, timeout time.Duration) (func(), error) {
    // Proper error handling
    if resourceName == "" {
        return nil, fmt.Errorf("resource name cannot be empty")
    }
    
    // Clear logic flow
    unlock, err := plm.acquireLockInternal(resourceName, timeout)
    if err != nil {
        return nil, fmt.Errorf("failed to acquire lock for %s: %w", resourceName, err)
    }
    
    return unlock, nil
}
```

### Documentation Requirements
- Document all public functions and types
- Include usage examples for complex features
- Update relevant markdown files (BUILD.md, README.md, etc.)
- Add architecture documentation for significant changes

### Testing Requirements
```go
func TestProcessLockingConcurrency(t *testing.T) {
    // Test concurrent access
    const numGoroutines = 10
    const resourceName = "test-resource"
    
    var wg sync.WaitGroup
    var successCount int64
    
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            lockMgr := GetProcessLockManager()
            unlock, err := lockMgr.TryAcquireLock(resourceName)
            if err == nil {
                atomic.AddInt64(&successCount, 1)
                unlock()
            }
        }()
    }
    
    wg.Wait()
    assert.Equal(t, int64(1), successCount, "Only one goroutine should acquire the lock")
}
```

## Platform-Specific Considerations

### Unix/Linux/macOS Features
- Full process-level locking with flock
- Signal handling for cleanup
- Proper file permission handling

### Windows Compatibility
- Thread-level synchronization fallback
- Windows-specific file handling
- Error message compatibility

## Pull Request Process

### Before Creating a PR
1. Ensure your branch is up to date with develop
2. Run all tests and ensure they pass
3. Build on multiple platforms if possible
4. Update documentation

### PR Description Template
```markdown
## Description
Brief description of the changes.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing on Unix/Linux/macOS
- [ ] Manual testing on Windows
- [ ] Concurrent operation testing

## Thread Safety
- [ ] Changes are thread-safe
- [ ] Process-level locking used appropriately
- [ ] No race conditions introduced
- [ ] Proper cleanup mechanisms

## Documentation
- [ ] Code comments added/updated
- [ ] Documentation files updated
- [ ] Examples provided if needed

## Checklist
- [ ] My code follows the project's style guidelines
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
```

### Review Process
1. **Automated Checks**: GitHub Actions will build and test your changes
2. **Code Review**: Maintainers will review for:
   - Thread safety correctness
   - Code quality and style
   - Documentation completeness
   - Test coverage
3. **Testing**: Manual testing on multiple platforms
4. **Approval**: At least one maintainer approval required

## Reporting Issues

### Bug Reports
Include:
- Operating system and version
- Go version
- azcopy version/commit hash
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs or error messages
- Whether issue occurs with concurrent operations

### Feature Requests
Include:
- Clear description of the desired functionality
- Use case and motivation
- Consideration for thread safety
- Potential implementation approach

## Development Environment

### Required Tools
```bash
# Install git-flow
brew install git-flow          # macOS
sudo apt-get install git-flow  # Ubuntu

# Install development tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### IDE Setup
Recommended VSCode extensions:
- Go extension
- GitLens
- GitHub Pull Requests and Issues

### Environment Variables
```bash
export GOOS=linux    # For cross-compilation testing
export GOARCH=amd64
export CGO_ENABLED=0 # For static builds
```

## Contact

- **Issues**: [GitHub Issues](https://github.com/VirtuallyScott/azure-storage-azcopy-threadsafe/issues)
- **Discussions**: [GitHub Discussions](https://github.com/VirtuallyScott/azure-storage-azcopy-threadsafe/discussions)

## License

By contributing to this project, you agree that your contributions will be licensed under the same license as the project.

---

Thank you for contributing to azcopy thread-safe edition! Your efforts help make file transfers safer and more reliable for concurrent environments.