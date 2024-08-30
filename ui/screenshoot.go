package ui

import (
	"fmt"
	"os/exec"
)

func screenshoot(path string) error {
	args := []string{"-s", "-m", "10", path}
	cmd := exec.Command("maim", args...)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("[screenshoot] maim finished with error: %w", err)
	}

	return nil
}
