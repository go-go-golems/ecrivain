package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEscapeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello\\\\World! #Test%", "Hello$\\backslash$World! \\#Test\\%"},
		{"Normal text", "Normal text"},
		{"$LaTeX$ & _example_{", "\\$LaTeX\\$ \\& \\_example\\_\\{"},
		{"{Hello} ^World#", "\\{Hello\\} \\^World\\#"},
		{"%Comment & $Math$", "\\%Comment \\& \\$Math\\$"},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "#%&~$_^{}\\\\",
			expected: "\\#\\%\\&\\~\\$\\_\\^\\{\\}$\\backslash$",
		},
		{
			input:    "#%#%#%",
			expected: "\\#\\%\\#\\%\\#\\%",
		},
		{
			input:    "A#B%C&D~E$F_G^H{I}",
			expected: "A\\#B\\%C\\&D\\~E\\$F\\_G\\^H\\{I\\}",
		},
		{
			input:    "A\\\\#B\\\\%C\\\\",
			expected: "A$\\backslash$\\#B$\\backslash$\\%C$\\backslash$",
		},
	}

	for _, test := range tests {
		output := escapeString(test.input)
		if output != test.expected {
			assert.Equal(t, test.expected, output)
		}
	}
}
