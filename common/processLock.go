//go:build !windows
// +build !windows

package common

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gofrs/flock"
)

// ProcessLockManager provides robust inter-process locking for azcopy operations
type ProcessLockManager struct {
	mu     sync.RWMutex
	locks  map[string]*flock.Flock
	folder string
}

var (
	globalLockManager *ProcessLockManager
	lockManagerOnce   sync.Once
)

// GetProcessLockManager returns the global process lock manager instance
func GetProcessLockManager() *ProcessLockManager {
	lockManagerOnce.Do(func() {
		lockDir := filepath.Join(AzcopyJobPlanFolder, ".locks")
		if err := os.MkdirAll(lockDir, DEFAULT_FILE_PERM); err != nil {
			// If we can't create the lock directory, log but continue
			// This allows azcopy to run even if locking fails
			GetLifecycleMgr().Info(fmt.Sprintf("Failed to create lock directory %s: %v", lockDir, err))
			lockDir = os.TempDir()
		}

		globalLockManager = &ProcessLockManager{
			locks:  make(map[string]*flock.Flock),
			folder: lockDir,
		}
	})
	return globalLockManager
}

// AcquireLock acquires an exclusive process lock for the given resource
// Returns a function to release the lock and an error if locking failed
func (plm *ProcessLockManager) AcquireLock(resourceName string, timeout time.Duration) (func(), error) {
	plm.mu.Lock()
	defer plm.mu.Unlock()

	lockFile := filepath.Join(plm.folder, fmt.Sprintf("%s.lock", resourceName))

	// Check if we already have a lock for this resource
	if existingLock, exists := plm.locks[resourceName]; exists {
		// Try to acquire the lock (it might have been released by another thread)
		locked, err := existingLock.TryLock()
		if err != nil || !locked {
			return nil, fmt.Errorf("resource %s already locked: %w", resourceName, err)
		}
		// Return the unlock function
		return func() { plm.releaseLock(resourceName) }, nil
	}

	// Create new lock
	lock := flock.New(lockFile)

	// Try to acquire with timeout
	ctx := make(chan error, 1)
	go func() {
		err := lock.Lock()
		ctx <- err
	}()

	select {
	case err := <-ctx:
		if err != nil {
			return nil, fmt.Errorf("failed to acquire lock for %s: %w", resourceName, err)
		}
		plm.locks[resourceName] = lock
		return func() { plm.releaseLock(resourceName) }, nil
	case <-time.After(timeout):
		// Cleanup the goroutine
		go func() { <-ctx }()
		return nil, fmt.Errorf("timeout acquiring lock for %s", resourceName)
	}
}

// TryAcquireLock attempts to acquire a lock without blocking
func (plm *ProcessLockManager) TryAcquireLock(resourceName string) (func(), error) {
	plm.mu.Lock()
	defer plm.mu.Unlock()

	lockFile := filepath.Join(plm.folder, fmt.Sprintf("%s.lock", resourceName))

	// Check if we already have a lock for this resource
	if existingLock, exists := plm.locks[resourceName]; exists {
		locked, err := existingLock.TryLock()
		if err != nil || !locked {
			return nil, fmt.Errorf("resource %s already locked: %w", resourceName, err)
		}
		return func() { plm.releaseLock(resourceName) }, nil
	}

	// Create new lock and try to acquire immediately
	lock := flock.New(lockFile)
	locked, err := lock.TryLock()
	if err != nil || !locked {
		return nil, fmt.Errorf("failed to acquire lock for %s: %w", resourceName, err)
	}

	plm.locks[resourceName] = lock
	return func() { plm.releaseLock(resourceName) }, nil
}

// releaseLock releases the lock for the given resource
// This is called by the unlock function returned by AcquireLock
func (plm *ProcessLockManager) releaseLock(resourceName string) {
	plm.mu.Lock()
	defer plm.mu.Unlock()

	if lock, exists := plm.locks[resourceName]; exists {
		lock.Unlock()
		delete(plm.locks, resourceName)
	}
}

// AcquireJobLock acquires a lock for a specific job ID
func (plm *ProcessLockManager) AcquireJobLock(jobID JobID, timeout time.Duration) (func(), error) {
	return plm.AcquireLock(fmt.Sprintf("job-%s", jobID.String()), timeout)
}

// AcquirePlanFileLock acquires a lock for a specific plan file
func (plm *ProcessLockManager) AcquirePlanFileLock(planFileName string, timeout time.Duration) (func(), error) {
	// Clean the filename to make it safe for filesystem
	safeName := filepath.Base(planFileName)
	return plm.AcquireLock(fmt.Sprintf("plan-%s", safeName), timeout)
}

// AcquireGlobalLock acquires a global lock for operations that need exclusive access
func (plm *ProcessLockManager) AcquireGlobalLock(operation string, timeout time.Duration) (func(), error) {
	return plm.AcquireLock(fmt.Sprintf("global-%s", operation), timeout)
}

// Cleanup releases all active locks and cleans up resources
func (plm *ProcessLockManager) Cleanup() {
	plm.mu.Lock()
	defer plm.mu.Unlock()

	for resourceName, lock := range plm.locks {
		lock.Unlock()
		delete(plm.locks, resourceName)
	}
}

// WithLock executes a function while holding a process lock
func (plm *ProcessLockManager) WithLock(resourceName string, timeout time.Duration, fn func() error) error {
	unlock, err := plm.AcquireLock(resourceName, timeout)
	if err != nil {
		return err
	}
	defer unlock()
	return fn()
}

// WithJobLock executes a function while holding a job-specific lock
func (plm *ProcessLockManager) WithJobLock(jobID JobID, timeout time.Duration, fn func() error) error {
	unlock, err := plm.AcquireJobLock(jobID, timeout)
	if err != nil {
		return err
	}
	defer unlock()
	return fn()
}

// WithPlanFileLock executes a function while holding a plan file lock
func (plm *ProcessLockManager) WithPlanFileLock(planFileName string, timeout time.Duration, fn func() error) error {
	unlock, err := plm.AcquirePlanFileLock(planFileName, timeout)
	if err != nil {
		return err
	}
	defer unlock()
	return fn()
}

// Default timeout for lock operations
const DefaultLockTimeout = 30 * time.Second
