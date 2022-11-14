package grpc

// ServerConfig -
type ServerConfig struct {
	Bind string `yaml:"bind" validate:"required,hostname_port"`
}

// ClientConfig -
type ClientConfig struct {
	ServerAddress string         `yaml:"server_address" validate:"required,hostname_port"`
	Subscriptions *Subscriptions `yaml:"subscriptions" validate:"omitempty"`
}

// Subscriptions -
type Subscriptions struct {
	Metadata bool `yaml:"head,omitempty"`
}

// Config -
type Config struct {
	Server *ServerConfig `yaml:"server" validate:"omitempty"`
	Client *ClientConfig `yaml:"client" validate:"omitempty"`
}
