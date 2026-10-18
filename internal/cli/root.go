package cli

import (
	"github.com/punkertron/gitframe/internal/git"

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
			// cmd.Printf("%+v\n", opt)

			files, err := git.GetFiles(opt.Repository, opt.Revision)
			if err != nil {
				return err
			}

			f := opt.FilterFiles(files)
			cmd.Println(f)

			return nil
		},
	}
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	setFlags(rootCmd.Flags(), &opt)

	return rootCmd
}
