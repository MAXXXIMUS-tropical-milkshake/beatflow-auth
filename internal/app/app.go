package app

import (
	"context"

	grpcapp "github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/app/gprc"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/config"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/postgres"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/service/auth"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/store/postgres/user"
)

type App struct {
	GRPCServer *grpcapp.App
	PG         *postgres.Postgres
}

func New(ctx context.Context, cfg *config.Config) *App {
	// Init logger
	logger.New(cfg.Log.Level)

	// Postgres connection
	pg, err := postgres.New(ctx, cfg.DB.URL)
	if err != nil {
		logger.Log().Fatal(ctx, "error with connection to database: %s", err.Error())
	}

	// Auth config
	authConfig := auth.NewConfig(cfg.JWTSecret, cfg.TokenTTL)

	// Store
	userStore := user.New(pg)

	// Service
	authService := auth.New(userStore, authConfig)

	// gRPC server
	gRPCApp := grpcapp.New(ctx, authService, cfg)

	return &App{
		GRPCServer: gRPCApp,
		PG:         pg,
	}
}
