package stats

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
