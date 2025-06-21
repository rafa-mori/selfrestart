package process

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

type ProcessManager struct{}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{}
}

func (pm *ProcessManager) GetCurrentPID() int {
	return os.Getpid()
}

func (pm *ProcessManager) KillCurrentProcess() error {
	pid := pm.GetCurrentPID()
	if pid <= 0 {
		return fmt.Errorf("invalid PID: %d", pid)
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("could not find process %d: %v", pid, err)
	}
	if err := process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("could not send interrupt to process %d: %v", pid, err)
	}
	if err := process.Release(); err != nil {
		return fmt.Errorf("could not release process %d: %v", pid, err)
	}
	deadline := time.Now().Add(10 * time.Second)
	sleepTicker := time.NewTicker(100 * time.Millisecond)
	beautySleep := 50 * time.Millisecond
	for {
		select {
		case <-time.After(deadline.Sub(time.Now())):
			if _, err := os.FindProcess(pid); err != nil {
				return fmt.Errorf("process %d did not terminate after interrupt", pid)
			}
		case <-sleepTicker.C:
			if _, err := os.FindProcess(pid); err != nil {
				return nil
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("process %d did not terminate after interrupt", pid)
			}
			if _, err := os.FindProcess(pid); err != nil {
				return nil
			}
		default:
			time.Sleep(beautySleep)
		}
	}
}

func (pm *ProcessManager) IsProcessRunning(pid int) (bool, error) {
	if pid <= 0 {
		return false, fmt.Errorf("invalid PID: %d", pid)
	}
	
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, fmt.Errorf("could not find process %d: %v", pid, err)
	}
	
	// On Unix systems, we can send signal 0 to check if process exists
	// On Windows, FindProcess always succeeds, so we use different approach
	if err := process.Signal(syscall.Signal(0)); err != nil {
		return false, nil // Process doesn't exist or we don't have permission
	}
	
	return true, nil
}
