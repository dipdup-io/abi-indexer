package sources

import (
	"context"

	"github.com/pkg/errors"
)

// errors
var (
	ErrNotFound = errors.New("metadata not found")
)

// Source -
type Source interface {
	Get(ctx context.Context, contract string) ([]byte, error)
	List(ctx context.Context) ([]string, error)
}

// Data -
type Data struct {
	Contract string
	Metadata []byte
}

// Type -
type Type string

// types
const (
	FSType       Type = "fs"
	SourcifyType Type = "sourcify"
)

// FactoryParams -
type FactoryParams struct {
	Sourcify *SourcifyConfig
	FS       *FileSystemConfig
}

// Factory -
func Factory(typ Type, params FactoryParams) (Source, error) {
	var abiSource Source
	switch typ {
	case FSType:
		if params.FS == nil {
			return nil, errors.New("you have to set 'fs_abi_source' section for file system ABI source")
		}
		abiSource = NewFileSystem(params.FS.Dir)
	case SourcifyType:
		if params.Sourcify == nil {
			return nil, errors.New("you have to set 'sourcify' section for Sourcify as ABI source")
		}
		abiSource = NewSourcify(params.Sourcify)
	default:
		return nil, errors.Errorf("invalid ABI source type: %s", typ)
	}
	return abiSource, nil
}
