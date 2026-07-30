package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"github.com/maksgudimov/catalog-service/internal/app/config/section"
	"github.com/maksgudimov/catalog-service/migration"
)

type (
	Client struct {
		_bunDB
		rawBunDB *bun.DB
		cfg      section.RepositoryPostgres
	}

	_bunDB = bun.IDB
)

func (c *Client) GetRawBunDB() *bun.DB {
	return c.rawBunDB
}

func (c *Client) Migrate(ctx context.Context) (oldVer, newVer int64, err error) {
	log.Printf("Starting database migrations...")

	migrations := migrate.NewMigrations()
	if err = migrations.Discover(migration.Postgres); err != nil {
		return 0, 0, fmt.Errorf("failed to discover migrations: %w", err)
	}

	opts := []migrate.MigratorOption{
		migrate.WithTableName(c.cfg.MigrationTable),
		migrate.WithLocksTableName(c.cfg.MigrationTable + "_lock"),
		migrate.WithMarkAppliedOnSuccess(true),
	}
	migrator := migrate.NewMigrator(c.GetRawBunDB(), migrations, opts...)

	if err = migrator.Init(ctx); err != nil {
		return 0, 0, fmt.Errorf("failed to init migrations table: %w", err)
	}
	log.Printf("Migrations table initialized: %s", c.cfg.MigrationTable)

	applied, err := migrator.AppliedMigrations(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get applied migrations: %w", err)
	}
	oldVer, err = maxMigrationVersion(applied)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to determine current migration version: %w", err)
	}

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return oldVer, 0, fmt.Errorf("failed migrate: %w", err)
	}

	newVer = oldVer
	if group != nil {
		newVerMaxValue, err := maxMigrationVersion(group.Migrations)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to determine current migration version: %w", err)
		}
		if newVerMaxValue > newVer {
			newVer = newVerMaxValue
			log.Printf("Database migrations completed: version %d -> %d", oldVer, newVer)
		} else {
			log.Printf("No new migrations")
		}

	}

	return oldVer, newVer, nil
}

func NewClient(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	var u url.URL
	u.Scheme = "postgres"
	u.Host = cfg.Address
	u.User = url.UserPassword(cfg.Username, cfg.Password)
	u.Path = cfg.Name

	args := make(url.Values)
	args.Set("sslmode", "disable")
	u.RawQuery = args.Encode()

	dsn := u.String()

	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithReadTimeout(cfg.ReadTimeout),
		pgdriver.WithWriteTimeout(cfg.WriteTimeout),
	))
	if sqlDB == nil {
		return nil, fmt.Errorf("failed to create sql.DB: driver returned nil")
	}

	sqlDB.SetMaxOpenConns(10)

	db := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())
	if db == nil {
		return nil, fmt.Errorf("failed to create bun.DB: bun returned nil")
	}

	log.Printf("Pinging database with timeout 2s...")

	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctxTimeout); err != nil {
		return nil, fmt.Errorf("failed to connect to the database at %s: %w", cfg.Address, err)
	}

	log.Printf("Database ping successful")

	return &Client{
		_bunDB:   db,
		rawBunDB: db,
		cfg:      cfg,
	}, nil
}

func maxMigrationVersion(ms migrate.MigrationSlice) (int64, error) {
	var maxValue int64

	for _, mg := range ms {
		v, err := strconv.ParseInt(mg.Name, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid migration name %q: %w", mg.Name, err)
		}

		if v > maxValue {
			maxValue = v
		}
	}

	return maxValue, nil
}
