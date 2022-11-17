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

// SubscriptionMetadata -
func SubscriptionMetadata(id uint64, metadata *storage.Metadata) *pb.SubscriptionMetadata {
	return &pb.SubscriptionMetadata{
		Subscription: &generalPB.SubscribeResponse{
			Id: id,
		},
		Metadata: Metadata(metadata),
	}
}

// HeadRequest -
func MetadataRequest() *generalPB.DefaultRequest {
	return new(generalPB.DefaultRequest)
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
