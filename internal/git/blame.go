package git

import (
	"fmt"
	"strings"
)

type CommitInfo struct {
	Author    string
	Committer string
	Lines     int // Lines attributed to this commit in this file.
}

type FileStatistics struct {
	Path    string
	Commits map[string]*CommitInfo // Commit hash -> contribution.
}

func GetFileStatistics(repository, revision, file string) (FileStatistics, error) {
	output, err := runGitCommand("-C", repository, "blame", "--porcelain", revision, "--", file)
	if err != nil {
		return FileStatistics{}, err
	}

	commits, err := parseBlame(output)
	if err != nil {
		return FileStatistics{}, fmt.Errorf("parse git blame for %q: %w", file, err)
	}

	// Empty files have no blame records. Attribute them to their last change.
	if len(commits) == 0 {
		output, err := runGitCommand(
			"-C", repository,
			"log", "-1",
			"--format=%H%x00%aN%x00%cN",
			revision, "--", ":(literal)"+file,
		)
		if err != nil {
			return FileStatistics{}, err
		}

		fields := strings.Split(strings.TrimSuffix(string(output), "\n"), "\x00")
		if len(fields) != 3 || fields[0] == "" {
			return FileStatistics{}, fmt.Errorf("git log returned no attribution for empty file %q", file)
		}
		commits[fields[0]] = &CommitInfo{Author: fields[1], Committer: fields[2]}
	}

	return FileStatistics{Path: file, Commits: commits}, nil
}

func parseBlame(output []byte) (map[string]*CommitInfo, error) {
	commits := make(map[string]*CommitInfo)
	if len(output) == 0 {
		return commits, nil
	}

	// An empty currentHash means that the next line must be a record header.
	var currentHash string
	for line := range strings.SplitSeq(strings.TrimSuffix(string(output), "\n"), "\n") {
		if currentHash == "" {
			fields := strings.Fields(line)
			if len(fields) != 3 && len(fields) != 4 {
				return nil, fmt.Errorf("invalid blame header %q", line)
			}
			currentHash = fields[0]
			if _, ok := commits[currentHash]; !ok {
				currentHash = strings.Clone(currentHash)
				commits[currentHash] = &CommitInfo{}
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
