package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/slaxor/replace"
	"github.com/spf13/pflag"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

var recursive bool
var all bool
var inplace bool
var enableRegex bool

func init() {
	pflag.BoolVarP(&recursive, "recursive", "r", false, "edit all files under the given directory, recursively (implies -i)")
	pflag.BoolVarP(&all, "all", "a", false, "edit all files, normally skip hidden files")
	pflag.BoolVarP(&inplace, "inplace", "i", false, "edit files in place")
	pflag.BoolVarP(&enableRegex, "regex", "e", false, "enable for old string to be a regular expression")
}

func main() {
	var err error
	var fs []string
	// var ds []string // Directories will be relevant for recursive mode when renaming directories
	pflag.Parse()
	args := pflag.Args()
	if len(args) < 2 {
		pflag.Usage()
		os.Exit(255)
	}
	o := args[0]
	if !enableRegex {
		o = regexp.QuoteMeta(o)
	}
	n := args[1]

	if recursive { // implies inplace = true
		var excl []string
		if len(args) < 3 {
			pflag.Usage()
			log.Printf("recursive: %t, arglen: %d", recursive, len(args))
			os.Exit(255)
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
			subvertFileInplace(f, o, n)
		}
	} else {
		switch {
		case len(args) == 2:
			subvertStdin(o, n)
		case len(args) > 2:
			var fn func(string, string, string)
			if inplace {
				fn = subvertFileInplace
			} else {
				fn = subvertFile
			}

			for _, f := range args[2:] {
				fn(f, o, n)
			}
		}
	}
}

func subvertStdin(o, n string) {
	i, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	r := replace.SubvertBytes(i, o, n)
	fmt.Printf("%s", r)
}

func subvertFile(f, o, n string) {
	i, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}
	r := replace.SubvertBytes(i, o, n)
	fmt.Printf("%s", r)
}

func subvertFileInplace(f, o, n string) {
	i, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}
	r := replace.SubvertBytes(i, o, n)
	err = ioutil.WriteFile(f, r, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
