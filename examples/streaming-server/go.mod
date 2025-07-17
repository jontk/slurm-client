module github.com/jontk/slurm-client/examples/streaming-server

go 1.19

require (
	github.com/jontk/slurm-client v0.0.0-00010101000000-000000000000
	github.com/jontk/slurm-client/pkg/streaming v0.0.0-00010101000000-000000000000
)

require github.com/gorilla/websocket v1.5.0 // indirect

replace github.com/jontk/slurm-client => ../../

replace github.com/jontk/slurm-client/pkg/streaming => ../../pkg/streaming