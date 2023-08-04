package replace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
func TestSubvert(t *testing.T) {
	tests := []struct {
		name     string
		oldName  string
		newName  string
		input    string
		expected string
	}{
		{
			name:     "Test lower case",
			oldName:  "old_foo",
			newName:  "new_foo",
			input:    "This is a string with old_foo",
			expected: "This is a string with new_foo",
		},
		{
			name:     "Test capitalized",
			oldName:  "OldFoo",
			newName:  "NewFoo",
			input:    "This is a string with OldFoo",
			expected: "This is a string with NewFoo",
		},
		{
			name:     "Test all caps",
			oldName:  "OLDFOO",
			newName:  "NEWFOO",
			input:    "This is a string with OLDFOO",
			expected: "This is a string with NEWFOO",
		},
		{
			name:     "Test all lower",
			oldName:  "oldfoo",
			newName:  "newfoo",
			input:    "This is a string with oldfoo",
			expected: "This is a string with newfoo",
		},
		{
			name:     "Test unmatched formats",
			oldName:  "OldFoo",
			newName:  "NewFoo",
			input:    "This is a string with OLDFOO",
			expected: "This is a string with NEWFOO",
		},
		{
			name:     "Test camel case to snake case",
			oldName:  "OldFoo",
			newName:  "NewFoo",
			input:    "This is a string with old_foo",
			expected: "This is a string with new_foo",
		},
		{
			name:     "Test snake case to camel case",
			oldName:  "old_foo",
			newName:  "new_foo",
			input:    "This is a string with OldFoo",
			expected: "This is a string with NewFoo",
		},
		{
			name:     "Test multi word occurrences",
			oldName:  "old_foo",
			newName:  "new_foo",
			input:    "This is a string with old Foo",
			expected: "This is a string with New Foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replaceFn := Subvert(tt.oldName, tt.newName)
			assert.Equal(t, tt.expected, replaceFn(tt.input))
		})
	}
}
*/
/*
	func TestToSnakeCase(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "Test with single word",
				input:    "Camel",
				expected: "camel",
			},
			{
				name:     "Test with multiple words",
				input:    "CamelCase",
				expected: "camel_case",
			},
			{
				name:     "Test with longer string",
				input:    "ThisIsATest",
				expected: "this_is_a_test",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, ToSnakeCase(tt.input))
			})
		}
	}

	func TestToCamelCase(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "Test with single word",
				input:    "camel",
				expected: "Camel",
			},
			{
				name:     "Test with multiple words",
				input:    "camel_case",
				expected: "CamelCase",
			},
			{
				name:     "Test with longer string",
				input:    "this_is_a_test",
				expected: "ThisIsATest",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, ToCamelCase(tt.input))
			})
		}
	}
*/

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

func TestSubvert(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		old      string
		new      string
		expected string
	}{
		{
			name:     "All lower case replacement",
			input:    "This is a test string with old_string.",
			old:      "old_string",
			new:      "new_string",
			expected: "This is a test string with new_string.",
		},
		{
			name:     "Title case replacement",
			input:    "This is a test string with Old_String.",
			old:      "old_string",
			new:      "new_string",
			expected: "This is a test string with New_String.",
		},
		{
			name:     "All upper case replacement",
			input:    "This is a test string with OLD STRING.",
			old:      "old_string",
			new:      "new_string",
			expected: "This is a test string with NEW STRING.",
		},
		{
			name:     "Mixed case replacement",
			input:    "This is a test string with old_string, Old_String and OLD STRING.",
			old:      "old_string",
			new:      "new_string",
			expected: "This is a test string with new_string, New_String and NEW STRING.",
		},
		{
			name:     "Multi line white spaces replacement",
			input:    "This is a test string with old_string, Old_String and OLD\n\n\tSTRING.",
			old:      "old_string",
			new:      "new_string",
			expected: "This is a test string with new_string, New_String and NEW\n\n\tSTRING.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Subvert(tt.input, tt.old, tt.new)
			assert.Equal(t, tt.expected, result)
			// if result != tt.expected {
			// t.Errorf("ToMulti() %q, want %v", result, tt.expected)
			// }
		})
	}
}

/*
func TestRegexps(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		old      string
		new      string
		expected string
	}{
		{
			name:     "happy path",
			input:    "this is a test string with old_string.",
			old:      "old_string",
			new:      "new_string",
			expected: "this is a test string with new_string.",
		},
		{
			name:     "take care of special characters",
			input:    "this is a test string with old_(?)string.",
			old:      "old_(?)string",
			new:      "new_string",
			expected: "this is a test string with new_string.",
		},
		{
			name:     "with title case",
			input:    "this is a test string with old String.",
			old:      "old_string",
			new:      "new_string",
			expected: "this is a test string with New String.",
		},
		{
			name:     "with multiple words",
			input:    "this is a test string with old string with extra words.",
			old:      "old_string_with_extra_words",
			new:      "new_string_and_more",
			expected: "this is a test string with new string and more.",
		},
		{
			name:     "with multiple occurrences",
			input:    "OldString oldstring old_string old string",
			old:      "old_string",
			new:      "new_string",
			expected: "NewString newstring new_string new string",
		},
		{
			name:     "with multiple white spaces",
			input:    "Old    String old\n\n\tstring old\tstring old\nstring",
			old:      "old_string",
			new:      "new_string",
			expected: "New    String new\n\n\tstring new\tstring new\nstring",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := myRE(tt.input, tt.old, tt.new)
			if result != tt.expected {
				t.Errorf("myRE:\n\tGot: %+q\n\tExp: %+q", result, tt.expected)
			}
		})
	}
}
*/
