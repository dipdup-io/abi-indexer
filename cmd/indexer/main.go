package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/dipdup-net/abi-indexer/internal/storage/postgres"
	"github.com/dipdup-net/abi-indexer/pkg/modules/grpc"
	"github.com/dipdup-net/abi-indexer/pkg/modules/metadata"

	"github.com/dipdup-net/go-lib/cmdline"
	"github.com/dipdup-net/go-lib/config"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	})

	args := cmdline.Parse()
	if args.Help {
		return
	}

	var cfg Config
	if err := config.Parse(args.Config, &cfg); err != nil {
		log.Panic().Err(err).Msg("parsing config file")
		return
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = zerolog.LevelInfoValue
	}

	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Panic().Err(err).Msg("parsing log level")
		return
	}
	zerolog.SetGlobalLevel(logLevel)

	ctx, cancel := context.WithCancel(context.Background())

	storage, err := postgres.Create(ctx, cfg.Database)
	if err != nil {
		log.Panic().Err(err).Msg("postgres connection error")
		cancel()
		return
	}
	metadataIndexer, err := metadata.NewMetadata(cfg.Metadata, storage.Metadata, storage.Events, storage.Methods, storage.Transactable)
	if err != nil {
		log.Panic().Err(err).Msg("creating indexer")
		cancel()
		return
	}

	grpcModule, err := grpc.NewServer(cfg.GRPC.Server, storage.Metadata)
	if err != nil {
		log.Panic().Err(err).Msg("creating grpc module")
		cancel()
		return
	}

	metadataIndexer.Subscribe(grpcModule.Subscriber)

	metadataIndexer.Start(ctx)
	grpcModule.Start(ctx)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-signals
	cancel()

	if err := metadataIndexer.Close(); err != nil {
		log.Panic().Err(err).Msg("closing metadata indexer")
	}
	if err := grpcModule.Close(); err != nil {
		log.Panic().Err(err).Msg("closing grpc server")
	}
	if err := storage.Close(); err != nil {
		log.Panic().Err(err).Msg("closing postgres connection")
	}

	close(signals)
}
