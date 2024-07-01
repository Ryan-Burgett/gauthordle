package git

import "github.com/josephnaberhaus/gauthordle/internal/command"

func IsGitInstalled() bool {
	_, err := command.Run("git", "--version")
	if err != nil {
		return false
	}

	return true
}

func IsInGitRepo() bool {
	_, err := command.Run("git", "status")
	if err != nil {
		return false
	}

	return true
}
