module github.com/jontk/slurm-client/cmd/slurm-cli

go 1.21

require (
	github.com/jontk/slurm-client v0.0.0
	github.com/spf13/cobra v1.8.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace github.com/jontk/slurm-client => ../..