package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func runGitCommand(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("git command failed: %w: %s",
				err, strings.TrimSpace(string(exitErr.Stderr)))
		}
		return nil, fmt.Errorf("could not run git: %w", err)
	}
	return output, nil
}
