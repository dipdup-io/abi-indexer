package grpc

import (
	"context"
	"io"
	"sync"

	"github.com/dipdup-net/abi-indexer/pkg/modules/grpc/pb"
	metadataModule "github.com/dipdup-net/abi-indexer/pkg/modules/metadata"
	"github.com/dipdup-net/indexer-sdk/pkg/messages"
	grpcModules "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/rs/zerolog/log"
)

// Client -
type Client struct {
	*grpcModules.AuthClient

	client        pb.MetadataServiceClient
	subscriptions *Subscriptions
	wg            *sync.WaitGroup
}

// NewClient -
func NewClient(cfg *ClientConfig) *Client {
	return &Client{
		AuthClient:    grpcModules.NewAuthClient(cfg.ServerAddress),
		subscriptions: cfg.Subscriptions,
		wg:            new(sync.WaitGroup),
	}
}

// Start -
func (client *Client) Start(ctx context.Context) {
	client.client = pb.NewMetadataServiceClient(client.Connection())

	client.wg.Add(1)
	go client.subscribeOnMetadata(ctx)
}

func (client *Client) subscribeOnMetadata(ctx context.Context) {
	defer client.wg.Done()

	if !client.subscriptions.Metadata {
		return
	}

	stream, err := client.client.SubscribeOnMetadata(ctx, MetadataRequest(client.SubscriptionID))
	if err != nil {
		log.Err(err).Msg("subscribe on metadata")
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			metadata, err := stream.Recv()
			if err == io.EOF {
				continue
			}
			if err != nil {
				log.Err(err).Msg("receiving metadata error")
				continue
			}

			log.Trace().Str("contract", metadata.Address).Msg("new metadata")
			client.Publisher().Notify(messages.NewMessage(metadataModule.TopicMetadata, metadata))
		}
	}
}

// Subscribe -
func (client *Client) Subscribe(s *messages.Subscriber, topic messages.Topic) {
	client.Publisher().Subscribe(s, topic)
}

// Unsubscribe -
func (client *Client) Unsubscribe(s *messages.Subscriber, topic messages.Topic) {
	client.Publisher().Unsubscribe(s, topic)
}

// GetMetadata -
func (client *Client) GetMetadata(ctx context.Context, address string) (*pb.Metadata, error) {
	return client.client.GetMetadata(ctx, &pb.GetMetadataRequest{
		Id:      client.SubscriptionID,
		Address: address,
	})
}

// ListMetadata -
func (client *Client) ListMetadata(ctx context.Context, limit, offset uint64, order generalPB.SortOrder) ([]*pb.Metadata, error) {
	response, err := client.client.ListMetadata(ctx, &pb.ListMetadataRequest{
		Id: client.SubscriptionID,
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
		Id: client.SubscriptionID,
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
		Id: client.SubscriptionID,
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
