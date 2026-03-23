# Process-Level Locking for azcopy

## Overview

This implementation provides robust inter-process locking for azcopy operations on Unix-like systems (Linux and macOS) using the `github.com/gofrs/flock` library. The locking mechanism prevents data corruption and race conditions when multiple azcopy processes run concurrently.

## Features

### Supported Platforms
- **Linux**: Full process-level locking using flock
- **macOS**: Full process-level locking using flock  
- **Windows**: No-op implementation (falls back to existing thread-level synchronization)

### Protected Operations

1. **Job Part Plan Files**
   - Creation of new plan files
   - Memory mapping of existing plan files
   - Prevents concurrent access that could corrupt job state

2. **Job Management**
   - Job creation and initialization
   - Job cleanup operations
   - Job lifecycle state changes

3. **Credential Management**
   - OAuth token caching operations
   - Prevents token corruption during concurrent authentication

4. **Memory Mapped Files**
   - MMF creation and mapping operations
   - Protects against race conditions in file system operations

## Architecture

### ProcessLockManager

The `ProcessLockManager` provides the core locking functionality:

```go
// Get the global lock manager instance
lockMgr := common.GetProcessLockManager()

// Acquire a lock with timeout
unlock, err := lockMgr.AcquireLock("resource-name", 30*time.Second)
if err != nil {
    return err
}
defer unlock()

// Protected operation here
```

### Lock Types

1. **Job Locks**: `AcquireJobLock(jobID, timeout)`
   - Protects job-specific operations
   - Named: `job-{jobID}`

2. **Plan File Locks**: `AcquirePlanFileLock(fileName, timeout)`
   - Protects plan file operations
   - Named: `plan-{filename}`

3. **Global Locks**: `AcquireGlobalLock(operation, timeout)`
   - Protects global operations
   - Named: `global-{operation}`

4. **Resource Locks**: `AcquireLock(name, timeout)`
   - General-purpose locking
   - Custom naming

### Lock Storage

Lock files are stored in:
- Primary: `{AzcopyJobPlanFolder}/.locks/`
- Fallback: System temporary directory

Lock files use the naming convention: `{resource-name}.lock`

## Usage Examples

### Job Operations
```go
func processJob(jobID common.JobID) error {
    lockMgr := common.GetProcessLockManager()
    
    return lockMgr.WithJobLock(jobID, common.DefaultLockTimeout, func() error {
        // Job processing logic here
        return nil
    })
}
```

### Plan File Operations
```go
func createPlanFile(fileName string) error {
    lockMgr := common.GetProcessLockManager()
    
    return lockMgr.WithPlanFileLock(fileName, common.DefaultLockTimeout, func() error {
        // Plan file creation logic here
        return nil
    })
}
```

### Global Operations
```go
func performGlobalOperation() error {
    lockMgr := common.GetProcessLockManager()
    
    return lockMgr.WithLock("global-cleanup", common.DefaultLockTimeout, func() error {
        // Global operation logic here
        return nil
    })
}
```

## Error Handling

The locking system is designed to be resilient:

1. **Lock Acquisition Failures**: Operations continue with warning logs, falling back to thread-level synchronization
2. **Timeout Handling**: Configurable timeouts prevent indefinite blocking
3. **Cleanup on Exit**: Signal handlers ensure locks are released on process termination
4. **Panic Recovery**: Cleanup handlers are protected against panics

## Configuration

### Default Timeout
```go
const DefaultLockTimeout = 30 * time.Second
```

### Custom Timeouts
```go
// Short timeout for non-critical operations
unlock, err := lockMgr.AcquireLock("resource", 5*time.Second)

// Longer timeout for critical operations  
unlock, err := lockMgr.AcquireLock("resource", 60*time.Second)
```

## Implementation Details

### Lock File Management
- Lock files are automatically created and cleaned up
- Stale lock files are handled by the flock library
- Process termination automatically releases locks

### Thread Safety
- Internal operations are protected by RWMutex
- Safe for concurrent use within a single process
- Prevents lock leakage through careful resource management

### Fallback Behavior
- If process-level locking fails, operations continue with warnings
- Existing thread-level synchronization remains active
- Graceful degradation ensures functionality is maintained

## Monitoring and Debugging

### Logging
Lock operations generate informational logs:
```
Failed to acquire lock for resource: timeout
Successfully acquired lock for plan-job123--00001.steV20
Released lock for job-{job-id}
```

### Lock Directory
Check the lock directory for active locks:
```bash
ls -la {AzcopyJobPlanFolder}/.locks/
```

### Troubleshooting

1. **Permission Issues**: Ensure azcopy has write access to the job plan folder
2. **Stale Locks**: The flock library automatically handles stale locks from terminated processes
3. **High Contention**: Consider reducing concurrency or increasing timeouts
4. **Disk Space**: Lock files are small but ensure adequate disk space

## Migration Notes

### Upgrading from Previous Versions
- No configuration changes required
- Existing job files are compatible
- Process locking is automatically enabled on Unix systems

### Performance Impact
- Minimal overhead for lock acquisition/release
- No impact when no other processes are running
- Timeout-based blocking prevents indefinite waits

## Best Practices

1. **Use Appropriate Timeouts**: Balance between avoiding deadlocks and allowing sufficient time for operations
2. **Handle Errors Gracefully**: Always check for lock acquisition failures and log appropriately
3. **Clean Resource Names**: Use descriptive, filesystem-safe names for custom locks
4. **Avoid Nested Locks**: Prevent deadlocks by avoiding complex lock hierarchies
5. **Monitor Lock Directory**: Periodically check for excessive lock files that might indicate issues

## Future Enhancements

- Windows support using named mutexes or other Windows-specific mechanisms
- Lock priority queuing for high-priority operations  
- Metrics collection for lock contention analysis
- Administrative tools for lock management and debugging