package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/slaxor/replace"
	"github.com/spf13/pflag"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func handlePanic() {
	if r := recover(); r != nil {
		_, fn, line, _ := runtime.Caller(3) // 2 steps up the stack frame
		log.Printf("[Panic] %s:%d - %v", fn, line, r)
	}
}

// func argCheck(args []string) (string, string, string) {
// if len(args) < 2 {
// pflag.Usage()
// fmt.Printf("Usage: %s [-i] <old string> <new string> [<file>]\n", os.Args[0])
// os.Exit(1)
// }
// }

func main() {
	defer handlePanic()
	var err error
	var inFile io.Reader
	var recursive bool
	var all bool
	var inplace bool
	pflag.BoolVarP(&recursive, "recursive", "r", false, "edit all files under the given directory, recursively (implies -i)")
	pflag.BoolVarP(&all, "all", "a", false, "edit all files, normally skip hidden files")
	pflag.BoolVarP(&inplace, "inplace", "i", false, "edit files in place")
	pflag.Parse()

	args := pflag.Args()
	if recursive {
		inplace = true
	}
	switch {
	case len(args) < 2:
		pflag.Usage()
		os.Exit(1)
	case len(args) == 2:
		inFile = os.Stdin
	case len(args) == 3:
		inFile, err = os.Open(args[2])
		if err != nil {
			log.Fatal(err)
		}

	default:

	}

	var i []byte
	o := args[0]
	n := args[1]

	if len(args) == 3 {
		inFile, err = os.Open(args[2])
		if err != nil {
			log.Fatal(err)
		}
	}
	i, err = ioutil.ReadAll(inFile)
	if err != nil {
		log.Fatal(err)
	}
	r := replace.SubvertBytes(i, o, n)
	if inplace && len(args) == 3 {
		err = ioutil.WriteFile(args[2], r, 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("%s", r)
	}
}
