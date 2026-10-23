package git

import (
	"context"
	"strings"
)

func ResolveRevision(ctx context.Context, path, commit string) (string, error) {
	output, err := runGitCommand(ctx, "-C", path, "rev-parse", "--verify", "--end-of-options", commit+"^{commit}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
