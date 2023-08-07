package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/slaxor/replace"
	"github.com/spf13/pflag"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

var recursive bool
var all bool
var inplace bool

func init() {
	log.Printf("Args: %v", os.Args)
	pflag.BoolVarP(&recursive, "recursive", "r", false, "edit all files under the given directory, recursively (implies -i)")
	pflag.BoolVarP(&all, "all", "a", false, "edit all files, normally skip hidden files")
	pflag.BoolVarP(&inplace, "inplace", "i", false, "edit files in place")
}

func main() {
	var err error
	var inFile io.Reader
	var i []byte
	var fs []string
	// var ds []string // Directories will be relevant for recursive mode when renaming directories

	pflag.Parse()
	args := pflag.Args()
	if recursive {
		var excl []string
		// inplace = true
		if len(args) < 3 {
			pflag.Usage()
			os.Exit(1)
		}
		if all {
			excl = []string{}
		} else {
			excl = []string{".git", ".svn", ".hg", ".bzr"}
		}
		fs, _, err = replace.Rlist(args[2], excl)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range fs {
			subvertFile(f, args[0], args[1])
		}
	} else {
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
		}

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
}

func subvertFile(fileName string, o, n string) {
	i, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	r := replace.SubvertBytes(i, o, n)
	err = ioutil.WriteFile(fileName, r, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
