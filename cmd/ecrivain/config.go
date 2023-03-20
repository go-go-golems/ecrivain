package main

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

type RegExpSet struct {
	Comment  string   `yaml:"comment"`
	Heading  string   `yaml:"heading"`
	Remove   string   `yaml:"remove,omitempty"`
	FileExts []string `yaml:"fileExts,omitempty"`
}

//go:embed regexps.yaml
var regexpsYaml []byte

type LanguageConfig struct {
	Languages map[string]RegExpSet `yaml:"languages"`
}

func LoadLanguages() (*LanguageConfig, error) {
	config := &LanguageConfig{
		Languages: map[string]RegExpSet{},
	}

	// load config from regexpsYaml
	err := yaml.Unmarshal(regexpsYaml, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (cfg *LanguageConfig) CreateBook(ext string, outFile File) (*Book, error) {
	for _, regExpSet := range cfg.Languages {
		for _, fileExt := range regExpSet.FileExts {
			if strings.EqualFold(ext, fileExt) {
				return NewBook([]string{}, outFile, regExpSet), nil
			}
		}
	}
	return nil, fmt.Errorf("no language found for extension: %s", ext)
}

func (cfg *LanguageConfig) DetectLanguage(file string) (string, error) {
	for lang, regExpSet := range cfg.Languages {
		for _, fileExt := range regExpSet.FileExts {
			if strings.HasSuffix(file, fileExt) {
				return lang, nil
			}
		}
	}
	return "", fmt.Errorf("no language found for file: %s", file)
}
