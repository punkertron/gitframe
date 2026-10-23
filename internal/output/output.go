package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"

	"github.com/punkertron/gitframe/internal/stats"
)

func PrintResult(results []stats.Result,
	format string,
	writer io.Writer) error {

	switch format {
	case "tabular":
		return tabularOutput(results, writer)
	case "csv":
		return csvOutput(results, writer)
	case "json":
		return jsonOutput(results, writer)
	case "json-lines":
		return jsonLinesOutput(results, writer)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func tabularOutput(results []stats.Result, writer io.Writer) error {
	w := new(tabwriter.Writer)

	w.Init(writer, 0, 0, 1, ' ', 0)
	if _, err := fmt.Fprintf(w, "Name\tLines\tCommits\tFiles\n"); err != nil {
		return err
	}
	for _, result := range results {
		if _, err := fmt.Fprintf(
			w, "%s\t%d\t%d\t%d\n", result.Name, result.Lines,
			result.Commits, result.Files); err != nil {
			return err
		}
	}
	return w.Flush()
}

func csvOutput(results []stats.Result, writer io.Writer) error {
	w := csv.NewWriter(writer)
	if err := w.Write([]string{"Name", "Lines", "Commits", "Files"}); err != nil {
		return err
	}

	for _, result := range results {
		record := []string{result.Name,
			strconv.Itoa(result.Lines),
			strconv.Itoa(result.Commits),
			strconv.Itoa(result.Files),
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func jsonOutput(results []stats.Result, writer io.Writer) error {
	return json.NewEncoder(writer).Encode(results)
}

func jsonLinesOutput(results []stats.Result, writer io.Writer) error {
	enc := json.NewEncoder(writer)

	for _, result := range results {
		if err := enc.Encode(result); err != nil {
			return err
		}
	}
	return nil
}
