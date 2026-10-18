package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func GetFiles(path, commit string) ([]string, error) {
	cmd := exec.Command("git", "-C", path, "ls-tree", "-r", "-z", "--name-only", commit)

	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf(
				"git ls-tree failed: %w %s",
				err, strings.TrimSpace(string(exitErr.Stderr)),
			)
		}
		return nil, fmt.Errorf("could not run git ls-tree: %w", err)
	}

	if len(output) == 0 {
		return nil, nil
	}
	files := strings.TrimSuffix(string(output), "\000")
	return strings.Split(files, "\000"), nil
}
