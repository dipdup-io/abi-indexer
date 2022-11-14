package grpc

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/dipdup-net/abi-indexer/pkg/modules/grpc/pb"
	metadataModule "github.com/dipdup-net/abi-indexer/pkg/modules/metadata"
	"github.com/dipdup-net/indexer-sdk/messages"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Client -
type Client struct {
	publisher *messages.Publisher
	conn      *gogrpc.ClientConn

	client     pb.MetadataServiceClient
	authClient pb.HelloServiceClient

	subscriptions  *Subscriptions
	subscriptionID string
	serverAddress  string

	wg *sync.WaitGroup
}

// NewClient -
func NewClient(cfg *ClientConfig) *Client {
	return &Client{
		publisher:     messages.NewPublisher(),
		subscriptions: cfg.Subscriptions,
		serverAddress: cfg.ServerAddress,
		wg:            new(sync.WaitGroup),
	}
}

// Connect -
func (client *Client) Connect(ctx context.Context) error {
	conn, err := gogrpc.Dial(
		client.serverAddress,
		gogrpc.WithTransportCredentials(insecure.NewCredentials()),
		gogrpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                (time.Duration(20) * time.Second),
			Timeout:             (time.Duration(10) * time.Second),
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return errors.Wrap(err, "dial connection")
	}
	client.conn = conn
	client.authClient = pb.NewHelloServiceClient(conn)
	client.client = pb.NewMetadataServiceClient(conn)

	hello, err := client.authClient.Hello(ctx, new(pb.HelloRequest))
	if err != nil {
		return errors.Wrap(err, "error after hello request")
	}
	client.subscriptionID = hello.Id
	return nil
}

// Start -
func (client *Client) Start(ctx context.Context) {
	client.wg.Add(1)
	go client.subscribeOnMetadata(ctx)
}

func (client *Client) subscribeOnMetadata(ctx context.Context) {
	defer client.wg.Done()

	if !client.subscriptions.Metadata {
		return
	}

	stream, err := client.client.SubscribeOnMetadata(ctx, MetadataRequest(client.subscriptionID))
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
			client.publisher.Notify(messages.NewMessage(metadataModule.TopicMetadata, metadata))
		}
	}
}

// Close -
func (client *Client) Close() error {
	if err := client.conn.Close(); err != nil {
		return err
	}
	return nil
}

// Subscribe -
func (client *Client) Subscribe(s *messages.Subscriber, topic messages.Topic) {
	client.publisher.Subscribe(s, topic)
}

// Unsubscribe -
func (client *Client) Unsubscribe(s *messages.Subscriber, topic messages.Topic) {
	client.publisher.Unsubscribe(s, topic)
}

// GetMetadata -
func (client *Client) GetMetadata(ctx context.Context, address string) (*pb.Metadata, error) {
	return client.client.GetMetadata(ctx, &pb.GetMetadataRequest{
		Id:      client.subscriptionID,
		Address: address,
	})
}

// ListMetadata -
func (client *Client) ListMetadata(ctx context.Context, limit, offset uint64, order pb.SortOrder) ([]*pb.Metadata, error) {
	response, err := client.client.ListMetadata(ctx, &pb.ListMetadataRequest{
		Id: client.subscriptionID,
		Page: &pb.Page{
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
func (client *Client) GetMetadataByMethodSinature(ctx context.Context, limit, offset uint64, order pb.SortOrder, signature string) ([]*pb.Metadata, error) {
	response, err := client.client.GetMetadataByMethodSinature(ctx, &pb.GetMetadataByMethodSinatureRequest{
		Id: client.subscriptionID,
		Page: &pb.Page{
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
func (client *Client) GetMetadataByTopic(ctx context.Context, limit, offset uint64, order pb.SortOrder, topic string) ([]*pb.Metadata, error) {
	response, err := client.client.GetMetadataByTopic(ctx, &pb.GetMetadataByTopicRequest{
		Id: client.subscriptionID,
		Page: &pb.Page{
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
