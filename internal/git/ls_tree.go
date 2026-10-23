package git

import (
	"fmt"
	"strings"
)

func GetFiles(path, commit string) ([]string, error) {
	output, err := runGitCommand("-C", path, "ls-tree", "-r", "-z", "--format=%(objecttype)%x09%(path)", commit)
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

		if objectType == "blob" {
			files = append(files, strings.Clone(filename))
		}
	}
	return files, nil
}
