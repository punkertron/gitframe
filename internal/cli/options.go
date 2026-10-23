package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/punkertron/gitframe/configs"

	"github.com/spf13/pflag"
)

type Options struct {
	Repository         string
	Revision           string
	OrderBy            string
	UseCommitter       bool
	Format             string
	Extensions         []string
	Languages          []string
	languageExtensions []string
	Exclude            []string
	RestrictTo         []string
}

func validateGlob(pattern string) error {
	_, err := filepath.Match(pattern, "")
	return err
}

func (opt *Options) Validate() error {
	if len(opt.Repository) == 0 {
		return fmt.Errorf("empty 'repository'")
	}

	if len(opt.Revision) == 0 {
		return fmt.Errorf("empty 'revision'")
	}

	switch opt.OrderBy {
	case "lines", "commits", "files":
	default:
		return fmt.Errorf("'order-by' must be one of ['lines', 'commits', 'files']")
	}

	switch opt.Format {
	case "tabular", "csv", "json", "json-lines":
	default:
		return fmt.Errorf("'format' must be one of ['tabular', 'csv', 'json', 'json-lines']")
	}

	opt.languageExtensions = nil
	if len(opt.Languages) > 0 {
		for i := 0; i < len(opt.Languages); i++ {
			language := strings.ToLower(opt.Languages[i])
			extensions, ok := configs.ExtensionsByLanguage[language]
			switch {
			case !ok:
				fmt.Fprintf(os.Stderr, "Warning: unsupported language %q\n", opt.Languages[i])
			case len(extensions) == 0:
				fmt.Fprintf(os.Stderr, "Warning: language %q has no extensions; ignoring it\n", opt.Languages[i])
			default:
				opt.languageExtensions = append(opt.languageExtensions, extensions...)
			}
		}
	}

	for _, pattern := range opt.Exclude {
		if err := validateGlob(pattern); err != nil {
			return fmt.Errorf("invalid exclude pattern %q: %w", pattern, err)
		}
	}

	for _, pattern := range opt.RestrictTo {
		if err := validateGlob(pattern); err != nil {
			return fmt.Errorf("invalid restrict-to pattern %q: %w", pattern, err)
		}
	}

	return nil
}

func (opt *Options) FilterFiles(files []string) []string {
	if len(files) == 0 {
		return nil
	}

	var res []string

	for _, file := range files {
		allowed := opt.shouldProcess(file)

		if allowed {
			res = append(res, file)
		}
	}

	return res
}

func (opt *Options) shouldProcess(fileName string) bool {
	// 1. Exclude
	for _, pattern := range opt.Exclude {
		if matched, _ := filepath.Match(pattern, fileName); matched {
			return false
		}
	}

	// 2. RestrictTo
	if len(opt.RestrictTo) > 0 {
		restricted := true
		for _, pattern := range opt.RestrictTo {
			if matched, _ := filepath.Match(pattern, fileName); matched {
				restricted = false
				break
			}
		}
		if restricted == true {
			return false
		}
	}

	// 3. Extensions
	if len(opt.Extensions) > 0 {
		if !slices.Contains(opt.Extensions, filepath.Ext(fileName)) {
			return false
		}
	}

	// 4. Languages.
	if len(opt.languageExtensions) > 0 {
		if !slices.Contains(opt.languageExtensions, filepath.Ext(fileName)) {
			return false
		}
	}

	return true
}

func defaultOptions() Options {
	return Options{
		Repository:   ".",
		Revision:     "HEAD",
		OrderBy:      "lines",
		UseCommitter: false,
		Format:       "tabular",
		Extensions:   []string{},
		Languages:    []string{},
		Exclude:      []string{},
		RestrictTo:   []string{},
	}
}

func setFlags(fs *pflag.FlagSet, opt *Options) {
	fs.StringVar(&opt.Repository, "repository", opt.Repository, "Path to Git repo. Default: current directory")
	fs.StringVar(&opt.Revision, "revision", opt.Revision, "Commit. Default: HEAD")
	fs.StringVar(&opt.OrderBy, "order-by", opt.OrderBy, "Sort results. Default: (lines, commits, files). If the keys are equal, the author with the lexicographically smaller name will appear higher")
	fs.BoolVar(&opt.UseCommitter, "use-committer", opt.UseCommitter, "A boolean flag that replaces the author with the committer in the author's calculations. Default: false")
	fs.StringVar(&opt.Format, "format", opt.Format, "Output format; one of tabular, csv, json, json-lines. Default: tabular")
	fs.StringSliceVar(&opt.Extensions, "extensions", opt.Extensions, "A list of extensions that narrows down the list of files to be included in the calculation; multiple restrictions are separated by commas, for example, '.go,.md'")
	fs.StringSliceVar(&opt.Languages, "languages", opt.Languages, "A list of languages (programming, markup, etc.) that narrows down the set of files included in the calculation; multiple restrictions are separated by commas, for example, 'go,markdown'")
	fs.StringSliceVar(&opt.Exclude, "exclude", opt.Exclude, "A set of Glob patterns that exclude files from the calculation, such as 'foo/*,bar/*'")
	fs.StringSliceVar(&opt.RestrictTo, "restrict-to", opt.RestrictTo, "A set of Glob patterns that include only files matching at least one of these glob patterns, such as 'foo/*,bar/*'")
}
