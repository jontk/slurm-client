module github.com/jontk/slurm-client/cmd/slurm-cli

go 1.22.5

toolchain go1.24.5

require github.com/spf13/cobra v1.10.2

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.6 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/jontk/slurm-client => ../..
