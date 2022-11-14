package storage

import (
	"context"

	"github.com/pkg/errors"
)

// Models -
type Models interface {
	*Metadata | *Method | *Event
}

// SortOrder -
type SortOrder string

// sort orders
const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// Table -
type Table[M Models] interface {
	Save(ctx context.Context, m M) error
	Update(ctx context.Context, m M) error
	List(ctx context.Context, limit, offset uint64, order SortOrder) ([]M, error)
}

// Transaction -
type Transaction interface {
	Flush(ctx context.Context) error
	Add(ctx context.Context, model any) error
	Update(ctx context.Context, model any) error
	Rollback(ctx context.Context) error
	BulkSave(ctx context.Context, models []any) error
	Close(ctx context.Context) error
}

// Model -
type Model interface {
	TableName() string
}

// ProcessTransactionError -
func ProcessTransactionError(ctx context.Context, err error, tx Transaction) error {
	processorErr := errors.Wrap(err, "transaction error")
	if err := tx.Rollback(ctx); err != nil {
		return errors.Wrap(processorErr, errors.Wrap(err, "rollback").Error())
	}
	return processorErr
}
