package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(ctx context.Context, cfg Config) (Wrapper, error) {
	cfg.setDefaultValues()

	var (
		pgxCfg *pgxpool.Config
		err    error
	)
	if pgxCfg, err = pgxpool.ParseConfig(cfg.dsn()); err != nil {
		return nil, fmt.Errorf("error parsing connString: %v", err)
	}

	pgxCfg.MaxConns = cfg.MaxConns
	pgxCfg.MaxConnLifetime = time.Duration(cfg.MaxConnLifetimeSec) * time.Second
	pgxCfg.MaxConnIdleTime = time.Duration(cfg.MaxConnIdleTimeSec) * time.Second

	var db *pgxpool.Pool
	if db, err = pgxpool.NewWithConfig(ctx, pgxCfg); err != nil {
		return nil, fmt.Errorf("error open connection: %w", err)
	}

	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error ping connection: %w", err)
	}

	return &wrapper{
		pgxPool: db,
	}, nil
}

type (
	Wrapper interface {
		Begin(ctx context.Context) (Tx, error)
		Close()
		Ping(ctx context.Context) error
		Exec(ctx context.Context, sql string, arguments ...interface{}) (CommandTag, error)
		Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)
		QueryRow(ctx context.Context, sql string, args ...interface{}) Row
		getPgxPool() *pgxpool.Pool
	}

	wrapper struct {
		pgxPool *pgxpool.Pool
	}

	SqlDB = sql.DB
)

func (w *wrapper) Begin(ctx context.Context) (Tx, error) {
	return w.pgxPool.Begin(ctx)
}

func (w *wrapper) Close() {
	w.pgxPool.Close()
}

func (w *wrapper) Ping(ctx context.Context) error {
	return w.pgxPool.Ping(ctx)
}

func (w *wrapper) Exec(ctx context.Context, sql string, args ...interface{}) (CommandTag, error) {
	return w.pgxPool.Exec(ctx, sql, args...)
}

func (w *wrapper) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	return w.pgxPool.Query(ctx, sql, args...)
}

func (w *wrapper) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
	return w.pgxPool.QueryRow(ctx, sql, args...)
}

func (w *wrapper) getPgxPool() *pgxpool.Pool {
	return w.pgxPool
}

func NewSqlDB(pgWrapper Wrapper) (*SqlDB, error) {
	var (
		db  *SqlDB
		err error

		connString = pgWrapper.getPgxPool().Config().ConnString()
	)
	if db, err = sql.Open("pgx", connString); err != nil {
		return nil, fmt.Errorf("error open connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error ping connection: %v", err)
	}

	return db, nil
}
