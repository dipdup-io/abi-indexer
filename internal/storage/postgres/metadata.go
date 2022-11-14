package postgres

import (
	"context"

	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

// Metadata -
type Metadata struct {
	*Table[*storage.Metadata]
}

// NewMetadata -
func NewMetadata(db *database.PgGo) *Metadata {
	return &Metadata{
		Table: NewTable[*storage.Metadata](db),
	}
}

// GetByAddress -
func (m *Metadata) GetByAddress(ctx context.Context, address string) (*storage.Metadata, error) {
	var response storage.Metadata
	err := m.db.DB().ModelContext(ctx, &response).Where("contract = ?", address).First()
	return &response, err
}

// GetByMethodSinature -
func (m *Metadata) GetByMethodSinature(ctx context.Context, signature string, limit, offset uint64, order storage.SortOrder) ([]*storage.Metadata, error) {
	return nil, nil
}

// GetByTopic -
func (m *Metadata) GetByTopic(ctx context.Context, topic string, limit, offset uint64, order storage.SortOrder) ([]*storage.Metadata, error) {
	return nil, nil
}
