package pkg

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type File interface {
	Reset()
	Beginning()
	Ending()
	AddHeading(level int, heading, fileName string, line int)
	AddCode(code, fileName string, startLine, endLine int)
	AddComment(comment, fileName string, startLine, endLine int)
	Write(w io.Writer) (int, error)
}

type Book struct {
	Files     []string
	RegExpSet RegExpSet
	OutFile   File
}

func NewBook(files []string, outFile File, regExpSet RegExpSet) *Book {
	return &Book{
		Files:     files,
		RegExpSet: regExpSet,
		OutFile:   outFile,
	}
}

func (p *Book) FormatBuffer() {
	p.OutFile.Reset()

	for _, file := range p.Files {
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file: %s\n", file)
			continue
		}

		content := string(data)
		if p.RegExpSet.Remove != "" {
			removeRe := regexp.MustCompile(p.RegExpSet.Remove)
			content = removeRe.ReplaceAllString(content, "")
		}

		commentRe := regexp.MustCompile(p.RegExpSet.Comment)
		headingRe := regexp.MustCompile(p.RegExpSet.Heading)

		lines := strings.Split(content, "\n")
		lineCounter := 0

		linesProcessed := p.doBeginning(file, lineCounter)
		lineCounter += linesProcessed
		lines = lines[linesProcessed:]

		for len(lines) > 0 {
			line := lines[0]

			// check if line is just whitespace and skip
			if len(strings.TrimSpace(line)) == 0 {
				lines = lines[1:]
				lineCounter++
				continue
			}

			if headingRe.MatchString(line) {
				p.doHeading(line, file, lineCounter)
				lines = lines[1:]
				lineCounter++
			} else if commentRe.MatchString(line) {
				linesProcessed := p.doComment(lines, file, lineCounter)
				lineCounter += linesProcessed
				lines = lines[linesProcessed:]
			} else {
				linesProcessed := p.doCode(lines, file, lineCounter)
				lineCounter += linesProcessed
				lines = lines[linesProcessed:]
			}
		}

		p.doEnding(file, lineCounter)
	}
}

func (p *Book) doHeading(line, fileName string, lineCounter int) {
	headingRe := regexp.MustCompile(p.RegExpSet.Heading)
	matches := headingRe.FindStringSubmatch(line)
	level := len(matches[1]) - 1
	heading := line[len(matches[0]):]
	p.OutFile.AddHeading(level, heading, fileName, lineCounter)
}

func (p *Book) doComment(lines []string, fileName string, lineCounter int) int {
	commentRe := regexp.MustCompile(p.RegExpSet.Comment)
	comment := ""
	commentLines := 0

	for len(lines) > 0 {
		line := lines[0]
		if !commentRe.MatchString(line) {
			break
		}
		comment += line + "\n"
		lines = lines[1:]
		commentLines++
		lineCounter++
	}

	p.OutFile.AddComment(comment, fileName, lineCounter, lineCounter+commentLines)
	return commentLines
}

func (p *Book) doCode(lines []string, fileName string, lineCounter int) int {
	code := ""
	codeLines := 0

	for len(lines) > 0 {
		line := lines[0]
		if p.RegExpSet.Comment != "" && regexp.MustCompile(p.RegExpSet.Comment).MatchString(line) {
			break
		}
		if p.RegExpSet.Heading != "" && regexp.MustCompile(p.RegExpSet.Heading).MatchString(line) {
			break
		}
		code += line + "\n"
		lines = lines[1:]
		codeLines++
		lineCounter++
	}

	p.OutFile.AddCode(code, fileName, lineCounter, lineCounter+codeLines)
	return codeLines
}

func (p *Book) doBeginning(file string, counter int) int {
	p.OutFile.Beginning()
	return 0
}

func (p *Book) doEnding(file string, counter int) {
	p.OutFile.Ending()
}
