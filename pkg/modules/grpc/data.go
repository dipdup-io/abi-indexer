package grpc

import (
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

type page struct {
	limit  uint64
	offset uint64
	order  storage.SortOrder
}

func newPage(req *pb.Page) *page {
	p := new(page)
	if req != nil {
		p.limit = req.Limit
		p.offset = req.Offset

		switch req.Order {
		case pb.SortOrder_ASC:
			p.order = storage.SortOrderAsc
		case pb.SortOrder_DESC:
			p.order = storage.SortOrderDesc
		default:
			p.order = storage.SortOrderAsc
		}
	}
	return p
}
