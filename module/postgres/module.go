package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-brick-template/go-brick-template/internal/application/runtime"
	"github.com/go-brick-template/go-brick-template/internal/config/options"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrModuleIsNotInitialized = errors.New("postgres module is not initialized")

type DB struct {
	master *pgxpool.Pool
	slave  *pgxpool.Pool
}

func (db *DB) Write(ctx context.Context) *pgxpool.Pool { return db.master }
func (db *DB) Read(ctx context.Context) *pgxpool.Pool  { return db.slave }

type Module struct {
	db *DB
}

func New(ctx context.Context, opts *options.PostgresOptions) (*Module, error) {
	if opts == nil {
		return nil, ErrModuleIsNotInitialized
	}
	db, err := openPools(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Module{db: db}, nil
}

func (m *Module) Init(_ context.Context, _ runtime.Info) error {
	return nil
}

func (m *Module) Run(_ context.Context, c *runtime.Closer) {
	c.Add("postgres", func() error {
		if m.db == nil {
			return nil
		}
		m.db.master.Close()
		if m.db.slave != m.db.master {
			m.db.slave.Close()
		}
		return nil
	})
}

func (m *Module) DB() *DB { return m.db }

func openPools(ctx context.Context, opts *options.PostgresOptions) (*DB, error) {
	poolCfg, err := pgxpool.ParseConfig(opts.URL)
	if err != nil {
		return nil, fmt.Errorf("parse postgres url: %w", err)
	}
	poolCfg.MaxConns = opts.MaxConns
	poolCfg.MinConns = opts.MinConns
	poolCfg.MaxConnLifetime = opts.MaxConnLifetime

	master, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	if err := master.Ping(ctx); err != nil {
		master.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return &DB{master: master, slave: master}, nil
}

func Ping(ctx context.Context, db *DB) error {
	if db == nil {
		return ErrModuleIsNotInitialized
	}
	return db.master.Ping(ctx)
}
