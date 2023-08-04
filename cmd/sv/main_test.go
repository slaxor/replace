package main

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	blnsOnce sync.Once
	blnsFile = "./blns.txt" // Choose a suitable path
	blnsURL  = "https://raw.githubusercontent.com/minimaxir/big-list-of-naughty-strings/master/blns.txt"
)

// getBLNS fetches the BLNS file from blnsURL, but only once.
// If a local copy already exists at blnsFile, it loads from there instead.
// Subsequent calls will return the cached content.
func getBLNS() ([]byte, error) {
	var data []byte
	var err error

	blnsOnce.Do(func() {
		// Try to load from local file
		data, err = ioutil.ReadFile(blnsFile)
		if err == nil {
			return // Successfully loaded from local file
		}

		// Local file not found or couldn't be read, try to download
		resp, err := http.Get(blnsURL)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		// Read the file content
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		// Store the file content locally
		err = ioutil.WriteFile(blnsFile, data, 0644)
	})

	return data, err
}

func TestMainE2EWithBLNS(t *testing.T) {
	// Open the blns file
	file, err := os.Open("blns.txt")
	if err != nil {
		t.Fatalf("Failed to open blns.txt: %v", err)
	}
	defer file.Close()

	// Read the blns file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Each line of the file is a naughty string
		naughty := scanner.Text()

		// Test the main function with the naughty string
		// This is just an example. You may need to adjust this code
		// depending on what exactly you want to test.
		err := main("--old", naughty, "--new", "new_foo", "--path", "/path/to/directory")
		assert.NoError(t, err, "main() failed with arg %q", naughty)
	}

	// Check for errors from the Scanner
	if err := scanner.Err(); err != nil {
		t.Fatalf("Failed to read blns.txt: %v", err)
	}
}

func TestMainE2E(t *testing.T) {
	// Compile the binary
	cmd := exec.Command("go", "build", "-o", "replace")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to compile binary: %v", err)
	}
	defer os.Remove("replace") // Cleanup the binary after the test

	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "e2e")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file with the old string
	filePath := filepath.Join(tempDir, "test.txt")
	err = ioutil.WriteFile(filePath, []byte("This is old_foo"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Run the binary
	cmd = exec.Command("./replace", "--old", "old_foo", "--new", "new_foo", "--path", tempDir)
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	// Check that the file contents have been replaced
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	assert.Equal(t, "This is new_foo", string(contents))
}
