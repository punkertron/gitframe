package cli

import "github.com/spf13/cobra"

const version = "0.1.0"

func NewRootCmd() *cobra.Command {
	f := defaultOptions()

	rootCmd := &cobra.Command{
		Use:          "gitframe",
		Short:        "Gitframe is a CLI application",
		Long:         "Gitframe is a command-line application built with Go",
		Args:         cobra.NoArgs,
		Version:      version,
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return f.Validate()
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("%+v\n", f)
		},
	}
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	setFlags(rootCmd.Flags(), &f)

	return rootCmd
}
