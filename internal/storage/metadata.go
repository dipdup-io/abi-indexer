package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IMetadata -
type IMetadata interface {
	storage.Table[*Metadata]

	GetByAddress(ctx context.Context, address string) (*Metadata, error)
	GetByMethod(ctx context.Context, signature string, limit, offset uint64, order storage.SortOrder) ([]*Metadata, error)
	GetByTopic(ctx context.Context, topic string, limit, offset uint64, order storage.SortOrder) ([]*Metadata, error)
}

// Metadata -
type Metadata struct {
	// nolint
	tableName struct{} `pg:"metadata"`

	ID uint64

	Contract   string `pg:",unique:metadata_contract,notnull"`
	Metadata   []byte
	JSONSchema []byte
}

// TableName -
func (Metadata) TableName() string {
	return "metadata"
}
