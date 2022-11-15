package grpc

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/abi-indexer/internal/storage/postgres"
	"github.com/dipdup-net/abi-indexer/pkg/modules/grpc/pb"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/rs/zerolog/log"
)

// Server -
type Server struct {
	*grpc.Server
	pb.UnimplementedMetadataServiceServer

	storage               postgres.Storage
	metadataSubscriptions *grpc.Subscriptions[*storage.Metadata, *pb.Metadata]

	wg *sync.WaitGroup
}

// NewServer -
func NewServer(cfg *grpc.ServerConfig, pg postgres.Storage) (*Server, error) {
	if cfg == nil {
		return nil, errors.New("configuration structure of gRPC server is nil")
	}

	server, err := grpc.NewServer(cfg)
	if err != nil {
		return nil, err
	}

	return &Server{
		Server:                server,
		storage:               pg,
		metadataSubscriptions: grpc.NewSubscriptions[*storage.Metadata, *pb.Metadata](),
		wg:                    new(sync.WaitGroup),
	}, nil
}

// Start -
func (server *Server) Start(ctx context.Context) {
	server.Server.Start(ctx)

	server.wg.Add(1)
	go server.listen(ctx)
}

func (server *Server) listen(ctx context.Context) {
	defer server.wg.Done()

	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

		case msg := <-server.Listen():
			switch typedMsg := msg.Data().(type) {
			case storage.Metadata:
				server.metadataSubscriptions.NotifyAll(&typedMsg, Metadata(&typedMsg))
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
	subscription := NewMetadataSubscription()
	server.metadataSubscriptions.Add(req.Id, subscription)

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case msg := <-subscription.Listen():
			if err := stream.Send(msg); err != nil {
				if err == io.EOF {
					return nil
				}
				log.Err(err).Msg("sending message error")
			}
		}
	}
}

// UnsubscribeFromMetadata -
func (server *Server) UnsubscribeFromMetadata(ctx context.Context, req *generalPB.DefaultRequest) (*generalPB.Message, error) {
	if err := server.metadataSubscriptions.Remove(req.Id); err != nil {
		return nil, err
	}

	return &generalPB.Message{
		Message: grpc.SuccessMessage,
	}, nil
}

// GetMetadata -
func (server *Server) GetMetadata(ctx context.Context, req *pb.GetMetadataRequest) (*pb.Metadata, error) {
	if req == nil {
		return nil, errors.New("invalid request")
	}

	reqCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	metadata, err := server.storage.Metadata.GetByAddress(reqCtx, req.Address)
	if err != nil {
		return nil, err
	}

	return Metadata(metadata), nil
}

// ListMetadata -
func (server *Server) ListMetadata(ctx context.Context, req *pb.ListMetadataRequest) (*pb.ListMetadataResponse, error) {
	p := newPage(req.GetPage())

	metadata, err := server.storage.Metadata.List(ctx, p.limit, p.offset, p.order)
	if err != nil {
		return nil, err
	}

	return ListMetadataResponse(metadata), nil
}

// GetMetadataByMethodSinature -
func (server *Server) GetMetadataByMethodSinature(ctx context.Context, req *pb.GetMetadataByMethodSinatureRequest) (*pb.ListMetadataResponse, error) {
	p := newPage(req.GetPage())

	metadata, err := server.storage.Metadata.GetByMethod(ctx, req.Signature, p.limit, p.offset, p.order)
	if err != nil {
		return nil, err
	}

	return ListMetadataResponse(metadata), nil
}

// GetMetadataByTopic -
func (server *Server) GetMetadataByTopic(ctx context.Context, req *pb.GetMetadataByTopicRequest) (*pb.ListMetadataResponse, error) {
	p := newPage(req.GetPage())

	metadata, err := server.storage.Metadata.GetByTopic(ctx, req.Topic, p.limit, p.offset, p.order)
	if err != nil {
		return nil, err
	}

	return ListMetadataResponse(metadata), nil
}
