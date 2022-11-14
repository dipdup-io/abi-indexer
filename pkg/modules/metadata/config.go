package metadata

import (
	"github.com/dipdup-net/abi-indexer/internal/sources"
	"github.com/dipdup-net/abi-indexer/internal/vm"
)

// Config -
type Config struct {
	SourceType   sources.Type              `yaml:"source_type" validate:"required,oneof=fs sourcify"`
	ThreadsCount int                       `yaml:"threads_count" validate:"omitempty,min=1"`
	VM           *vm.Config                `yaml:"vm"`
	Sourcify     *sources.SourcifyConfig   `yaml:"sourcify"`
	FS           *sources.FileSystemConfig `yaml:"fs"`
}
