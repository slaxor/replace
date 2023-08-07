# Replace

This project is a Go package and CLI tool to perform smart
case-preserving search and replace operations on text, inspired by the
Subvert command from [vim-abolish](https://github.com/tpope/vim-abolish)
by Tim Pope.

Smart means given `sv old_string new_string` replaces old_string ->
new_string (of course), OldString -> NewString, old string ->
new string, OLDSTRING -> NEWSTRING, etc.

perhaps some of you find this helpful too. Let me know.

## Installation

To install the package, use `go get`:

```shell
go get github.com/slaxor/replace
```

If you don't care for the package. You can just install the command.
To install the CLI tool, clone the repository and build the binary:

```shell
go install github.com/slaxor/replace/cmd/sv@master
```

Or even go to the realease page and fetch the binary of your choice

## Usage

### As a Package

Here's a basic example of using the package in your own Go code:

```go
import "github.com/slaxor/replace"

func main() {
    input := "Hello, old world!"
    old := "old_world"
    new := "new_world"
    output := replace.Subvert(input, old, new)
    // output is now "Hello, new world!"
}
```

### As a CLI Tool

Here's a basic example of using the CLI tool:

`sv -r old_world new_world /path/to/directory`

This will recursively search for occurrences of "old_world" (in various casing styles) in all text files in the specified directory, and replace them with "new_world" (preserving the original casing style).
By default, certain directories (like .git) are excluded from the search. Use the -a flag to include all directories.

## Todo

1. [ ] option to rename files in a recursive replacement
2. [ ] make exclusions user accessable
3. [ ] recursive pattern matching on file exclusions

## Contribute

You know the drill I guess

-   Fork
-   Branch
-   Make Pull Request
-   I will review
-   I will merge or tell you why not
-   Done
