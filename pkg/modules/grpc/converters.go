package grpc

import (
	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/abi-indexer/pkg/modules/grpc/pb"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
)

// Metadata -
func Metadata(metadata *storage.Metadata) *pb.Metadata {
	return &pb.Metadata{
		Address:    metadata.Contract,
		Metadata:   metadata.Metadata,
		JsonSchema: metadata.JSONSchema,
	}
}

// HeadRequest -
func MetadataRequest(id string) *generalPB.DefaultRequest {
	request := new(generalPB.DefaultRequest)
	request.Id = id
	return request
}

// ListMetadataResponse -
func ListMetadataResponse(metadata []*storage.Metadata) *pb.ListMetadataResponse {
	response := &pb.ListMetadataResponse{
		Metadata: make([]*pb.Metadata, 0),
	}
	for i := range metadata {
		response.Metadata = append(response.Metadata, Metadata(metadata[i]))
	}
	return response
}
