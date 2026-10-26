package sort

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/punkertron/gitfame/internal/stats"
)

type Comparator func(x, y stats.Result) int

func SortData(statistics map[string]*stats.PersonStatistics, orderBy string) ([]stats.Result, error) {
	res := make([]stats.Result, 0, len(statistics))

	for name, stat := range statistics {
		res = append(res, stats.Result{
			Name:    name,
			Lines:   stat.Lines,
			Commits: len(stat.Commits),
			Files:   len(stat.Files),
		})
	}

	byLines := func(x, y stats.Result) int {
		return cmp.Compare(y.Lines, x.Lines)
	}

	byCommits := func(x, y stats.Result) int {
		return cmp.Compare(y.Commits, x.Commits)
	}

	byFiles := func(x, y stats.Result) int {
		return cmp.Compare(y.Files, x.Files)
	}

	byName := func(x, y stats.Result) int {
		return cmp.Compare(x.Name, y.Name)
	}

	orders := map[string][]Comparator{
		"":        {byLines, byCommits, byFiles, byName}, // default
		"lines":   {byLines, byCommits, byFiles, byName},
		"commits": {byCommits, byLines, byFiles, byName},
		"files":   {byFiles, byLines, byCommits, byName},
	}

	comparators, exists := orders[orderBy]
	if !exists {
		return nil, fmt.Errorf("unsupported OrderBy flag: %s", orderBy)
	}

	slices.SortFunc(res, func(x, y stats.Result) int {
		for _, comp := range comparators {
			if c := comp(x, y); c != 0 {
				return c
			}
		}
		return 0
	})
	return res, nil
}
