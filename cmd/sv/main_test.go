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
	o := <-outC
	e := <-errC
	return o, e
}

func TestStdioExec(t *testing.T) {
	t.Skip("Regularly testing test helpers seems a bit silly")
	otest := strings.Repeat("out ", 10)
	etest := strings.Repeat("err ", 10)
	o, e := stdioExec(func() {
		i, _ := ioutil.ReadAll(os.Stdin)
		fmt.Fprintf(os.Stdout, "%s", i[:len(otest)])
		fmt.Fprintf(os.Stderr, "%s", i[len(otest):])
	}, otest+etest)
	assert.Equal(t, otest, o)
	assert.Equal(t, etest, e)
}

func stdioExec(fn func(), i string) (string, string) {
	return stdoutAndStderrExec(func() {
		stdinExec(fn, i)
	})
}

func TestStdio(t *testing.T) {
	os.Args = []string{"sv", "old_string", "new_string"}
	o, e := stdioExec(main, "old_string")
	assert.Equal(t, "new_string", o)
	assert.Equal(t, "", e)
}

func TestFile(t *testing.T) {
	testFile := filepath.Join(t.TempDir(), "test.txt")
	f, err := os.Create(testFile)
	assert.NoError(t, err)
	f.WriteString("old_string")
	os.Args = []string{"sv", "old_string", "new_string", testFile}

	o, e := stdioExec(main, "old_string")
	assert.Equal(t, "new_string", o)
	assert.Equal(t, "", e)
}

func TestFileInplace(t *testing.T) {
	testFile := filepath.Join(t.TempDir(), "test.txt")
	f, err := os.Create(testFile)
	assert.NoError(t, err)
	f.WriteString("old_string")
	os.Args = []string{"sv", "-i", "old_string", "new_string", testFile}

	o, e := stdoutAndStderrExec(main)
	newContent, _ := os.ReadFile(testFile)
	assert.Equal(t, "", o)
	assert.Equal(t, "", e)
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

	o, e := stdoutAndStderrExec(main)
	assert.Equal(t, "", o)
	assert.Equal(t, "", e)
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

	o, e := stdoutAndStderrExec(main)
	assert.Equal(t, "", o)
	assert.Equal(t, "", e)
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
