package postgres

import (
	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Events -
type Events struct {
	*postgres.Table[*storage.Event]
}

// NewEvents -
func NewEvents(db *database.PgGo) *Events {
	return &Events{
		Table: postgres.NewTable[*storage.Event](db),
	}
}
