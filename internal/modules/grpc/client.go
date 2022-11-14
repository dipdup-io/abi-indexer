package grpc

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/dipdup-net/abi-indexer/internal/messages"
	"github.com/dipdup-net/abi-indexer/internal/modules/grpc/pb"
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
func NewClient(cfg ClientConfig) *Client {
	return &Client{
		publisher:     messages.NewPublisher(),
		subscriptions: cfg.Subscriptions,
		serverAddress: cfg.ServerAddress,
		wg:            new(sync.WaitGroup),
	}
}

// Start -
func (client *Client) Start(ctx context.Context) {
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
		log.Err(err).Msg("dial connection")
		return
	}
	client.conn = conn
	client.authClient = pb.NewHelloServiceClient(conn)
	client.client = pb.NewMetadataServiceClient(conn)

	hello, err := client.authClient.Hello(ctx, new(pb.HelloRequest))
	if err != nil {
		log.Err(err).Msg("error after hello request")
		return
	}
	client.subscriptionID = hello.Id

	client.wg.Add(1)
	go client.subscribeOnMetadata(ctx, client.subscriptions.Head)

}

func (client *Client) subscribeOnMetadata(ctx context.Context, active bool) {
	defer client.wg.Done()

	if !active {
		return
	}

	stream, err := client.client.SubscribeOnMetadata(ctx, MetadataRequest(client.subscriptionID))
	if err != nil {
		log.Err(err).Msg("subscribe on head")
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
				log.Err(err).Msg("receiving head error")
				continue
			}

			log.Trace().Str("contract", metadata.Address).Msg("new metadata")
			client.publisher.Notify(messages.NewMessage(messages.TopicMetadata, metadata))
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
