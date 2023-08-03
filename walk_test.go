package replace

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func syncfs(fd int) error {
	_, _, err := syscall.Syscall(306, uintptr(fd), 0, 0)
	if err != 0 {
		return err
	}
	return nil
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

func TestWalkAndSubvert(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "old_foo")
	err := os.Mkdir(subDir, 0755)
	assert.NoError(t, err)

	testFile := filepath.Join(subDir, "test.txt")
	err = ioutil.WriteFile(testFile, []byte("This is a string with old_foo and OldFoo"), 0644)
	assert.NoError(t, err)

	excludedDir := filepath.Join(dir, ".git")
	err = os.Mkdir(excludedDir, 0755)
	assert.NoError(t, err)

	excludedFile := filepath.Join(excludedDir, "test.txt")
	err = ioutil.WriteFile(excludedFile, []byte("This is a string with old_foo and OldFoo"), 0644)
	assert.NoError(t, err)

	err = WalkAndSubvert(dir, []string{".git"}, "old_foo", "new_foo")
	assert.NoError(t, err)
	// t.Run("Test subvert function application", func(t *testing.T) {
	// })

	t.Run("Test directory renaming", func(t *testing.T) {
		t.Skip("Not implemented")
		// Check if the directory was renamed
		_, err = os.Stat(subDir)
		assert.Error(t, err) // Expect an error because the directory should have been renamed

		newSubDir := filepath.Join(dir, "new_foo")
		_, err = os.Stat(newSubDir)
		assert.NoError(t, err)
	})

	t.Run("Test file content replacement", func(t *testing.T) {
		// Check if the file content was replaced
		content, err := ioutil.ReadFile(testFile)
		assert.NoError(t, err)
		assert.Equal(t, "This is a string with new_foo and NewFoo", string(content))
	})

	t.Run("Test exclusion of directories", func(t *testing.T) {
		t.Skip("Not implemented")
		// Check if the excluded directory was not renamed
		_, err = os.Stat(excludedDir)
		assert.NoError(t, err)
	})

	t.Run("Test exclusion of file content replacement", func(t *testing.T) {
		// Check if the excluded file content was not replaced
		excludedContent, err := ioutil.ReadFile(excludedFile)
		assert.NoError(t, err)
		assert.Equal(t, "This is a string with old_foo and OldFoo", string(excludedContent))
	})

	t.Run("directory is renamed", func(t *testing.T) {
		t.Skip("Not implemented")
		_, err := os.Stat(subDir)
		assert.Error(t, err)

		newSubDir := filepath.Join(dir, "new_foo")
		_, err = os.Stat(newSubDir)
		assert.NoError(t, err)
	})

	t.Run("file is renamed and content is replaced", func(t *testing.T) {
		t.Skip("Not implemented")
		newTestFile := filepath.Join(dir, "new_foo", "new_foo.txt")
		content, err := ioutil.ReadFile(newTestFile)
		assert.NoError(t, err)
		assert.Equal(t, "This is a string with new_foo and NewFoo", string(content))
	})
}
