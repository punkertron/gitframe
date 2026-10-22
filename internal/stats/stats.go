package stats

type PersonStatistics struct {
	Lines   int
	Commits map[string]struct{}
	Files   map[string]struct{}
}

type Result struct {
	Name    string
	Lines   int
	Commits int
	Files   int
}
