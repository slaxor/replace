package replace

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupFiles(t *testing.T, tempDir string, files []string) {
	var err error
	for _, file := range files[:5] { // for regular files
		err := os.MkdirAll(filepath.Join(tempDir, filepath.Dir(file)), 0755)
		if err != nil {
			t.Fatal(err)
		}
		_, err = os.Create(filepath.Join(tempDir, file))
		if err != nil {
			t.Fatal(err)
		}
	}

	// Creating a device file
	// err := syscall.Mknod(filepath.Join(tempDir, files[5]), 0600, int(unix.Mkdev(uint32(1), uint32(3)))) // /dev/null device file
	// if err != nil {
	// t.Fatal(err)
	// }

	// Creating a pipe
	err = syscall.Mkfifo(filepath.Join(tempDir, files[5]), 0600)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRlist(t *testing.T) {
	tempDir := t.TempDir()

	// Create test files and directories
	allFiles := []string{
		"test1.txt",
		"test2.txt",
		".git/ignore.txt",
		"subdir/test3.txt",
		"subdir/.hidden/test4.txt",
		"pipe",
	}

	setupFiles(t, tempDir, allFiles)

	cases := []struct {
		name          string
		exclusions    []string
		expectedFiles []string // Expected files
	}{
		{
			name:       "No exclusions",
			exclusions: []string{},
			expectedFiles: []string{
				"test1.txt",
				"test2.txt",
				".git/ignore.txt",
				"subdir/test3.txt",
				"subdir/.hidden/test4.txt",
			},
		},
		{
			name:       "Exclude .git",
			exclusions: []string{".git"},
			expectedFiles: []string{
				"test1.txt",
				"test2.txt",
				"subdir/test3.txt",
				"subdir/.hidden/test4.txt",
			},
		},
		{
			name:       "Exclude hidden",
			exclusions: []string{"*.hidden"},
			expectedFiles: []string{
				"test1.txt",
				"test2.txt",
				".git/ignore.txt",
				"subdir/test3.txt",
			},
		},
		{
			name:          "Exclude all",
			exclusions:    []string{"*"},
			expectedFiles: []string{}, // all files excluded
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files, _, err := Rlist(tempDir, tc.exclusions)
			if err != nil {
				t.Fatal(err)
			}
			expectedFiles := make([]string, len(tc.expectedFiles))
			for i, file := range tc.expectedFiles {
				expectedFiles[i] = filepath.Join(tempDir, file)
			}
			assert.ElementsMatch(t, expectedFiles, files)
		})
	}
}
