package pkg

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type TxtFile struct {
	Title  string
	Author string
	Buffer strings.Builder
}

func NewTxtFile(title, author string) *TxtFile {
	return &TxtFile{
		Title:  title,
		Author: author,
	}
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

func (t *TxtFile) Beginning() string {
	return ""
}

func (t *TxtFile) Ending() string {
	return ""
}

func (t *TxtFile) Write(w io.Writer) (int, error) {
	return w.Write([]byte(t.Buffer.String()))
}
