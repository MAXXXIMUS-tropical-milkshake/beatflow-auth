package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	_defaultMaxPoolSize  = 10
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	DB *sql.DB
}

func New(ctx context.Context, dbURL string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	var db *sql.DB
	var err error

	for pg.connAttempts > 0 {
		db, err = sql.Open("pgx", dbURL)
		if err == nil && db.Ping() == nil {
			db.SetMaxOpenConns(pg.maxPoolSize)
			db.SetConnMaxLifetime(time.Hour)

			pg.DB = db
			break
		}

		logger.Log().Debug(ctx,
			"postgres is trying to connect, attempts left: %d", pg.connAttempts,
		)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		logger.Log().Fatal(ctx, "failed to connect to database: %s", err.Error())
		return nil, err
	}

	return pg, nil
}

func (p *Postgres) Close(ctx context.Context) {
	if err := p.DB.Close(); err != nil {
		logger.Log().Info(ctx, "Error closing database connection: %s", err.Error())
	}
}
