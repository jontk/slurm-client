module github.com/jontk/slurm-client/pkg/streaming

go 1.22.5

toolchain go1.24.5

require (
	github.com/gorilla/websocket v1.5.0
	github.com/jontk/slurm-client v0.0.0-00010101000000-000000000000
)

replace github.com/jontk/slurm-client => ../../
