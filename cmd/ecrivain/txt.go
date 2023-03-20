package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type TxtFile struct {
	Title  string
	Author string
	Buffer strings.Builder
}

func (t *TxtFile) Reset() {
	t.Buffer.Reset()
}

func (t *TxtFile) Add(s string) {
	t.Buffer.WriteString(s)
}

func (t *TxtFile) AddHeading(level int, heading, fileName string, line int) {
	heading = strings.Join(strings.Fields(heading), " ")
	t.Buffer.WriteString(fmt.Sprintf("++ %s\n", heading))
}

func (t *TxtFile) AddComment(comment, fileName string, startLine, endLine int) {
	t.Buffer.WriteString(comment)
}

func (t *TxtFile) AddCode(code, fileName string, startLine, endLine int) {
	code = strings.TrimRight(code, " \t\n")
	code = strings.ReplaceAll(code, "\t", "   ")
	t.Buffer.WriteString(fmt.Sprintf("\n== %s (%s:%d) ================\n%s\n", filepath.Base(fileName), fileName, startLine, code))
	t.Buffer.WriteString("=========================================================\n\n")
}
