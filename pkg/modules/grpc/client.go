package grpc

import (
	"context"
	"errors"
	"sync"

	"github.com/dipdup-net/abi-indexer/pkg/modules/grpc/pb"
	"github.com/dipdup-net/indexer-sdk/pkg/messages"
	grpcModules "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/rs/zerolog/log"
)

// Client -
type Client struct {
	*grpcModules.Client

	client        pb.MetadataServiceClient
	subscriptions *Subscriptions
	wg            *sync.WaitGroup
}

// NewClient -
func NewClient(cfg *ClientConfig) *Client {
	return &Client{
		Client:        grpcModules.NewClient(cfg.ServerAddress),
		subscriptions: cfg.Subscriptions,
		wg:            new(sync.WaitGroup),
	}
}

// Start -
func (client *Client) Start(ctx context.Context) {
	client.client = pb.NewMetadataServiceClient(client.Connection())
}

// SubscribeOnMetadata -
func (client *Client) SubscribeOnMetadata(ctx context.Context, s *messages.Subscriber) (messages.SubscriptionID, error) {
	if client.subscriptions != nil && !client.subscriptions.Metadata {
		return 0, nil
	}

	stream, err := client.client.SubscribeOnMetadata(ctx, MetadataRequest())
	if err != nil {
		return 0, err
	}

	return grpcModules.Subscribe[*pb.SubscriptionMetadata](
		client.Publisher(),
		s,
		stream,
		client.handleNewMetadata,
		client.wg,
	)
}

func (client *Client) handleNewMetadata(ctx context.Context, data *pb.SubscriptionMetadata, id messages.SubscriptionID) error {
	log.Trace().Str("contract", data.Metadata.Address).Msg("new metadata")
	client.Publisher().Notify(messages.NewMessage(id, data))
	return nil
}

// UnsubscribeFromMetadata -
func (client *Client) UnsubscribeFromMetadata(ctx context.Context, s *messages.Subscriber, id messages.SubscriptionID) error {
	if client.subscriptions != nil && !client.subscriptions.Metadata {
		return nil
	}

	subscriptionID, ok := id.(uint64)
	if !ok {
		return errors.New("invalid subscription id")
	}

	if _, err := client.client.UnsubscribeFromMetadata(ctx, &generalPB.UnsubscribeRequest{
		Id: subscriptionID,
	}); err != nil {
		return err
	}

	client.Publisher().Unsubscribe(s, id)
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
