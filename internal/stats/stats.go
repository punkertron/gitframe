package stats

type CommitInfo struct {
	Author    string
	Committer string
	Lines     int // Lines attributed to this commit in this file.
}

type FileStatistics struct {
	Path    string
	Commits map[string]*CommitInfo // Commit hash -> contribution.
}

type PersonStatistics struct {
	Lines   int
	Commits map[string]struct{}
	Files   map[string]struct{}
}

type Result struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}
