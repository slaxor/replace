package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	blnsOnce sync.Once
	blnsFile = "./blns.txt" // Choose a suitable path
	blnsURL  = "https://raw.githubusercontent.com/minimaxir/big-list-of-naughty-strings/master/blns.txt"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// getBLNS fetches the BLNS file from blnsURL, but only once.
// If a local copy already exists at blnsFile, it loads from there instead.
// Subsequent calls will return the cached content.
func getBLNS() ([]byte, error) {
	var data []byte
	var err error
	blnsOnce.Do(func() {
		data, err = ioutil.ReadFile(blnsFile)
		if err == nil {
			return // Successfully loaded from local file
		}
		resp, err := http.Get(blnsURL)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		err = ioutil.WriteFile(blnsFile, data, 0644)
	})
	return data, err
}

func TestStdinExec(t *testing.T) {
	t.Skip("Regularly testing test helpers seems a bit silly")
	stdinExec(func() {
		si, err := ioutil.ReadAll(os.Stdin)
		assert.NoError(t, err)
		log.Printf("%T %[1]v", si)
	}, strings.Repeat("in ", 10))
}

func stdinExec(fn func(), input string) {
	orig := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	defer func() { os.Stdin = orig }()
	w.Write([]byte(input))
	go func() {
	}()
	w.Close()
	fn()
}

func TestStdoutAndStderrExec(t *testing.T) {
	t.Skip("Regularly testing test helpers seems a bit silly")
	otest := strings.Repeat("out ", 10)
	etest := strings.Repeat("err ", 10)
	o, e := stdoutAndStderrExec(func() {
		fmt.Fprint(os.Stdout, otest)
		fmt.Fprint(os.Stderr, etest)
	})
	assert.Equal(t, otest, o)
	assert.Equal(t, etest, e)
}

func stdoutAndStderrExec(fn func()) (string, string) {
	origOut := os.Stdout
	origErr := os.Stderr
	ro, wo, _ := os.Pipe()
	re, we, _ := os.Pipe()
	os.Stdout = wo
	os.Stderr = we
	defer func() {
		os.Stdout = origOut
		os.Stderr = origErr
	}()
	fn()
	outC := make(chan string)
	errC := make(chan string)
	go func() {
		var oBuf bytes.Buffer
		var eBuf bytes.Buffer
		io.Copy(&oBuf, ro)
		io.Copy(&eBuf, re)
		outC <- oBuf.String()
		errC <- eBuf.String()
	}()
	wo.Close()
	we.Close()
	return <-outC, <-errC
}

func TestStdioExec(t *testing.T) {
	t.Skip("Regularly testing test helpers seems a bit silly")
	otest := strings.Repeat("out ", 10)
	etest := strings.Repeat("err ", 10)
	so, se := stdioExec(func() {
		i, _ := ioutil.ReadAll(os.Stdin)
		fmt.Fprintf(os.Stdout, "%s", i[:len(otest)])
		fmt.Fprintf(os.Stderr, "%s", i[len(otest):])
	}, otest+etest)
	assert.Equal(t, otest, so)
	assert.Equal(t, etest, se)
}

func stdioExec(fn func(), i string) (string, string) {
	return stdoutAndStderrExec(func() {
		stdinExec(fn, i)
	})
}

func TestStdio(t *testing.T) {
	os.Args = []string{"sv", "old_string", "new_string"}
	so, se := stdioExec(main, "old_string")
	assert.Equal(t, "new_string", so)
	assert.Equal(t, "", se)
}

func TestFile(t *testing.T) {
	testFile := filepath.Join(t.TempDir(), "test.txt")
	f, err := os.Create(testFile)
	assert.NoError(t, err)
	f.WriteString("old_string")
	os.Args = []string{"sv", "old_string", "new_string", testFile}

	so, se := stdioExec(main, "old_string")
	assert.Equal(t, "new_string", so)
	assert.Equal(t, "", se)
}

func TestFileInplace(t *testing.T) {
	testFile := filepath.Join(t.TempDir(), "test.txt")
	f, err := os.Create(testFile)
	assert.NoError(t, err)
	f.WriteString("old_string")
	os.Args = []string{"sv", "-i", "old_string", "new_string", testFile}

	so, se := stdoutAndStderrExec(main)
	newContent, _ := os.ReadFile(testFile)
	assert.Equal(t, "", so)
	assert.Equal(t, "", se)
	assert.Equal(t, "new_string", string(newContent))
}

func TestRecursive(t *testing.T) {
	tempDir := t.TempDir()
	testFiles := []string{
		"test1.txt",
		"testdir/test1.txt",
		"testdir/testsub/test1.txt",
		".git/test1.txt",
		".git/testdir/test1.txt",
		".regularhidden/testdir/testsub/test1.txt",
	}
	mkFakeProject(t, "OldString", tempDir, testFiles)
	os.Args = []string{"sv", "-r", "old_string", "new_string", tempDir}

	so, se := stdoutAndStderrExec(main)
	assert.Equal(t, "", so)
	assert.Equal(t, "", se)
	for _, tf := range testFiles {
		f := filepath.Join(tempDir, tf)
		newContent, _ := os.ReadFile(f)
		if strings.HasPrefix(tf, ".git/") {
			assert.Equalf(t, "OldString", string(newContent), "File %s", f)
		} else {
			assert.Equalf(t, "NewString", string(newContent), "File %s", f)
		}
	}
}

func TestRecursiveAll(t *testing.T) {
	tempDir := t.TempDir()
	testFiles := []string{
		"test1.txt",
		"testdir/test1.txt",
		"testdir/testsub/test1.txt",
		".git/test1.txt",
		".git/testdir/test1.txt",
		".regularhidden/testdir/testsub/test1.txt",
	}
	mkFakeProject(t, "OldString", tempDir, testFiles)
	os.Args = []string{"sv", "-r", "-a", "old_string", "new_string", tempDir}

	so, se := stdoutAndStderrExec(main)
	assert.Equal(t, "", so)
	assert.Equal(t, "", se)
	for _, tf := range testFiles {
		f := filepath.Join(tempDir, tf)
		newContent, _ := os.ReadFile(f)
		assert.Equalf(t, "NewString", string(newContent), "File %s", f)
	}
}

func mkFakeProject(t *testing.T, testString string, tempDir string, files []string) {
	for _, file := range files {
		err := os.MkdirAll(filepath.Join(tempDir, filepath.Dir(file)), 0755)
		if err != nil {
			t.Fatal(err)
		}
		f, err := os.Create(filepath.Join(tempDir, file))
		if err != nil {
			t.Fatal(err)
		}
		f.WriteString(testString)
		f.Close()
	}
}

func resetArgs() {
	recursive = false
	all = false
	inplace = false
	enableRegex = false
}

func TestRegex(t *testing.T) {
	tests := []struct {
		name string
		args []string
		i    string
		exp  string
	}{
		{
			name: "enabled regex",
			args: []string{"sv", "-e", "o[a-z]+_s[a-z]+", "new_string"},
			i:    "open shut, outward_southward, overcomeSucceed, OvalShape, output signal",
			exp:  "new string, new_string, newString, NewString, new string",
		},
		{
			name: "diabled regex",
			args: []string{"sv", "o[a-z]+_s[a-z]+", "new_string"},
			i:    "open shut, outward_southward, overcomeSucceed, OvalShape, output signal",
			exp:  "open shut, outward_southward, overcomeSucceed, OvalShape, output signal",
		},
		{
			name: "diabled regex with a match anyway",
			args: []string{"sv", "o[a-z]+_s[a-z]+", "new_string"},
			i:    "open shut, outward_southward, o[a-z]+ s[a-z]+, OvalShape, output signal",
			exp:  "open shut, outward_southward, new string, OvalShape, output signal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetArgs()
			os.Args = tt.args
			so, se := stdioExec(main, tt.i)
			assert.Equal(t, tt.exp, so)
			assert.Equal(t, "", se)
		})
	}
}
