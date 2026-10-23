package git

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

func GetFiles(ctx context.Context, path, commit string) ([]string, error) {
	output, err := runGitCommand(ctx, "-C", path, "ls-tree", "-r", "-z", "--format=%(objecttype)%x09%(path)", commit)
	if err != nil {
		return nil, err
	}

	if len(output) == 0 {
		return nil, nil
	}

	files := []string{}
	for record := range strings.SplitSeq(strings.TrimSuffix(string(output), "\000"), "\000") {
		objectType, filename, ok := strings.Cut(record, "\t")
		if !ok {
			return nil, fmt.Errorf("invalid ls-tree record %q", record)
		}

		if strings.HasPrefix(filename, "\"") {
			filename, err = strconv.Unquote(filename)
			if err != nil {
				return nil, fmt.Errorf("decode ls-tree filename: %w", err)
			}
		}

		if objectType == "blob" {
			files = append(files, strings.Clone(filename))
		}
	}
	return files, nil
}
