package postgres

import (
	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

// Methods -
type Methods struct {
	*Table[*storage.Method]
}

// NewMethods -
func NewMethods(db *database.PgGo) *Methods {
	return &Methods{
		Table: NewTable[*storage.Method](db),
	}
}
