package grpc

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/abi-indexer/pkg/modules/grpc/pb"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/rs/zerolog/log"
)

// Server -
type Server struct {
	*grpc.Server
	pb.UnimplementedMetadataServiceServer

	metadata              storage.IMetadata
	metadataSubscriptions *grpc.Subscriptions[*storage.Metadata, *pb.SubscriptionMetadata]

	wg *sync.WaitGroup
}

// NewServer -
func NewServer(cfg *grpc.ServerConfig, metadata storage.IMetadata) (*Server, error) {
	if cfg == nil {
		return nil, errors.New("configuration structure of gRPC server is nil")
	}

	server, err := grpc.NewServer(cfg)
	if err != nil {
		return nil, err
	}

	return &Server{
		Server:                server,
		metadata:              metadata,
		metadataSubscriptions: grpc.NewSubscriptions[*storage.Metadata, *pb.SubscriptionMetadata](),
		wg:                    new(sync.WaitGroup),
	}, nil
}

// Start -
func (server *Server) Start(ctx context.Context) {
	pb.RegisterMetadataServiceServer(server.Server.Server(), server)

	server.Server.Start(ctx)

	server.wg.Add(1)
	go server.listen(ctx)
}

func (server *Server) listen(ctx context.Context) {
	defer server.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-server.Listen():
			switch typedMsg := msg.Data().(type) {
			case storage.Metadata:
				server.metadataSubscriptions.NotifyAll(&typedMsg, SubscriptionMetadata)
			default:
				log.Warn().Msgf("unknown message type: %T", typedMsg)
			}
		}
	}
}

// Close -
func (server *Server) Close() error {
	return server.Server.Close()
}

////////////////////////////////////////////////
//////////////    HANDLERS    //////////////////
////////////////////////////////////////////////

// SubscribeOnMetadata -
func (server *Server) SubscribeOnMetadata(req *generalPB.DefaultRequest, stream pb.MetadataService_SubscribeOnMetadataServer) error {
	return grpc.DefaultSubscribeOn[*storage.Metadata, *pb.SubscriptionMetadata](
		stream,
		server.metadataSubscriptions,
		NewMetadataSubscription(),
	)
}

// UnsubscribeFromMetadata -
func (server *Server) UnsubscribeFromMetadata(ctx context.Context, req *generalPB.UnsubscribeRequest) (*generalPB.UnsubscribeResponse, error) {
	return grpc.DefaultUnsubscribe(ctx, server.metadataSubscriptions, req.Id)
}

// GetMetadata -
func (server *Server) GetMetadata(ctx context.Context, req *pb.GetMetadataRequest) (*pb.Metadata, error) {
	reqCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	metadata, err := server.metadata.GetByAddress(reqCtx, req.Address)
	if err != nil {
		return nil, err
	}

	return Metadata(metadata), nil
}

// ListMetadata -
func (server *Server) ListMetadata(ctx context.Context, req *pb.ListMetadataRequest) (*pb.ListMetadataResponse, error) {
	p := newPage(req.GetPage())

	metadata, err := server.metadata.List(ctx, p.limit, p.offset, p.order)
	if err != nil {
		return nil, err
	}

	return ListMetadataResponse(metadata), nil
}

// GetMetadataByMethodSinature -
func (server *Server) GetMetadataByMethodSinature(ctx context.Context, req *pb.GetMetadataByMethodSinatureRequest) (*pb.ListMetadataResponse, error) {
	p := newPage(req.GetPage())

	metadata, err := server.metadata.GetByMethod(ctx, req.Signature, p.limit, p.offset, p.order)
	if err != nil {
		return nil, err
	}

	return ListMetadataResponse(metadata), nil
}

// GetMetadataByTopic -
func (server *Server) GetMetadataByTopic(ctx context.Context, req *pb.GetMetadataByTopicRequest) (*pb.ListMetadataResponse, error) {
	p := newPage(req.GetPage())

	metadata, err := server.metadata.GetByTopic(ctx, req.Topic, p.limit, p.offset, p.order)
	if err != nil {
		return nil, err
	}

	return ListMetadataResponse(metadata), nil
}
