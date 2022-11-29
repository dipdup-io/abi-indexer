package grpc

import "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"

// ClientConfig -
type ClientConfig struct {
	ServerAddress string         `yaml:"server_address" validate:"required"`
	Subscriptions *Subscriptions `yaml:"subscriptions" validate:"omitempty"`
}

// Subscriptions -
type Subscriptions struct {
	Metadata bool `yaml:"head,omitempty"`
}

// Config -
type Config struct {
	Server *grpc.ServerConfig `yaml:"server" validate:"omitempty"`
	Client *ClientConfig      `yaml:"client" validate:"omitempty"`
}
