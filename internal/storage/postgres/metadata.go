package postgres

import (
	"context"

	models "github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Metadata -
type Metadata struct {
	*postgres.Table[*models.Metadata]
}

// NewMetadata -
func NewMetadata(db *database.PgGo) *Metadata {
	return &Metadata{
		Table: postgres.NewTable[*models.Metadata](db),
	}
}

// GetByAddress -
func (m *Metadata) GetByAddress(ctx context.Context, address string) (*models.Metadata, error) {
	var response models.Metadata
	err := m.DB().ModelContext(ctx, &response).Where("contract = ?", address).First()
	return &response, err
}

// GetByMethod -
func (m *Metadata) GetByMethod(ctx context.Context, signature string, limit, offset uint64, order storage.SortOrder) ([]*models.Metadata, error) {
	var methods []*models.Method
	query := m.DB().ModelContext(ctx, &methods).
		Relation("Metadata").
		Where("signature = ?", signature).
		Where("metadata_id is not null")

	postgres.Pagination(query, limit, offset, order)

	if err := query.Select(); err != nil {
		return nil, err
	}

	response := make([]*models.Metadata, len(methods))
	for i := range methods {
		response[i] = methods[i].Metadata
	}
	return response, nil
}

// GetByTopic -
func (m *Metadata) GetByTopic(ctx context.Context, topic string, limit, offset uint64, order storage.SortOrder) ([]*models.Metadata, error) {
	var events []*models.Event
	query := m.DB().ModelContext(ctx, &events).
		Relation("Metadata").
		Where("signature_id = ?", topic).
		Where("metadata_id is not null")

	postgres.Pagination(query, limit, offset, order)

	if err := query.Select(); err != nil {
		return nil, err
	}

	response := make([]*models.Metadata, len(events))
	for i := range events {
		response[i] = events[i].Metadata
	}
	return response, nil
}
