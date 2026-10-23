package cli

import (
	"github.com/punkertron/gitframe/internal/git"
	"github.com/punkertron/gitframe/internal/output"
	"github.com/punkertron/gitframe/internal/sort"
	"github.com/punkertron/gitframe/internal/stats"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

func NewRootCmd() *cobra.Command {
	opt := defaultOptions()

	rootCmd := &cobra.Command{
		Use:           "gitframe",
		Short:         "Gitframe is a CLI application",
		Long:          "Gitframe is a command-line application built with Go",
		Args:          cobra.NoArgs,
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opt.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, &opt)
		},
	}
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	setFlags(rootCmd.Flags(), &opt)

	return rootCmd
}

func run(cmd *cobra.Command, opt *Options) error {
	revision, err := git.ResolveRevision(opt.Repository, opt.Revision)
	if err != nil {
		return err
	}

	files, err := git.GetFiles(opt.Repository, revision)
	if err != nil {
		return err
	}

	filtered := opt.FilterFiles(files)

	people := make(map[string]*stats.PersonStatistics)

	for _, file := range filtered {
		stat, err := git.GetFileStatistics(opt.Repository, revision, file)
		if err != nil {
			return err
		}

		for commit, info := range stat.Commits {
			name := info.Author
			if opt.UseCommitter {
				name = info.Committer
			}

			st, exists := people[name]
			if !exists {
				st = &stats.PersonStatistics{
					Lines:   0,
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

	res, err := sort.SortData(people, opt.OrderBy)
	if err != nil {
		return err
	}

	return output.PrintResult(res, opt.Format, cmd.OutOrStdout())
}
