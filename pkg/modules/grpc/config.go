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
	Head bool        `yaml:"head,omitempty"`
	Logs *LogsConfig `yaml:"logs"`
	Txs  *TxsConfig  `yaml:"txs"`
}

// IsEmpty -
func (s *Subscriptions) IsEmpty() bool {
	return !s.Head && s.Logs == nil && s.Txs == nil
}

// LogsConfig -
type LogsConfig struct {
	Contracts []string `yaml:"contracts"`
	Topics    []string `yaml:"topics"`
}

// Logs -
type TxsConfig struct {
	From    []string `yaml:"from"`
	To      []string `yaml:"to"`
	Methods []string `yaml:"methods"`
}

// Config -
type Config struct {
	Server *ServerConfig `yaml:"server" validate:"omitempty"`
	Client *ClientConfig `yaml:"client" validate:"omitempty"`
}
