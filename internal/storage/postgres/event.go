package postgres

import (
	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

// Events -
type Events struct {
	*Table[*storage.Event]
}

// NewEvents -
func NewEvents(db *database.PgGo) *Events {
	return &Events{
		Table: NewTable[*storage.Event](db),
	}
}
