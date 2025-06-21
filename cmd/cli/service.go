package cli

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/rafa-mori/selfrestart"
	gl "github.com/rafa-mori/selfrestart/logger"
	"github.com/spf13/cobra"
)

func ServiceCmdList() []*cobra.Command {
	return []*cobra.Command{
		startCommand(),
		restartCommand(),
		statusCommand(),
		checkCommand(),
	}
}

func startCommand() *cobra.Command {
	var debug bool
	var daemon bool

	var startCmd = &cobra.Command{
		Use: "start",
		Annotations: GetDescriptions([]string{
			"Start the SelfRestart service.",
			"This command starts the SelfRestart service and begins monitoring for restart signals.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			if debug {
				gl.SetDebug(true)
				gl.Log("debug", "Debug mode enabled")
			}

			sr := selfrestart.New()
			
			// Check if Go is installed
			if !sr.IsGolangInstalled() {
				gl.Log("error", "Go is not installed or not found in PATH")
				os.Exit(1)
			}

			platformInfo := sr.GetPlatformInfo()
			gl.Log("info", fmt.Sprintf("Platform: %s/%s", platformInfo.OS, platformInfo.Arch))
			gl.Log("info", fmt.Sprintf("Current PID: %d", sr.GetCurrentPID()))

			if daemon {
				gl.Log("info", "Starting in daemon mode...")
				// Setup signal handling for restart
				gl.Log("info", "Send SIGUSR1 to restart: kill -USR1 " + fmt.Sprintf("%d", sr.GetCurrentPID()))
			}

			gl.Log("success", "SelfRestart service started successfully")
			
			if daemon {
				// Keep running
				for {
					time.Sleep(5 * time.Second)
					gl.Log("debug", "Service running...")
				}
			}
		},
	}

	startCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	startCmd.Flags().BoolVarP(&daemon, "daemon", "", false, "Run as daemon")

	return startCmd
}

func restartCommand() *cobra.Command {
	var wait bool
	var pidFlag int

	var restartCmd = &cobra.Command{
		Use: "restart",
		Annotations: GetDescriptions([]string{
			"Restart the current process or a specified PID.",
			"This command restarts the current process safely using the SelfRestart mechanism.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			sr := selfrestart.New()

			var targetPID int
			if pidFlag > 0 {
				targetPID = pidFlag
				gl.Log("info", fmt.Sprintf("Attempting to restart process with PID: %d", targetPID))
				
				// Send SIGUSR1 signal to the target process
				if err := syscall.Kill(targetPID, syscall.SIGUSR1); err != nil {
					gl.Log("error", fmt.Sprintf("Failed to send restart signal to PID %d: %v", targetPID, err))
					os.Exit(1)
				}
				gl.Log("success", fmt.Sprintf("Restart signal sent to PID %d", targetPID))
			} else {
				targetPID = sr.GetCurrentPID()
				gl.Log("info", fmt.Sprintf("Restarting current process (PID: %d)", targetPID))
				
				if err := sr.Restart(); err != nil {
					gl.Log("error", fmt.Sprintf("Failed to restart: %v", err))
					os.Exit(1)
				}
				
				if wait {
					gl.Log("info", "Waiting for restart to complete...")
					time.Sleep(2 * time.Second)
				}
				
				gl.Log("success", "Process restart initiated")
				os.Exit(0)
			}
		},
	}

	restartCmd.Flags().BoolVarP(&wait, "wait", "w", false, "Wait for restart to complete")
	restartCmd.Flags().IntVarP(&pidFlag, "pid", "p", 0, "PID of process to restart")

	return restartCmd
}

func statusCommand() *cobra.Command {
	var pidFlag int

	var statusCmd = &cobra.Command{
		Use: "status",
		Annotations: GetDescriptions([]string{
			"Check the status of a process.",
			"This command checks if a process is running and provides information about it.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			sr := selfrestart.New()

			var targetPID int
			if pidFlag > 0 {
				targetPID = pidFlag
			} else {
				targetPID = sr.GetCurrentPID()
			}

			gl.Log("info", fmt.Sprintf("Checking status of PID: %d", targetPID))

			running, err := sr.IsProcessRunning(targetPID)
			if err != nil {
				gl.Log("error", fmt.Sprintf("Error checking process status: %v", err))
				os.Exit(1)
			}

			if running {
				gl.Log("success", fmt.Sprintf("Process %d is running", targetPID))
			} else {
				gl.Log("warn", fmt.Sprintf("Process %d is not running", targetPID))
			}

			platformInfo := sr.GetPlatformInfo()
			gl.Log("info", fmt.Sprintf("Platform: %s/%s", platformInfo.OS, platformInfo.Arch))
		},
	}

	statusCmd.Flags().IntVarP(&pidFlag, "pid", "p", 0, "PID to check (default: current process)")

	return statusCmd
}

func checkCommand() *cobra.Command {
	var checkCmd = &cobra.Command{
		Use: "check",
		Aliases: []string{"health"},
		Annotations: GetDescriptions([]string{
			"Check system requirements and Go installation.",
			"This command verifies that all requirements are met for SelfRestart to function properly.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			sr := selfrestart.New()

			gl.Log("info", "Checking system requirements...")

			// Check Go installation
			if sr.IsGolangInstalled() {
				gl.Log("success", "✅ Go is installed and accessible")
			} else {
				gl.Log("error", "❌ Go is not installed or not found in PATH")
			}

			// Check platform support
			platformInfo := sr.GetPlatformInfo()
			gl.Log("info", fmt.Sprintf("Platform: %s/%s", platformInfo.OS, platformInfo.Arch))

			supportedPlatforms := []string{"linux", "darwin", "windows"}
			supported := false
			for _, p := range supportedPlatforms {
				if platformInfo.OS == p {
					supported = true
					break
				}
			}

			if supported {
				gl.Log("success", "✅ Platform is supported")
			} else {
				gl.Log("warn", "⚠️  Platform may not be fully supported")
			}

			// Check permissions
			if _, err := os.Stat("/tmp"); err != nil {
				gl.Log("error", "❌ Cannot access /tmp directory")
			} else {
				gl.Log("success", "✅ Temporary directory is accessible")
			}

			gl.Log("success", "System check completed")
		},
	}

	return checkCmd
}
