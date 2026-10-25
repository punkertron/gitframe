package cli

import (
	"github.com/spf13/cobra"

	"github.com/punkertron/gitframe/internal/git"
	"github.com/punkertron/gitframe/internal/output"
	"github.com/punkertron/gitframe/internal/sort"
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
			if err := opt.Validate(); err != nil {
				return err
			}
			return opt.prepareLanguages(cmd.ErrOrStderr())
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
	revision, err := git.ResolveRevision(cmd.Context(), opt.Repository, opt.Revision)
	if err != nil {
		return err
	}

	files, err := git.GetFiles(cmd.Context(), opt.Repository, revision)
	if err != nil {
		return err
	}

	filtered := opt.FilterFiles(files)

	people, err := collectStatistics(cmd.Context(), cmd.ErrOrStderr(), opt, revision, filtered)
	if err != nil {
		return err
	}

	res, err := sort.SortData(people, opt.OrderBy)
	if err != nil {
		return err
	}

	return output.PrintResult(res, opt.Format, cmd.OutOrStdout())
}
