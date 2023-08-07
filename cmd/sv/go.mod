module github.com/slaxor/replace/cmd/sv

go 1.20

require (
	github.com/slaxor/replace v0.0.0-20230804121628-1c9ef3c50562
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/slaxor/replace => ../..
