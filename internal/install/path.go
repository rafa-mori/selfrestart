package install

import (
	"os"
	"os/user"
	"path/filepath"
)

func GetUserShellRC() string {
	home, _ := os.UserHomeDir()
	shell := os.Getenv("SHELL")
	shellRC := ".bashrc"
	if shell != "" && len(shell) > 2 && shell[len(shell)-2:] == "sh" {
		shellRC = ".zshrc"
	}
	return filepath.Join(home, shellRC)
}

func GetCurrentUser() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.Username
}
