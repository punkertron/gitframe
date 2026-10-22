package cli

import (
	"github.com/punkertron/gitframe/internal/git"
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
	files, err := git.GetFiles(opt.Repository, opt.Revision)
	if err != nil {
		return err
	}

	filtered := opt.FilterFiles(files)
	//cmd.Println(filtered)

	people := make(map[string]*stats.PersonStatistics)

	for _, file := range filtered {
		stat, err := git.GetFileStatistics(opt.Repository, opt.Revision, file)
		if err != nil {
			return err
		}

		//cmd.Println(stats.Path)

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

			//cmd.Println(commit, info)
		}
	}

	// for name, person := range people {
	// 	cmd.Println(name, len(person.Commits), len(person.Files), person.Lines)
	// }

	res, err := sort.SortData(people, opt.OrderBy)
	if err != nil {
		return err
	}
	cmd.Println(res)

	return nil
}
