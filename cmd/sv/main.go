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
		logMessage := fmt.Sprintf("[Panic] %s:%d - %v", fn, line, r)
		log.Print(logMessage)
	}
}
func main() {
	defer handlePanic()
	var err error
	var inplace bool
	pflag.BoolVarP(&inplace, "inplace", "i", false, "edit files in place")
	pflag.Parse()

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
	if inplace && len(args) == 3 {
		err = ioutil.WriteFile(args[2], r, 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("%s", r)
	}
}
