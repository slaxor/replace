package main

import (
	"fmt"
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
		logMessage := fmt.Sprintf("[Panic Recovery] %s:%d - %v", fn, line, r)
		log.Print(logMessage)
	}
}
func main() {
	defer handlePanic()
	var err error

	var verbose bool
	var inplace bool

	pflag.BoolVarP(&verbose, "verbose", "v", false, "enable verbose mode")
	pflag.BoolVarP(&inplace, "inplace", "i", false, "edit files in place")

	// Parse the flags
	pflag.Parse()

	log.Printf("Verbose mode enabled: %t", verbose)
	log.Printf("Inplace mode enabled: %t", inplace)
	args := pflag.Args()
	log.Printf("Args: %v", args)

	if len(args) < 2 {
		fmt.Printf("Usage: %v <old string> <new string> [<file>]\n", os.Args[0])
		os.Exit(1)
	}
	var i []byte
	o := args[0]
	n := args[1]

	inFile := os.Stdin
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
	if inplace {
		err = ioutil.WriteFile(args[3], r, 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("%s", r)
	}
}