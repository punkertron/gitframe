package configs

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

type Language struct {
	Name       string   `json:"name"`
	Extensions []string `json:"extensions"`
}

//go:embed language_extensions.json
var data []byte

var ExtensionsByLanguage = mustLoadLanguages()

func mustLoadLanguages() map[string][]string {
	var languages []Language

	if err := json.Unmarshal(data, &languages); err != nil {
		panic(fmt.Errorf("failed to load languages: %w", err))
	}

	m := make(map[string][]string, len(languages))
	for _, language := range languages {
		m[strings.ToLower(language.Name)] = language.Extensions
	}

	return m
}
