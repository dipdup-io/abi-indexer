package storage

import "github.com/dipdup-net/indexer-sdk/pkg/storage"

// IEvent -
type IEvent interface {
	storage.Table[*Event]
}

// Event -
type Event struct {
	// nolint
	tableName struct{} `pg:"events"`

	ID uint64

	Name        string
	Signature   string
	SignatureID []byte
	MetadataID  uint64
	Anonymous   bool `pg:"default:false"`

	Metadata *Metadata `pg:",rel:has-one"`
}

// TableName -
func (Event) TableName() string {
	return "events"
}
