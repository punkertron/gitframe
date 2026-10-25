package cli

import (
	"context"
	"io"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/punkertron/gitframe/internal/git"
	"github.com/punkertron/gitframe/internal/stats"
)

type fileResult struct {
	stat stats.FileStatistics
	file string
}

func collectStatistics(ctx context.Context, statusWriter io.Writer, opt *Options, revision string, files []string) (map[string]*stats.PersonStatistics, error) {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	people := make(map[string]*stats.PersonStatistics)

	var progressWG sync.WaitGroup
	done := make(chan struct{})
	var cnt atomic.Int64

	if opt.Progress && len(files) > 0 {
		progressWG.Add(1)
		go func() {
			defer progressWG.Done()
			printProgress(len(files), statusWriter, done, &cnt)
		}()
	}

	for result := range startWorkers(ctx, cancel, opt.Repository, revision, files) {
		addFileStatistics(people, result.file, result.stat, opt.UseCommitter)
		cnt.Add(1)
	}
	close(done)
	progressWG.Wait()

	if err := context.Cause(ctx); err != nil {
		return nil, err
	}
	return people, nil
}

func startWorkers(ctx context.Context, cancel context.CancelCauseFunc, repository, revision string, files []string) <-chan fileResult {
	workersCnt := min(runtime.GOMAXPROCS(0), len(files))
	filesCh := make(chan string)
	resultsCh := make(chan fileResult)

	var wg sync.WaitGroup
	wg.Add(workersCnt + 1)
	for range workersCnt {
		go worker(ctx, cancel, repository, revision, filesCh, resultsCh, &wg)
	}

	go func() {
		defer wg.Done()
		defer close(filesCh)
		for _, file := range files {
			select {
			case filesCh <- file:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	return resultsCh
}

func addFileStatistics(people map[string]*stats.PersonStatistics, file string, fileStat stats.FileStatistics, useCommitter bool) {
	for commit, info := range fileStat.Commits {
		name := info.Author
		if useCommitter {
			name = info.Committer
		}

		st, exists := people[name]
		if !exists {
			st = &stats.PersonStatistics{
				Commits: make(map[string]struct{}),
				Files:   make(map[string]struct{}),
			}

			people[name] = st
		}

		st.Lines += info.Lines
		st.Files[file] = struct{}{}
		st.Commits[commit] = struct{}{}
	}
}

func worker(ctx context.Context, cancel context.CancelCauseFunc, repository, revision string,
	files <-chan string,
	results chan<- fileResult,
	wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range files {
		if ctx.Err() != nil {
			return
		}
		stat, err := git.GetFileStatistics(ctx, repository, revision, file)
		if err != nil {
			cancel(err)
			return
		}
		select {
		case results <- fileResult{stat: stat, file: file}:
		case <-ctx.Done():
			return
		}
	}
}
