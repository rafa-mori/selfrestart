package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Installer struct{}

func NewInstaller() *Installer {
	return &Installer{}
}

func (i *Installer) IsInPath(target string) (bool, error) {
	path := os.Getenv("PATH")
	if path == "" {
		return false, fmt.Errorf("PATH environment variable not set")
	}

	pathSeparator := ":"
	if strings.Contains(path, ";") { // Windows uses semicolon
		pathSeparator = ";"
	}

	paths := strings.Split(path, pathSeparator)
	for _, dir := range paths {
		if dir == "" {
			continue
		}

		// Check if the target executable exists in this directory
		executable := filepath.Join(dir, target)
		if _, err := os.Stat(executable); err == nil {
			return true, nil
		}

		// Also check with .exe extension on Windows
		if pathSeparator == ";" {
			executable = filepath.Join(dir, target+".exe")
			if _, err := os.Stat(executable); err == nil {
				return true, nil
			}
		}
	}

	return false, nil
}

func (i *Installer) InstallGoUnix() (bool, error) {
	installed, err := i.IsInPath("go")
	if err != nil {
		return false, err
	}
	if installed {
		return true, nil
	}
	cmd := exec.Command("sh", "-c", "curl -sSL https://golang.org/dl/go1.20.5.linux-amd64.tar.gz | sudo tar -C /usr/local -xzf -")
	if err := cmd.Run(); err != nil {
		return false, err
	}
	installed, err = i.IsInPath("go")
	return installed, err
}
