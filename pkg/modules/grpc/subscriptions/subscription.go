package subscriptions

import (
	"io"

	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/abi-indexer/pkg/modules/grpc/pb"
)

// Types -
type Types interface {
	*storage.Metadata
}

// PBTypes -
type PBTypes interface {
	*pb.Metadata
}

// Subscription -
type Subscription[T Types, P PBTypes] interface {
	Filter(typ T) bool
	Send(msg P)
	Listen() <-chan P
	io.Closer
}

// Subscriptions -
type Subscriptions struct {
	ID       string
	Metadata Subscription[*storage.Metadata, *pb.Metadata]
}

// Close -
func (s *Subscriptions) Close() error {
	if s.Metadata != nil {
		if err := s.Metadata.Close(); err != nil {
			return err
		}
	}

	return nil
}
