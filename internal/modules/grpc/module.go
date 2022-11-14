package grpc

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/dipdup-net/abi-indexer/internal/messages"
	"github.com/dipdup-net/abi-indexer/internal/modules/grpc/pb"
	"github.com/dipdup-net/abi-indexer/internal/modules/grpc/subscriptions"
	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/abi-indexer/internal/storage/postgres"
	"github.com/rs/zerolog/log"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// Server -
type Server struct {
	*messages.Subscriber
	pb.UnimplementedMetadataServiceServer

	bind string

	storage     postgres.Storage
	subscribers map[string]*subscriptions.Subscriptions
	subsMx      sync.RWMutex

	wg *sync.WaitGroup
}

// NewServer -
func NewServer(cfg ServerConfig, pg postgres.Storage) (*Server, error) {
	subscriber, err := messages.NewSubscriber()
	if err != nil {
		return nil, err
	}
	return &Server{
		bind:        cfg.Bind,
		storage:     pg,
		Subscriber:  subscriber,
		subscribers: make(map[string]*subscriptions.Subscriptions),
		// TODO: map subscribers by topics
		wg: new(sync.WaitGroup),
	}, nil
}

// Start -
func (module *Server) Start(ctx context.Context) {
	module.wg.Add(1)
	go module.grpc(ctx)

	module.wg.Add(1)
	go module.listen(ctx)
}

func (module *Server) grpc(ctx context.Context) {
	defer module.wg.Done()

	log.Info().Str("bind", module.bind).Msg("running grpc...")

	listener, err := net.Listen("tcp", module.bind)
	if err != nil {
		log.Err(err).Msg("net.Listen")
		return
	}
	grpcServer := gogrpc.NewServer(
		gogrpc.StatsHandler(module),
		gogrpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    20 * time.Second,
				Timeout: 10 * time.Second,
			},
		),
		gogrpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             10 * time.Second,
				PermitWithoutStream: true,
			},
		))
	pb.RegisterMetadataServiceServer(grpcServer, module)
	if err := grpcServer.Serve(listener); err != nil {
		log.Err(err).Msg("grpcServer.Serve")
	}
}

func (module *Server) listen(ctx context.Context) {
	defer module.wg.Done()

	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

		case msg := <-module.Listen():
			switch typedMsg := msg.Data().(type) {
			case *storage.Metadata:
				module.metadataHandler(typedMsg)
			default:
				log.Warn().Msgf("unknown message type: %T", typedMsg)
			}
		}
	}
}

func (module *Server) metadataHandler(metadata *storage.Metadata) {
	module.subsMx.RLock()
	{
		for _, subscriber := range module.subscribers {
			if subscriber.Metadata != nil && subscriber.Metadata.Filter(metadata) {
				subscriber.Metadata.Send(Metadata(metadata))
			}
		}
	}
	module.subsMx.RUnlock()
}

// Close -
func (module *Server) Close() error {
	if err := module.Subscriber.Close(); err != nil {
		return err
	}
	return nil
}
