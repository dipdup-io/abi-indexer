package main

import (
	"github.com/dipdup-net/abi-indexer/pkg/modules/grpc"
	"github.com/dipdup-net/abi-indexer/pkg/modules/metadata"
	"github.com/dipdup-net/go-lib/config"
)

// Config -
type Config struct {
	config.Config `yaml:",inline"`
	LogLevel      string          `yaml:"log_level" validate:"omitempty,oneof=debug trace info warn error fatal panic"`
	Metadata      metadata.Config `yaml:"metadata"`
	GRPC          grpc.Config     `yaml:"grpc"`
}

// Substitute -
func (c *Config) Substitute() error {
	if err := c.Config.Substitute(); err != nil {
		return err
	}
	return nil
}
