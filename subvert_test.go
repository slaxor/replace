package replace

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Test with single word",
			input:    "Camel",
			expected: "Camel",
		},
		{
			name:     "Test with multiple words",
			input:    "CamelCase",
			expected: "Camel_Case",
		},
		{
			name:     "Test with spaces",
			input:    "this is a Test",
			expected: "this_is_a_Test",
		},
		{
			name:     "Test aleady normalized",
			input:    normalize("this_is_a_Test"),
			expected: "this_is_a_Test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, normalize(tt.input))
		})
	}
}

// Lots of edge cases aren't working yet
func TestSubvert(t *testing.T) {
	tests := []struct {
		name string
		i    string
		exp  string
		o    string
		n    string
	}{
		{
			name: "lower case (happy path)",
			i:    "This is a string with old_string",
			exp:  "This is a string with new_string",
			o:    "old_string",
			n:    "new_string",
		},
		{
			name: "camel case",
			i:    "This is a string with OldString",
			exp:  "This is a string with NewString",
			o:    "old_string",
			n:    "new_string",
		},
		{
			name: "all caps",
			i:    "This is a string with OLDSTRING",
			exp:  "This is a string with NEWSTRING",
			o:    "old_string",
			n:    "new_string",
		},
		{
			name: "all lower",
			i:    "This is a string with oldstring",
			exp:  "This is a string with newstring",
			o:    "old_string",
			n:    "new_string",
		},
		{
			name: "camel to snake",
			i:    "This is a string with old_string",
			exp:  "This is a string with new_string",
			o:    "OldString",
			n:    "NewString",
		},
		{
			name: "multi word",
			i:    "This is a string with old string",
			exp:  "This is a string with new string",
			o:    "old_string",
			n:    "new_string",
		},
		{
			name: "initial lower case", // go private should be kept that way
			i:    "This is a string with oldString",
			exp:  "This is a string with newString",
			o:    "old_string",
			n:    "new_string",
		},
		{
			name: "all but initial caps", // its a bit silly, but you never know what people do
			i:    "This is a string with oLDSTRING",
			exp:  "This is a string with nEWSTRING",
			o:    "old_string",
			n:    "new_string",
		},
		{
			name: "CamelCase replacement",
			i:    "This is a string with old string",
			exp:  "This is a string with new string",
			o:    "OldString",
			n:    "NewString",
		},
		{
			name: "multi word replacement",
			i:    "This is a test string with old_String.",
			exp:  "This is a test string with new_String.",
			o:    "old string",
			n:    "new string",
		},
		{
			name: "more than one replacement",
			i:    "old string. old_string. OldString. Old String. OLD STRING. old-string. oldString.",
			exp:  "new string. new_string. NewString. New String. NEW STRING. old-string. newString.",
			o:    "old_string",
			n:    "new_string",
		},
		{
			name: "uneven replacement",
			i:    "This is a test string with old String.",
			exp:  "This is a test string with new String With An Extra.",
			o:    "old_string",
			n:    "newStringWithAnExtra",
		},
		{
			name: "multi line",
			i:    "This is a test string with old_string, Old_String and OLD\n\n\tSTRING.",
			exp:  "This is a test string with new_string, New_String and NEW\n\n\tSTRING.",
			o:    "old_string",
			n:    "new_string",
		},
		{
			name: "regex match",
			i:    "open shut, outward_southward, overcomeSucceed, OvalShape, output signal",
			exp:  "new string, new_string, newString, NewString, new string",
			o:    "o[a-z]+ s[]a-z]+",
			n:    "new_string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Subvert(tt.i, tt.o, tt.n)
			assert.Equal(t, tt.exp, result)
		})
	}
}

func TestSubvertFileContent(t *testing.T) {
	dir := t.TempDir()
	tests := []struct {
		name     string
		file     string
		perm     os.FileMode
		content  []byte
		o        string
		n        string
		expected []byte
	}{
		{
			name:     "happy path",
			file:     "test_happy_path.txt",
			perm:     0640,
			content:  []byte("This is a test string with old_string."),
			o:        "old_string",
			n:        "new_string",
			expected: []byte("This is a test string with new_string."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := filepath.Join(dir, tt.file)
			err := ioutil.WriteFile(f, tt.content, tt.perm)
			assert.NoError(t, err)
			err = SubvertFileContent(f, tt.o, tt.n)
			assert.NoError(t, err)
			result, err := os.ReadFile(f)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result, "expected\n%s\nis not equal to\n%s\n", tt.expected, result)
			// assert.FileModePerm(tt.perm, f)
		})
	}
}
