package subscriptions

import (
	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/abi-indexer/pkg/modules/grpc/pb"
)

// Head -
type Metadata struct {
	data chan *pb.Metadata
}

// NewHead -
func NewMetadata() *Metadata {
	return &Metadata{
		data: make(chan *pb.Metadata, 1024),
	}
}

// Filter -
func (m *Metadata) Filter(*storage.Metadata) bool {
	return true
}

// Send -
func (m *Metadata) Send(data *pb.Metadata) {
	m.data <- data
}

// Close -
func (m *Metadata) Close() error {
	close(m.data)
	return nil
}

// Listen -
func (m *Metadata) Listen() <-chan *pb.Metadata {
	return m.data
}
