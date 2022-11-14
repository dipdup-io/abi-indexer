package vm

import (
	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/abi-indexer/internal/vm/evm"
	"github.com/pkg/errors"
)

// Decoder -
type Decoder interface {
	Methods() ([]storage.Method, error)
	Events() ([]storage.Event, error)
}

// Type -
type Type string

// vm types
const (
	TypeEVM Type = "evm"
)

// VirtualMachine -
type VirtualMachine interface {
	Decoder

	JSONSchema() ([]byte, error)
}

// Config -
type Config struct {
	Type Type `yaml:"type" validate:"required,oneof=evm"`
}

// Factory -
func Factory(typ Type, abi []byte) (VirtualMachine, error) {
	switch typ {
	case TypeEVM:
		return evm.NewVM(abi)
	default:
		return nil, errors.Errorf("unknown virtual machine: %s", typ)
	}
}
