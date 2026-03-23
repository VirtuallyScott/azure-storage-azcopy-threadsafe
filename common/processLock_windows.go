//go:build windows
// +build windows

package common

import (
	"sync"
	"time"
)

// ProcessLockManager provides a no-op implementation for Windows
// since Windows file locking behavior is different and the user requested Unix-specific locking
type ProcessLockManager struct {
	mu     sync.RWMutex
	folder string
}

var (
	globalLockManager *ProcessLockManager
	lockManagerOnce   sync.Once
)

// GetProcessLockManager returns the global process lock manager instance
func GetProcessLockManager() *ProcessLockManager {
	lockManagerOnce.Do(func() {
		globalLockManager = &ProcessLockManager{}
		GetLifecycleMgr().Info("Process-level locking is not implemented on Windows - using thread-level synchronization only")
	})
	return globalLockManager
}

// AcquireLock is a no-op on Windows - returns immediately
func (plm *ProcessLockManager) AcquireLock(resourceName string, timeout time.Duration) (func(), error) {
	// No-op unlock function
	return func() {}, nil
}

// TryAcquireLock is a no-op on Windows - returns immediately
func (plm *ProcessLockManager) TryAcquireLock(resourceName string) (func(), error) {
	return func() {}, nil
}

// AcquireJobLock is a no-op on Windows
func (plm *ProcessLockManager) AcquireJobLock(jobID JobID, timeout time.Duration) (func(), error) {
	return func() {}, nil
}

// AcquirePlanFileLock is a no-op on Windows
func (plm *ProcessLockManager) AcquirePlanFileLock(planFileName string, timeout time.Duration) (func(), error) {
	return func() {}, nil
}

// AcquireGlobalLock is a no-op on Windows
func (plm *ProcessLockManager) AcquireGlobalLock(operation string, timeout time.Duration) (func(), error) {
	return func() {}, nil
}

// Cleanup is a no-op on Windows
func (plm *ProcessLockManager) Cleanup() {}

// WithLock executes the function without locking on Windows
func (plm *ProcessLockManager) WithLock(resourceName string, timeout time.Duration, fn func() error) error {
	return fn()
}

// WithJobLock executes the function without locking on Windows
func (plm *ProcessLockManager) WithJobLock(jobID JobID, timeout time.Duration, fn func() error) error {
	return fn()
}

// WithPlanFileLock executes the function without locking on Windows
func (plm *ProcessLockManager) WithPlanFileLock(planFileName string, timeout time.Duration, fn func() error) error {
	return fn()
}

// Default timeout for lock operations (unused on Windows)
const DefaultLockTimeout = 30 * time.Second
