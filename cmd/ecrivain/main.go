package main

import (
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
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

type Config struct {
	Languages map[string]RegExpSet `yaml:"languages"`
}

func (cfg Config) CreateBook(ext string, outFile File) *Pbook {
	for lang, regExpSet := range cfg.Languages {
		for _, fileExt := range regExpSet.FileExts {
			if strings.ToLower(ext) == strings.ToLower(fileExt) {
				return NewPbook([]string{}, outFile, regExpSet)
			}
		}
	}
	return nil
}

func escapeString(aStr string) string {
	aStr = strings.Replace(aStr, "\\", "$\\backslash$", -1)
	escapeRe := regexp.MustCompile("([#%&~$_^{}])")
	return escapeRe.ReplaceAllStringFunc(aStr, func(match string) string {
		return "\\" + match
	})
}

func main() {
	var (
		texClass   string
		pbookType  string
		title      string
		author     string
		output     string
		style      string
		includeToc bool
	)

	rootCmd := &cobra.Command{
		Use:   "pbook [flags] file ...",
		Short: "A tool to generate pbooks",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := Config{}
			err := yaml.Unmarshal(regexpsYaml, &cfg)
			if err != nil {
				fmt.Printf("Error parsing YAML: %v", err)
				return
			}

			for _, file := range args {
				// get file extension
				name, ext := filepath.Split(file)

				var outFile File

				p := cfg.CreateBook(ext, NewTexFile(output, title, author, style, includeToc))
				if p == nil {
					fmt.Printf("Error creating pbook for file: %s\n", file)
				}

				p.FormatBuffer()

				if p.OutFile.Title != "" {
					p.OutFile.Buffer.WriteString("\\end{document}\n")
				}

				if p.OutFile.Author != "" {
					p.OutFile.Buffer.WriteString(fmt.Sprintf("\\author{%s}\n", p.OutFile.Author))
				}

				if p.OutFile.Title != "" {
					p.OutFile.Buffer.WriteString(fmt.Sprintf("\\title{%s}\n", p.OutFile.Title))
					p.OutFile.Buffer.WriteString("\\maketitle\n")
				}

			}
			return nil
		},
	}

	rootCmd.Flags().StringVarP(&texClass, "class", "c", "TexFile", "Output class (TexFile|IdqTexFile|TxtFile)")
	rootCmd.Flags().StringVarP(&pbookType, "type", "T", "", "Pbook type (C|Lisp)")
	rootCmd.Flags().StringVarP(&title, "title", "t", "", "Title")
	rootCmd.Flags().StringVarP(&author, "author", "a", "", "Author")
	rootCmd.Flags().BoolVarP(&includeToc, "toc", "O", true, "Include table of contents")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output file")
	rootCmd.Flags().StringVarP(&style, "style", "s", "article", "Style")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
