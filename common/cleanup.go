package common

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	cleanupHandlers []func()
	cleanupMux      sync.Mutex
	signalHandlerInitialized sync.Once
)

// RegisterCleanupHandler registers a function to be called on process exit
func RegisterCleanupHandler(handler func()) {
	cleanupMux.Lock()
	defer cleanupMux.Unlock()
	
	cleanupHandlers = append(cleanupHandlers, handler)
	
	// Initialize signal handler on first registration
	signalHandlerInitialized.Do(func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		
		go func() {
			<-c
			performCleanup()
			os.Exit(1)
		}()
	})
}

// performCleanup executes all registered cleanup handlers
func performCleanup() {
	cleanupMux.Lock()
	defer cleanupMux.Unlock()
	
	for _, handler := range cleanupHandlers {
		if handler != nil {
			func() {
				defer func() {
					if r := recover(); r != nil {
						// Log the panic but continue with other cleanup handlers
						GetLifecycleMgr().Info(fmt.Sprintf("Cleanup handler panicked: %v", r))
					}
				}()
				handler()
			}()
		}
	}
}

// PerformFinalCleanup should be called before the application exits
// This is a public function that can be called from main() or other exit points
func PerformFinalCleanup() {
	performCleanup()
}

// init automatically registers the process lock manager cleanup
func init() {
	RegisterCleanupHandler(func() {
		if globalLockManager != nil {
			globalLockManager.Cleanup()
		}
	})
}