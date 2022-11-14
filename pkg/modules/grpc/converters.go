package grpc

import (
	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/abi-indexer/pkg/modules/grpc/pb"
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
func MetadataRequest(id string) *pb.DefaultRequest {
	request := new(pb.DefaultRequest)
	request.Id = id
	return request
}

// ToMetadataModel -
func ToMetadataModel(metadata *pb.Metadata) *storage.Metadata {
	return &storage.Metadata{
		Contract:   metadata.GetAddress(),
		Metadata:   metadata.GetMetadata(),
		JSONSchema: metadata.GetJsonSchema(),
	}
}