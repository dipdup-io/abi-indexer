package postgres

import (
	"context"
	"strings"

	models "github.com/dipdup-net/abi-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// Storage -
type Storage struct {
	*postgres.Storage

	Metadata models.IMetadata
	Methods  models.IMethod
	Events   models.IEvent
}

// Create -
func Create(ctx context.Context, cfg config.Database) (*Storage, error) {
	strg, err := postgres.Create(ctx, cfg, initDatabase)
	if err != nil {
		return nil, nil
	}

	return &Storage{
		Storage:  strg,
		Metadata: NewMetadata(strg.Connection()),
		Events:   NewEvents(strg.Connection()),
		Methods:  NewMethods(strg.Connection()),
	}, nil
}

func initDatabase(ctx context.Context, conn *database.PgGo) error {
	if _, err := conn.DB().ExecContext(ctx, "create role posgrest_anon nologin"); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			if err := conn.Close(); err != nil {
				return err
			}
			return err
		}
	}

	if _, err := conn.DB().ExecContext(ctx, "grant usage on schema public to posgrest_anon;"); err != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		return err
	}

	for _, data := range []storage.Model{
		&models.Metadata{}, &models.Method{}, &models.Event{},
	} {
		if err := conn.DB().WithContext(ctx).Model(data).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		}); err != nil {
			if err := conn.Close(); err != nil {
				return err
			}
			return err
		}

		if _, err := conn.DB().
			WithParam("SCHEMA", pg.Ident("public")).
			WithParam("NAME", pg.Ident(data.TableName())).
			ModelContext(ctx, data).
			Exec("grant select on ?SCHEMA.?NAME to posgrest_anon;"); err != nil {
			if err := conn.Close(); err != nil {
				return err
			}
			return err
		}
	}
	return createIndices(ctx, conn)
}

func createIndices(ctx context.Context, conn *database.PgGo) error {
	return conn.DB().RunInTransaction(ctx, func(tx *pg.Tx) error {
		// Methods
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS methods_metadata_id ON methods (metadata_id)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS methods_signature_id ON methods (signature_id)`); err != nil {
			return err
		}

		// Events
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS events_metadata_id ON events (metadata_id)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS events_signature_id ON events (signature_id)`); err != nil {
			return err
		}

		return nil
	})
}
