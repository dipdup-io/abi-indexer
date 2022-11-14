package postgres

import (
	"context"

	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

// Table -
type Table[M storage.Models] struct {
	db *database.PgGo
}

// NewTable -
func NewTable[M storage.Models](db *database.PgGo) *Table[M] {
	return &Table[M]{db}
}

// Save -
func (s *Table[M]) Save(ctx context.Context, m M) error {
	_, err := s.db.DB().ModelContext(ctx, m).Returning("id").Insert()
	return err
}

// Update -
func (s *Table[M]) Update(ctx context.Context, m M) error {
	_, err := s.db.DB().ModelContext(ctx, m).WherePK().Update()
	return err
}

// List -
func (s *Table[M]) List(ctx context.Context, limit, offset uint64, order storage.SortOrder) ([]M, error) {
	var models []M
	query := s.db.DB().ModelContext(ctx, &models)

	if limit == 0 {
		limit = 10
	}

	query.Limit(int(limit)).Offset(int(offset))

	switch order {
	case storage.SortOrderAsc:
		query.Order("id asc")
	case storage.SortOrderDesc:
		query.Order("id desc")
	default:
		query.Order("id asc")
	}

	err := query.Select(&models)
	return models, err
}
