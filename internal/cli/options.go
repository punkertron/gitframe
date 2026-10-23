package cli

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/pflag"

	"github.com/punkertron/gitframe/configs"
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
		return errors.New("empty 'repository'")
	}

	if len(opt.Revision) == 0 {
		return errors.New("empty 'revision'")
	}

	switch opt.OrderBy {
	case "lines", "commits", "files":
	default:
		return errors.New("'order-by' must be one of ['lines', 'commits', 'files']")
	}

	switch opt.Format {
	case "tabular", "csv", "json", "json-lines":
	default:
		return errors.New("'format' must be one of ['tabular', 'csv', 'json', 'json-lines']")
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

func (opt *Options) prepareLanguages(warnings io.Writer) error {
	opt.languageExtensions = nil
	for _, name := range opt.Languages {
		extensions, ok := configs.ExtensionsByLanguage[strings.ToLower(name)]
		switch {
		case !ok:
			if _, err := fmt.Fprintf(warnings, "Warning: unsupported language %q\n", name); err != nil {
				return fmt.Errorf("write language warning: %w", err)
			}
		case len(extensions) == 0:
			if _, err := fmt.Fprintf(warnings, "Warning: language %q has no extensions; ignoring it\n", name); err != nil {
				return fmt.Errorf("write language warning: %w", err)
			}
		default:
			opt.languageExtensions = append(opt.languageExtensions, extensions...)
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
		if opt.shouldProcess(file) {
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
		if restricted {
			return false
		}
	}

	extension := filepath.Ext(fileName)
	if len(opt.Extensions) > 0 && !slices.Contains(opt.Extensions, extension) {
		return false
	}

	if len(opt.languageExtensions) > 0 && !slices.Contains(opt.languageExtensions, extension) {
		return false
	}

	return true
}

func defaultOptions() Options {
	return Options{
		Repository: ".",
		Revision:   "HEAD",
		OrderBy:    "lines",
		Format:     "tabular",
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
