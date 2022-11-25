package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/dipdup-net/abi-indexer/internal/storage/postgres"
	"github.com/dipdup-net/abi-indexer/pkg/modules/grpc"
	"github.com/dipdup-net/abi-indexer/pkg/modules/metadata"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"

	"github.com/dipdup-net/go-lib/config"
)

var (
	rootCmd = &cobra.Command{
		Use:   "indexer",
		Short: "DipDup indexer",
	}
)

func main() {
	configPath := rootCmd.PersistentFlags().StringP("config", "c", "dipdup.yml", "path to YAML config file")
	if err := rootCmd.Execute(); err != nil {
		log.Panic().Err(err).Msg("command line execute")
		return
	}
	if err := rootCmd.MarkFlagRequired("config"); err != nil {
		log.Panic().Err(err).Msg("config command line arg is required")
		return
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	})

	var cfg Config
	if err := config.Parse(*configPath, &cfg); err != nil {
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

	if err := modules.Connect(metadataIndexer, grpcModule, metadata.OutputMetadata, metadata.OutputMetadata); err != nil {
		log.Panic().Err(err).Msg("connecting modules error")
		cancel()
		return
	}

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
