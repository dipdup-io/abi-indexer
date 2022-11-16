package postgres

import (
	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Methods -
type Methods struct {
	*postgres.Table[*storage.Method]
}

// NewMethods -
func NewMethods(db *database.PgGo) *Methods {
	return &Methods{
		Table: postgres.NewTable[*storage.Method](db),
	}
}
