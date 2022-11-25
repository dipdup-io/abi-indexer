package grpc

import (
	"context"
	"sync"

	"github.com/dipdup-net/abi-indexer/pkg/modules/grpc/pb"
	"github.com/dipdup-net/abi-indexer/pkg/modules/metadata"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	grpcModules "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Client -
type Client struct {
	*grpcModules.Client

	output        *modules.Output
	client        pb.MetadataServiceClient
	subscriptions *Subscriptions
	wg            *sync.WaitGroup
}

// NewClient -
func NewClient(cfg *ClientConfig) *Client {
	return &Client{
		Client:        grpcModules.NewClient(cfg.ServerAddress),
		output:        modules.NewOutput(metadata.OutputMetadata),
		subscriptions: cfg.Subscriptions,
		wg:            new(sync.WaitGroup),
	}
}

// Start -
func (client *Client) Start(ctx context.Context) {
	client.client = pb.NewMetadataServiceClient(client.Connection())
}

// Name -
func (client *Client) Name() string {
	return "metadata_grpc_client"
}

// Input -
func (client *Client) Input(name string) (*modules.Input, error) {
	return nil, errors.Wrap(modules.ErrUnknownInput, name)
}

// Output -
func (client *Client) Output(name string) (*modules.Output, error) {
	if name != metadata.OutputMetadata {
		return nil, errors.Wrap(modules.ErrUnknownOutput, name)
	}
	return client.output, nil
}

// AttachTo -
func (client *Client) AttachTo(name string, input *modules.Input) error {
	output, err := client.Output(name)
	if err != nil {
		return err
	}
	output.Attach(input)
	return nil
}

// SubscribeOnMetadata -
func (client *Client) SubscribeOnMetadata(ctx context.Context) (uint64, error) {
	if client.subscriptions != nil && !client.subscriptions.Metadata {
		return 0, nil
	}

	stream, err := client.client.SubscribeOnMetadata(ctx, MetadataRequest())
	if err != nil {
		return 0, err
	}

	return grpcModules.Subscribe[*pb.SubscriptionMetadata](
		stream,
		client.handleNewMetadata,
		client.wg,
	)
}

func (client *Client) handleNewMetadata(ctx context.Context, data *pb.SubscriptionMetadata, id uint64) error {
	log.Trace().Str("contract", data.Metadata.Address).Msg("new metadata")
	client.output.Push(data)
	return nil
}

// UnsubscribeFromMetadata -
func (client *Client) UnsubscribeFromMetadata(ctx context.Context, id uint64) error {
	if client.subscriptions != nil && !client.subscriptions.Metadata {
		return nil
	}
	if _, err := client.client.UnsubscribeFromMetadata(ctx, &generalPB.UnsubscribeRequest{
		Id: id,
	}); err != nil {
		return err
	}

	return nil
}

// GetMetadata -
func (client *Client) GetMetadata(ctx context.Context, address string) (*pb.Metadata, error) {
	return client.client.GetMetadata(ctx, &pb.GetMetadataRequest{
		Address: address,
	})
}

// ListMetadata -
func (client *Client) ListMetadata(ctx context.Context, limit, offset uint64, order generalPB.SortOrder) ([]*pb.Metadata, error) {
	response, err := client.client.ListMetadata(ctx, &pb.ListMetadataRequest{
		Page: &generalPB.Page{
			Limit:  limit,
			Offset: offset,
			Order:  order,
		},
	})
	if err != nil {
		return nil, err
	}
	return response.Metadata, nil
}

// GetMetadataByMethodSinature -
func (client *Client) GetMetadataByMethodSinature(ctx context.Context, limit, offset uint64, order generalPB.SortOrder, signature string) ([]*pb.Metadata, error) {
	response, err := client.client.GetMetadataByMethodSinature(ctx, &pb.GetMetadataByMethodSinatureRequest{
		Page: &generalPB.Page{
			Limit:  limit,
			Offset: offset,
			Order:  order,
		},
		Signature: signature,
	})
	if err != nil {
		return nil, err
	}
	return response.Metadata, nil
}

// GetMetadataByTopic -
func (client *Client) GetMetadataByTopic(ctx context.Context, limit, offset uint64, order generalPB.SortOrder, topic string) ([]*pb.Metadata, error) {
	response, err := client.client.GetMetadataByTopic(ctx, &pb.GetMetadataByTopicRequest{
		Page: &generalPB.Page{
			Limit:  limit,
			Offset: offset,
			Order:  order,
		},
		Topic: topic,
	})
	if err != nil {
		return nil, err
	}
	return response.Metadata, nil
}
