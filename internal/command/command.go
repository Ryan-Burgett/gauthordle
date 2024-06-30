package command

import (
	"fmt"
	"os"
	"os/exec"
)

func Run(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	c.Stderr = os.Stderr

	result, err := c.Output()
	if err != nil {
		return "", fmt.Errorf("error when running %s: %w", c.String(), err)
	}

	return string(result), nil
}
