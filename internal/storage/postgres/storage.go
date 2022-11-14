package postgres

import (
	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

// Storage -
type Storage struct {
	Metadata storage.IMetadata
	Methods  storage.IMethod
	Events   storage.IEvent

	db *database.PgGo
}

// Close -
func (s Storage) Close() error {
	return s.db.Close()
}

// BeginTransaction -
func (s *Storage) BeginTransaction() (storage.Transaction, error) {
	tx, err := s.db.DB().Begin()
	if err != nil {
		return nil, err
	}
	return &Transaction{tx}, nil
}
