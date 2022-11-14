package grpc

import (
	"github.com/dipdup-net/abi-indexer/internal/modules/grpc/pb"
	"github.com/dipdup-net/abi-indexer/internal/storage"
)

// Metadata -
func Metadata(metadata *storage.Metadata) *pb.Metadata {
	return &pb.Metadata{}
}

// HeadRequest -
func MetadataRequest(id string) *pb.DefaultRequest {
	request := new(pb.DefaultRequest)
	request.Id = id
	return request
}
