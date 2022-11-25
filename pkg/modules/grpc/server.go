package grpc

import (
	"context"
	"sync"
	"time"

	"github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/abi-indexer/pkg/modules/grpc/pb"
	"github.com/dipdup-net/abi-indexer/pkg/modules/metadata"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/pkg/errors"
)

// Server -
type Server struct {
	*grpc.Server
	pb.UnimplementedMetadataServiceServer

	input *modules.Input

	metadata              storage.IMetadata
	metadataSubscriptions *grpc.Subscriptions[*storage.Metadata, *pb.SubscriptionMetadata]

	wg *sync.WaitGroup
}

// NewServer -
func NewServer(
	cfg *grpc.ServerConfig,
	metadataRepo storage.IMetadata,
) (*Server, error) {
	if cfg == nil {
		return nil, errors.New("configuration structure of gRPC server is nil")
	}

	server, err := grpc.NewServer(cfg)
	if err != nil {
		return nil, err
	}

	return &Server{
		Server:                server,
		input:                 modules.NewInput(metadata.OutputMetadata),
		metadataSubscriptions: grpc.NewSubscriptions[*storage.Metadata, *pb.SubscriptionMetadata](),
		metadata:              metadataRepo,
		wg:                    new(sync.WaitGroup),
	}, nil
}

// Name -
func (server *Server) Name() string {
	return "metadata_grpc_server"
}

// Input -
func (server *Server) Input(name string) (*modules.Input, error) {
	if name != metadata.OutputMetadata {
		return nil, errors.Wrap(modules.ErrUnknownInput, name)
	}
	return server.input, nil
}

// Output -
func (server *Server) Output(name string) (*modules.Output, error) {
	return nil, errors.Wrap(modules.ErrUnknownOutput, name)
}

// AttachTo -
func (server *Server) AttachTo(name string, input *modules.Input) error {
	return errors.Wrap(modules.ErrUnknownOutput, name)
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
		case msg, ok := <-server.input.Listen():
			if !ok {
				return
			}
			model, ok := msg.(storage.Metadata)
			if !ok {
				continue
			}
			server.metadataSubscriptions.NotifyAll(&model, SubscriptionMetadata)
		}
	}
}

// Close -
func (server *Server) Close() error {
	if err := server.input.Close(); err != nil {
		return err
	}
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
