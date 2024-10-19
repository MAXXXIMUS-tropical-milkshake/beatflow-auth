package app

import (
	"context"

	grpcapp "github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/app/gprc"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/config"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/postgres"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/redis"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/service/auth"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/service/user"
	userstore "github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/store/postgres/user"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/store/redis/refreshtoken"
)

type App struct {
	GRPCServer *grpcapp.App
	PG         *postgres.Postgres
	RDB        *redis.Redis
}

func New(ctx context.Context, cfg *config.Config) *App {
	// Init logger
	logger.New(cfg.Log.Level)

	// Postgres connection
	pg, err := postgres.New(ctx, cfg.DB.URL)
	if err != nil {
		logger.Log().Fatal(ctx, "error with connection to database: %s", err.Error())
	}

	// Redis connection
	rdb, err := redis.New(ctx, redis.Config{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err != nil {
		logger.Log().Fatal(ctx, "error with connection to redis: %s", err.Error())
	}

	// Auth config
	authConfig := auth.NewConfig(cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)

	// Store
	userStore := userstore.New(pg)
	refreshTokenStore := refreshtoken.New(rdb)

	// Service
	authService := auth.New(userStore, refreshTokenStore, authConfig)
	userService := user.New(userStore)

	// gRPC server
	gRPCApp := grpcapp.New(ctx, authService, userService, authConfig, cfg)

	return &App{
		GRPCServer: gRPCApp,
		PG:         pg,
		RDB:        rdb,
	}
}
