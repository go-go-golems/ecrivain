package main

import "fmt"

type TexFile struct {
	TxtFile
	Style                     string
	IncludeToc                bool
	BookSectioningCommands    []string
	ArticleSectioningCommands []string
	Buffer                    string
}

func NewTexFile(title string, author string, style string, includeToc bool) *TexFile {
	return &TexFile{
		TxtFile:                   TxtFile{Title: title, Author: author},
		Style:                     style,
		IncludeToc:                includeToc,
		BookSectioningCommands:    []string{"chapter", "section", "subsection", "subsubsection"},
		ArticleSectioningCommands: []string{"section", "subsection", "subsubsection"},
		Buffer:                    "",
	}
}

func (tf *TexFile) Beginning() string {
	authorString := ""
	if tf.Author != "" {
		authorString = fmt.Sprintf("\\author{%s}\n", tf.Author)
	}
	tocString := ""
	if tf.IncludeToc {
		tocString = "\\tableofcontents\n"
	}
	return fmt.Sprintf(`\documentclass[notitlepage,a4paper]{%s}
\usepackage{fancyvrb,color,palatino}
\definecolor{gray}{gray}{0.6}
\title{%s}
%s
\begin{document}
\maketitle
%s`, tf.Style, escapeString(tf.Title), authorString, tocString)
}

func (tf *TexFile) Ending() string {
	return "\\end{document}\n"
}

func (tf *TexFile) SectioningCommand(level int) string {
	if tf.Style == "article" {
		return tf.ArticleSectioningCommands[min(level, len(tf.ArticleSectioningCommands))]
	} else if tf.Style == "book" {
		return tf.BookSectioningCommands[min(level, len(tf.BookSectioningCommands))]
	}
	return ""
}

func min(level int, i int) int {
	if level < i {
		return level
	}
	return i
}
