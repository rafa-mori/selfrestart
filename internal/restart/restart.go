package restart

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Restarter struct{}

func NewRestarter() *Restarter {
	return &Restarter{}
}

func (r *Restarter) CreateAndExecRestartScript(oldPID int, binPath string) error {
	script := fmt.Sprintf(`#!/bin/sh
LOG="/tmp/selfrestart.log"
echo "[trap] Preparing restart..." >> $LOG
trap '
  echo "[trap] Old process finished. Trying to restart..." >> $LOG
  if [ -x "%s" ]; then
    echo "[trap] Executing new binary: %s" >> $LOG
    "%s" &
  else
    echo "[trap] New binary not found or not executable." >> $LOG
  fi
' EXIT

echo "[info] Waiting for process %d to finish..." >> $LOG
while kill -0 %d 2>/dev/null; do
  sleep 0.5
done

exit
`, binPath, binPath, binPath, oldPID, oldPID)

	tmpPath := filepath.Join(os.TempDir(), "restart_helper.sh")
	if err := os.WriteFile(tmpPath, []byte(script), 0755); err != nil {
		return fmt.Errorf("could not write restart script: %v", err)
	}

	cmd := exec.Command("sh", tmpPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not start restart script: %v", err)
	}

	prc := cmd.Process
	if err := prc.Release(); err != nil {
		return fmt.Errorf("could not detach restart process: %v", err)
	}

	if len(os.Args) > 1 && os.Args[1] == "--wait" {
		if err := cmd.Wait(); err != nil {
			return cmd.Start()
		}
	}

	return nil
}
