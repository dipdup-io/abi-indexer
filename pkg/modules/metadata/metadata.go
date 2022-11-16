package metadata

import (
	"context"

	"github.com/dipdup-net/abi-indexer/internal/sources"
	models "github.com/dipdup-net/abi-indexer/internal/storage"

	"github.com/dipdup-net/abi-indexer/internal/vm"
	"github.com/dipdup-net/indexer-sdk/pkg/messages"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/workerpool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// topics
const (
	TopicMetadata messages.Topic = "metadata"
)

// Metadata -
type Metadata struct {
	publisher *messages.Publisher

	repo         models.IMetadata
	events       models.IEvent
	methods      models.IMethod
	transactable storage.Transactable

	source sources.Source
	vmType vm.Type

	pool *workerpool.TimedPool[string]
}

// NewMetadata -
func NewMetadata(
	cfg Config,
	metadataRepo models.IMetadata,
	events models.IEvent,
	methods models.IMethod,
	transactable storage.Transactable,
) (*Metadata, error) {
	src, err := sources.Factory(cfg.SourceType, sources.FactoryParams{
		Sourcify: cfg.Sourcify,
		FS:       cfg.FS,
	})
	if err != nil {
		return nil, err
	}

	if cfg.ThreadsCount <= 0 {
		cfg.ThreadsCount = 10
	}

	metadata := &Metadata{
		repo:         metadataRepo,
		methods:      methods,
		events:       events,
		transactable: transactable,
		source:       src,
		vmType:       cfg.VM.Type,
		publisher:    messages.NewPublisher(),
	}

	metadata.pool = workerpool.NewTimedPool(
		metadata.dispatcher,
		metadata.worker,
		metadata.errorHandler,
		cfg.ThreadsCount,
		3600000)

	return metadata, nil
}

// Start -
func (metadata *Metadata) Start(ctx context.Context) {
	metadata.pool.Start(ctx)
}

func (metadata *Metadata) errorHandler(ctx context.Context, err error) {
	log.Err(err).Msg("error during fetching metadata")
}

func (metadata *Metadata) dispatcher(ctx context.Context) ([]string, error) {
	return metadata.source.List(ctx)
}

func (metadata *Metadata) worker(ctx context.Context, task string) {
	if err := metadata.processData(ctx, task); err != nil {
		log.Err(err).Msg("processing metadata error")
	}
}

func (metadata *Metadata) processData(ctx context.Context, address string) error {
	_, err := metadata.repo.GetByAddress(ctx, address)
	switch {
	case err == nil:
		return nil
	case !metadata.repo.IsNoRows(err):
		return err
	}

	log.Info().Str("address", address).Msg("new metadata was found")

	data, err := metadata.source.Get(ctx, address)
	if err != nil {
		return errors.Wrap(err, address)
	}

	model := models.Metadata{
		Contract: address,
		Metadata: data,
	}

	machine, err := vm.Factory(metadata.vmType, model.Metadata)
	if err != nil {
		return err
	}

	schema, err := machine.JSONSchema()
	if err != nil {
		return err
	}

	model.JSONSchema = schema

	methods, err := machine.Methods()
	if err != nil {
		return err
	}

	events, err := machine.Events()
	if err != nil {
		return err
	}

	if err := metadata.save(ctx, model, methods, events); err != nil {
		return err
	}

	metadata.publisher.Notify(messages.NewMessage(TopicMetadata, model))

	return nil
}

func (metadata *Metadata) save(ctx context.Context, model models.Metadata, methods []models.Method, events []models.Event) error {
	tx, err := metadata.transactable.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Close(ctx); err != nil {
			log.Err(err).Msg("closing postgres transaction error")
		}
	}()

	if err := tx.Add(ctx, &model); err != nil {
		return tx.HandleError(ctx, err)
	}

	if len(methods) > 0 {
		data := make([]any, len(methods))
		for i := range methods {
			methods[i].MetadataID = model.ID
			data[i] = &methods[i]
		}

		if err := tx.BulkSave(ctx, data); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if len(events) > 0 {
		data := make([]any, len(events))
		for i := range events {
			events[i].MetadataID = model.ID
			data[i] = &events[i]
		}

		if err := tx.BulkSave(ctx, data); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	return tx.Flush(ctx)
}

// Close -
func (metadata *Metadata) Close() error {
	log.Info().Msg("closing metadata indexer...")

	if err := metadata.pool.Close(); err != nil {
		return err
	}

	return nil
}

// Subscribe -
func (metadata *Metadata) Subscribe(s *messages.Subscriber, topic messages.Topic) {
	metadata.publisher.Subscribe(s, topic)
}

// Unsubscribe -
func (metadata *Metadata) Unsubscribe(s *messages.Subscriber, topic messages.Topic) {
	metadata.publisher.Unsubscribe(s, topic)
}
