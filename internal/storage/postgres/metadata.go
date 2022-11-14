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
