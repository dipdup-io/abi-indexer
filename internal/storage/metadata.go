package storage

import "context"

// IMetadata -
type IMetadata interface {
	Table[*Metadata]

	GetByAddress(ctx context.Context, address string) (*Metadata, error)
	GetByMethodSinature(ctx context.Context, signature string, limit, offset uint64, order SortOrder) ([]*Metadata, error)
	GetByTopic(ctx context.Context, topic string, limit, offset uint64, order SortOrder) ([]*Metadata, error)
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
