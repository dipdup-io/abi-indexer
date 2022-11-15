package grpc

import (
	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/abi-indexer/pkg/modules/grpc/pb"
)

// MetadataSubscription -
type MetadataSubscription struct {
	data chan *pb.Metadata
}

// NewMetadataSubscription -
func NewMetadataSubscription() *MetadataSubscription {
	return &MetadataSubscription{
		data: make(chan *pb.Metadata, 1024),
	}
}

// Filter -
func (m *MetadataSubscription) Filter(*storage.Metadata) bool {
	return true
}

// Send -
func (m *MetadataSubscription) Send(data *pb.Metadata) {
	m.data <- data
}

// Close -
func (m *MetadataSubscription) Close() error {
	close(m.data)
	return nil
}

// Listen -
func (m *MetadataSubscription) Listen() <-chan *pb.Metadata {
	return m.data
}
