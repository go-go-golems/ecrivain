package pkg

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
	}
}

func (t *TexFile) Beginning() {
	authorString := ""
	if t.Author != "" {
		authorString = fmt.Sprintf("\\author{%s}\n", t.Author)
	}
	tocString := ""
	if t.IncludeToc {
		tocString = "\\tableofcontents\n"
	}
	t.Buffer.WriteString(fmt.Sprintf(`
\documentclass[notitlepage,a4paper]{%s}
\usepackage{fancyvrb,color,palatino}
\definecolor{gray}{gray}{0.6}
\title{%s}
%s
\begin{document}
\maketitle
%s`, t.Style, escapeString(t.Title), authorString, tocString))
}

func (t *TexFile) Ending() {
	t.Buffer.WriteString("\\end{document}\n")
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

func (t *TexFile) AddHeading(level int, heading, fileName string, line int) {
	whitespace := regexp.MustCompile(`\s+`)
	heading = whitespace.ReplaceAllString(heading, " ")
	sectionCmd := t.SectioningCommand(level)
	escapedHeading := escapeString(heading)
	t.Buffer.WriteString(fmt.Sprintf("\\%s{%s}\n", sectionCmd, escapedHeading))
}

func (t *TexFile) AddComment(comment, fileName string, startLine, endLine int) {
	escapedComment := escapeString(comment)
	codeComment := regexp.MustCompile("`([^`']*)'")
	comment = codeComment.ReplaceAllString(escapedComment, "{\\\\tt \\1}")
	quoteComment := regexp.MustCompile("\"([^\"]*)\"")
	t.Buffer.WriteString(quoteComment.ReplaceAllString(comment, "``\\1''"))
}

func (t *TexFile) AddCode(code, fileName string, startLine, endLine int) {
	code = strings.TrimRight(code, "\n")
	code = strings.ReplaceAll(code, "\\end{Verbatim}", "\\\\_end{Verbatim}")
	code = strings.ReplaceAll(code, "\t", "   ")
	t.Buffer.WriteString("\n\\begin{Verbatim}[fontsize=\\small,frame=leftline,framerule=0.9mm," +
		"rulecolor=\\color{gray},framesep=5.1mm,xleftmargin=5mm,fontfamily=cmtt]\n")
	t.Buffer.WriteString(code)
	t.Buffer.WriteString("\n\\end{Verbatim}\n")
}
