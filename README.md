# Replace

This project is a Go package and CLI tool to perform smart
case-preserving search and replace operations on text, inspired by the
Subvert command from [vim-abolish](https://github.com/tpope/vim-abolish)
by Tim Pope.

## Installation

To install the package, use `go get`:

```shell
go get github.com/slaxor/replace
```

To install the CLI tool, clone the repository and build the binary:

```shell
git clone https://github.com/slaxor/replace.git
cd replace
go build -o replace .
```

Usage
As a Package
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

As a CLI Tool
Here's a basic example of using the CLI tool:

`replace -o old_world -n new_world -p /path/to/directory`
This will recursively search for occurrences of "old_world" (in various casing styles) in all text files in the specified directory, and replace them with "new_world" (preserving the original casing style).

By default, certain directories (like .git) are excluded from the search. Use the -a flag to include all directories:

`replace -a -o old_world -n new_world -p /path/to/directory`
Credits
The search and replace functionality in this project was inspired by the Subvert command from vim-abolish by Tim Pope. Thank you, Tim, for creating such a useful tool!
