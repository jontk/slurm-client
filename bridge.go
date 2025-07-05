package slurm

import (
	"github.com/jontk/slurm-client/internal/factory"
)

// clientBridge adapts the factory.SlurmClient to our main SlurmClient interface
type clientBridge struct {
	factoryClient factory.SlurmClient
}

// newClientBridge creates a bridge between factory and main interfaces
func newClientBridge(factoryClient factory.SlurmClient) SlurmClient {
	return &clientBridge{factoryClient: factoryClient}
}

// Version returns the API version
func (c *clientBridge) Version() string {
	return c.factoryClient.Version()
}

// Jobs returns the JobManager
func (c *clientBridge) Jobs() JobManager {
	return &jobManagerBridge{c.factoryClient.Jobs()}
}

// Nodes returns the NodeManager
func (c *clientBridge) Nodes() NodeManager {
	return &nodeManagerBridge{c.factoryClient.Nodes()}
}

// Partitions returns the PartitionManager
func (c *clientBridge) Partitions() PartitionManager {
	return &partitionManagerBridge{c.factoryClient.Partitions()}
}

// Info returns the InfoManager
func (c *clientBridge) Info() InfoManager {
	return &infoManagerBridge{c.factoryClient.Info()}
}

// Close closes the client
func (c *clientBridge) Close() error {
	return c.factoryClient.Close()
}

// Manager bridges to convert between factory and main interfaces

type jobManagerBridge struct {
	factory.JobManager
}

type nodeManagerBridge struct {
	factory.NodeManager
}

type partitionManagerBridge struct {
	factory.PartitionManager
}

type infoManagerBridge struct {
	factory.InfoManager
}