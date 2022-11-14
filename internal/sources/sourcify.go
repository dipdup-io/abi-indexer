package sources

import (
	"context"
	"path/filepath"
	"strconv"
	"time"

	sourcify "github.com/dipdup-net/sourcify-api"
)

// SourcifyConfig -
type SourcifyConfig struct {
	BaseURL string `yaml:"base_url" validate:"required,url"`
	Timeout int    `yaml:"timeout" validate:"required,min=0"`
	ChainID uint64 `yaml:"chain_id" validate:"required,min=1"`
}

// Sourcify -
type Sourcify struct {
	api     *sourcify.API
	chainID string
	timeout time.Duration
}

// NewSourcify -
func NewSourcify(cfg *SourcifyConfig) *Sourcify {
	src := new(Sourcify)

	src.api = sourcify.NewAPI(cfg.BaseURL)
	src.chainID = strconv.FormatUint(cfg.ChainID, 10)
	src.timeout = time.Second * time.Duration(cfg.Timeout)

	return src
}

// Get -
func (s *Sourcify) Get(ctx context.Context, contract string) ([]byte, error) {
	fileTree, err := s.api.GetFiles(ctx, s.chainID, contract)
	if err != nil {
		return nil, err
	}
	for i := range fileTree.Files {
		if ext := filepath.Ext(fileTree.Files[i].Name); ext == ".json" {
			metadata, err := sourcify.ParseMetadata(fileTree.Files[i].Content)
			if err != nil {
				return nil, err
			}
			return metadata.Output.ABI, nil
		}
	}
	return nil, ErrNotFound
}

// List -
func (s *Sourcify) List(ctx context.Context) ([]string, error) {
	contracts, err := s.api.GetContractAddresses(ctx, s.chainID)
	if err != nil {
		return nil, err
	}

	addresses := make([]string, 0)
	addresses = append(addresses, contracts.Full...)
	addresses = append(addresses, contracts.Partial...)

	return addresses, nil
}
