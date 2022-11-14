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

// GetByMethod -
func (m *Metadata) GetByMethod(ctx context.Context, signature string, limit, offset uint64, order storage.SortOrder) ([]*storage.Metadata, error) {
	var methods []*storage.Method
	query := m.db.DB().ModelContext(ctx, &methods).
		Relation("Metadata").
		Where("signature = ?", signature).
		Where("metadata_id is not null")

	addPagination(query, limit, offset, order)

	if err := query.Select(); err != nil {
		return nil, err
	}

	response := make([]*storage.Metadata, len(methods))
	for i := range methods {
		response[i] = methods[i].Metadata
	}
	return response, nil
}

// GetByTopic -
func (m *Metadata) GetByTopic(ctx context.Context, topic string, limit, offset uint64, order storage.SortOrder) ([]*storage.Metadata, error) {
	var events []*storage.Event
	query := m.db.DB().ModelContext(ctx, &events).
		Relation("Metadata").
		Where("topic = ?", topic).
		Where("metadata_id is not null")

	addPagination(query, limit, offset, order)

	if err := query.Select(); err != nil {
		return nil, err
	}

	response := make([]*storage.Metadata, len(events))
	for i := range events {
		response[i] = events[i].Metadata
	}
	return response, nil
}
