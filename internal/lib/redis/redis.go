package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
	"github.com/redis/go-redis/v9"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
	_defaultMaxPoolSize  = 10
)

type Redis struct {
	connAttempts int
	connTimeout  time.Duration
	maxPoolSize  int
	*redis.Client
}

type Config struct {
	Addr, Password string
	DB             int
}

func New(ctx context.Context, config Config, opts ...Option) (*Redis, error) {
	rdb := &Redis{
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
		maxPoolSize:  _defaultMaxPoolSize,
	}

	// Custom options
	for _, opt := range opts {
		opt(rdb)
	}

	var db *redis.Client
	var err error

	for rdb.connAttempts > 0 {
		db = redis.NewClient(&redis.Options{
			Addr:     config.Addr,
			Password: config.Password,
			DB:       config.DB,
			PoolSize: rdb.maxPoolSize,
		})

		_, err = db.Ping(ctx).Result()
		if err != nil {
			logger.Log().Debug(ctx,
				"redis is trying to connect, attempts left: %d", rdb.connAttempts,
			)
		} else {
			rdb.Client = db
			break
		}

		time.Sleep(rdb.connTimeout)

		rdb.connAttempts--
	}

	if err != nil {
		logger.Log().Fatal(ctx, "failed to connect to redis")
		return nil, err
	}

	return rdb, nil
}

func (r *Redis) Close() error {
	if err := r.Client.Close(); err != nil {
		return fmt.Errorf("error closing redis client: %w", err)
	}

	return nil
}
