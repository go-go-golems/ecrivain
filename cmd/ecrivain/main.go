package main

import (
	"context"
	_ "embed"
	"fmt"
	clay "github.com/go-go-golems/clay/pkg"
	"github.com/go-go-golems/ecrivain/pkg"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/help"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
)

type EcrivainSettings struct {
	OutputFile string `glazed.parameter:"output-file"`
	OutputType string `glazed.parameter:"output-type"`
	Title      string `glazed.parameter:"title"`
	Author     string `glazed.parameter:"author"`
	Style      string `glazed.parameter:"style"`
	IncludeToc bool   `glazed.parameter:"include-toc"`
}

// NOTE(manuel, 2023-03-20) More flags to add:
// - globs
// - watcher

func NewEcrivainLayer() (layers.ParameterLayer, error) {
	ecrivainLayer, err := layers.NewParameterLayer(
		"ecrivain",
		"Ecrivain rendering options",
		layers.WithFlags(
			parameters.NewParameterDefinition(
				"output-file",
				parameters.ParameterTypeString,
				parameters.WithShortFlag("o"),
				parameters.WithHelp("Output file"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"output-type",
				parameters.ParameterTypeChoice,
				parameters.WithShortFlag("t"),
				parameters.WithHelp("Output type"),
				parameters.WithChoices([]string{"tex", "txt"}),
			),
			parameters.NewParameterDefinition(
				"title",
				parameters.ParameterTypeString,
				parameters.WithShortFlag("T"),
				parameters.WithHelp("Title"),
			),
			parameters.NewParameterDefinition(
				"author",
				parameters.ParameterTypeString,
				parameters.WithShortFlag("a"),
				parameters.WithHelp("Author"),
			),
			parameters.NewParameterDefinition(
				"style",
				parameters.ParameterTypeString,
				parameters.WithShortFlag("s"),
				parameters.WithHelp("Style"),
				parameters.WithDefault("article"),
			),
			parameters.NewParameterDefinition(
				"include-toc",
				parameters.ParameterTypeBool,
				parameters.WithShortFlag("i"),
				parameters.WithHelp("Include table of contents"),
				parameters.WithDefault(true),
			),
		))
	if err != nil {
		return nil, err
	}

	return ecrivainLayer, nil
}

type RenderCommand struct {
	*cmds.CommandDescription
}

func (r *RenderCommand) RunIntoWriter(
	ctx context.Context,
	parsedLayers map[string]*layers.ParsedParameterLayer,
	ps map[string]interface{},
	w io.Writer) error {
	ecrivainLayer, ok := parsedLayers["ecrivain"]
	if !ok {
		return fmt.Errorf("ecrivain layer not found")
	}
	settings := &EcrivainSettings{}
	err := parameters.InitializeStructFromParameters(settings, ecrivainLayer.Parameters)
	if err != nil {
		return err
	}

	if settings.OutputType == "" {
		// get settings OutputFile's extension
		ext := filepath.Ext(settings.OutputFile)

		switch ext {
		case ".tex":
			settings.OutputType = "tex"
		case ".txt":
			settings.OutputType = "txt"
		default:
			return fmt.Errorf("unknown output type: %s", ext)
		}
	}

	var outFile pkg.File

	switch settings.OutputType {
	case "tex":
		outFile = pkg.NewTexFile(
			settings.Title,
			settings.Author,
			settings.Style,
			settings.IncludeToc,
		)

	case "txt":
		outFile = pkg.NewTxtFile(settings.Title, settings.Author)

	default:
		return fmt.Errorf("ecrivain type: %s", settings.OutputType)
	}

	config, err := pkg.LoadLanguages()
	if err != nil {
		return err
	}

	language := ""

	files := ps["source-files"].([]string)
	for _, file := range files {
		lang, err := config.DetectLanguage(file)
		if err != nil {
			return err
		}

		if language == "" {
			language = lang
		} else if language != lang {
			return fmt.Errorf("cannot mix languages: %s and %s", language, lang)
		}
	}

	book, err := config.CreateBook(files, language, outFile)
	if err != nil {
		return err
	}

	// NOTE(manuel, 2023-03-20) iterating through files should probably be done here, not FormatBuffer()
	book.FormatBuffer()

	_, err = outFile.Write(w)
	if err != nil {
		return err
	}

	return nil
}

func NewRenderCommand() (*RenderCommand, error) {
	ecrivainLayer, err := NewEcrivainLayer()
	if err != nil {
		return nil, err
	}

	description := cmds.NewCommandDescription("render",
		cmds.WithShort("Render a set of source files into a single document"),
		cmds.WithLayers(ecrivainLayer),
		cmds.WithArguments(
			parameters.NewParameterDefinition(
				"source-files",
				parameters.ParameterTypeStringList,
				parameters.WithHelp("Source files"),
				parameters.WithRequired(true),
			),
		),
	)

	return &RenderCommand{
		CommandDescription: description,
	}, nil
}

func main() {
	var err error

	helpSystem := help.NewHelpSystem()
	//err := helpSystem.LoadSectionsFromFS(docFS, ".")
	//cobra.CheckErr(err)

	rootCmd := &cobra.Command{
		Use:   "ecrivain [flags]",
		Short: "A tool to generate beautiful books out of beautiful source code",
		Args:  cobra.MinimumNArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Reinitalize the logger after the config has been loaded
			err := clay.InitLogger()
			cobra.CheckErr(err)
		},
	}

	helpSystem.SetupCobraRootCommand(rootCmd)

	err = clay.InitViper("cliopatra", rootCmd)
	cobra.CheckErr(err)
	err = clay.InitLogger()
	cobra.CheckErr(err)

	renderCommand, err := NewRenderCommand()
	cobra.CheckErr(err)
	cobraRenderCommand, err := cli.BuildCobraCommandFromWriterCommand(renderCommand)
	cobra.CheckErr(err)

	rootCmd.AddCommand(cobraRenderCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
