package git

import (
	"strings"
)

func ResolveRevision(path, commit string) (string, error) {
	output, err := runGitCommand("-C", path, "rev-parse", "--verify", "--end-of-options", commit+"^{commit}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
