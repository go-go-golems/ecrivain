package pkg

import (
	_ "embed"
	"gopkg.in/yaml.v3"
	"strings"
)

type RegExpSet struct {
	Comment  string   `yaml:"commentRe"`
	Heading  string   `yaml:"headingRe"`
	Remove   string   `yaml:"removeRe,omitempty"`
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
	err := yaml.Unmarshal(regexpsYaml, &config.Languages)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (cfg *LanguageConfig) CreateBook(files []string, language string, outFile File) (*Book, error) {
	regExpSet, ok := cfg.Languages[language]
	if !ok {
		return nil, errors.Errorf("language not found: %s", language)
	}
	return NewBook(files, outFile, regExpSet), nil
}

func (cfg *LanguageConfig) DetectLanguage(file string) (string, error) {
	for lang, regExpSet := range cfg.Languages {
		for _, fileExt := range regExpSet.FileExts {
			if strings.HasSuffix(file, fileExt) {
				return lang, nil
			}
		}
	}
	return "", errors.Errorf("no language found for file: %s", file)
}
