package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

type File interface {
	Reset()
	Beginning() string
	Ending() string
	AddHeading(level int, heading, fileName string, line int)
	AddCode(code, fileName string, startLine, endLine int)
	AddComment(comment, fileName string, startLine, endLine int)
	EscapeString(aStr string) string
}

type Pbook struct {
	Files     []string
	RegExpSet RegExpSet
	OutFile   File
}

func NewPbook(files []string, outFile File, regExpSet RegExpSet) *Pbook {
	return &Pbook{
		Files:     files,
		RegExpSet: regExpSet,
		OutFile:   outFile,
	}
}

func (p *Pbook) FormatBuffer() {
	p.OutFile.Reset()

	for _, file := range p.Files {
		data, err := ioutil.ReadFile(file)
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

		for len(lines) > 0 {
			line := lines[0]

			if headingRe.MatchString(line) {
				p.doHeading(line, file, lineCounter)
				lines = lines[1:]
				lineCounter++
			} else if commentRe.MatchString(line) {
				lineCounter += p.doComment(lines, file, lineCounter)
				lines = lines[lineCounter:]
			} else {
				lineCounter += p.doCode(lines, file, lineCounter)
				lines = lines[lineCounter:]
			}
		}
	}
}

func (p *Pbook) doHeading(line, fileName string, lineCounter int) {
	headingRe := regexp.MustCompile(p.RegExpSet.Heading)
	matches := headingRe.FindStringSubmatch(line)
	level := len(matches[1]) - 1
	heading := line[len(matches[0]):]
	p.OutFile.AddHeading(level, heading, fileName, lineCounter)
}

func (p *Pbook) doComment(lines []string, fileName string, lineCounter int) int {
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

func (p *Pbook) doCode(lines []string, fileName string, lineCounter int) int {
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
