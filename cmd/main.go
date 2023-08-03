package main

import (
	"fmt"
	"os"

	"github.com/slaxor/replace"
	flag "github.com/spf13/pflag"
)

func main() {
	var oldStr, newStr, path string
	var all bool
	var exclusions []string

	// Define the flags
	flag.StringVarP(&oldStr, "old", "o", "", "The old string to replace")
	flag.StringVarP(&newStr, "new", "n", "", "The new string to replace with")
	flag.StringVarP(&path, "path", "p", ".", "The path of directory to perform replacements")
	flag.BoolVarP(&all, "all", "a", false, "Include system directories")

	// Parse the flags
	flag.Parse()

	if !all {
		exclusions = []string{".git"}
	}

	err := replace.WalkAndSubvert(path, exclusions, oldStr, newStr)
	if err != nil {
		fmt.Printf("An error occurred: %v", err)
		os.Exit(1)
	}
}
