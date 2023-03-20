package main

import (
	"fmt"
	"regexp"
	"strings"
)

type TexFile struct {
	TxtFile
	Style                     string
	IncludeToc                bool
	BookSectioningCommands    []string
	ArticleSectioningCommands []string
	Buffer                    string
}

func escapeRepl(_ string, match string) string {
	return "\\" + match
}

// escapeString takes an input string and returns a new string with special
// characters escaped using a backslash, except for literal backslashes.
//
// Here's a brief explanation of the special characters' roles in LaTeX:
//
// - #: Represents a parameter in a macro or command definition.
// - %: Indicates the start of a comment. Text following this symbol on the same line is ignored by the LaTeX compiler.
// - &: Used to separate columns in a table or matrix environment.
// - ~: Inserts a non-breaking space, which prevents line breaks between the tilde and the next character.
// - $: Encloses inline mathematical expressions or symbols.
// - _: Used to create subscripts in mathematical expressions.
// - ^: Used to create superscripts in mathematical expressions.
// - {: Marks the beginning of an argument or group in LaTeX commands.
// - }: Marks the end of an argument or group in LaTeX commands.
func escapeString(aStr string) string {
	aStr = strings.ReplaceAll(aStr, "\\\\", ":\\backslash:")
	specialChars := regexp.MustCompile("([#%&~$_^{}])")

	res := specialChars.ReplaceAllStringFunc(aStr, func(match string) string {
		return escapeRepl(aStr, match)
	})

	res = strings.ReplaceAll(res, ":\\backslash:", "$\\backslash$")

	return res
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
	return fmt.Sprintf(`
\documentclass[notitlepage,a4paper]{%s}
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

func min(level int, i int) int {
	if level < i {
		return level
	}
	return i
}

func (tf *TexFile) SectioningCommand(level int) string {
	if tf.Style == "article" {
		return tf.ArticleSectioningCommands[min(level, len(tf.ArticleSectioningCommands))]
	} else if tf.Style == "book" {
		return tf.BookSectioningCommands[min(level, len(tf.BookSectioningCommands))]
	}
	return ""
}
