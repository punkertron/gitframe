package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/punkertron/gitframe/internal/stats"
)

func GetFileStatistics(ctx context.Context, repository, revision, file string) (stats.FileStatistics, error) {
	output, err := runGitCommand(ctx, "-C", repository, "blame", "--porcelain", revision, "--", file)
	if err != nil {
		return stats.FileStatistics{}, err
	}

	commits, err := parseBlame(ctx, output)
	if err != nil {
		return stats.FileStatistics{}, fmt.Errorf("parse git blame for %q: %w", file, err)
	}

	// Empty files have no blame records. Attribute them to their last change.
	if len(commits) == 0 {
		output, err := runGitCommand(
			ctx, "-C", repository,
			"log", "-1",
			"--format=%H%x00%aN%x00%cN",
			revision, "--", ":(literal)"+file,
		)
		if err != nil {
			return stats.FileStatistics{}, err
		}

		fields := strings.Split(strings.TrimSuffix(string(output), "\n"), "\x00")
		if len(fields) != 3 || fields[0] == "" {
			return stats.FileStatistics{}, fmt.Errorf("git log returned no attribution for empty file %q", file)
		}
		commits[fields[0]] = &stats.CommitInfo{Author: fields[1], Committer: fields[2]}
	}

	return stats.FileStatistics{Path: file, Commits: commits}, nil
}

func parseBlame(ctx context.Context, output []byte) (map[string]*stats.CommitInfo, error) {
	commits := make(map[string]*stats.CommitInfo)
	if len(output) == 0 {
		return commits, nil
	}

	// An empty currentHash means that the next line must be a record header.
	var currentHash string
	for line := range strings.SplitSeq(strings.TrimSuffix(string(output), "\n"), "\n") {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if currentHash == "" {
			fieldCount := 0
			for field := range strings.FieldsSeq(line) {
				if fieldCount == 0 {
					currentHash = field
				}
				fieldCount++
			}
			if fieldCount != 3 && fieldCount != 4 {
				return nil, fmt.Errorf("invalid blame header %q", line)
			}
			if _, ok := commits[currentHash]; !ok {
				currentHash = strings.Clone(currentHash)
				commits[currentHash] = &stats.CommitInfo{}
			}
			continue
		}

		info := commits[currentHash]
		switch {
		case strings.HasPrefix(line, "\t"):
			// Every source line counts, including a tab-only blank line.
			info.Lines++
			currentHash = ""
		case strings.HasPrefix(line, "author "):
			info.Author = strings.Clone(strings.TrimPrefix(line, "author "))
		case strings.HasPrefix(line, "committer "):
			info.Committer = strings.Clone(strings.TrimPrefix(line, "committer "))
		}
		// Other metadata is irrelevant to these statistics. Metadata omitted
		// from later records is already stored in the map for this hash.
	}
	if currentHash != "" {
		return nil, fmt.Errorf("incomplete blame record for commit %q", currentHash)
	}

	return commits, nil
}
