package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/app"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	ctx := context.Background()

	application := app.New(ctx, cfg)

	// Closing DBs
	defer application.PG.Close(ctx)
	defer application.RDB.Close()

	go func() { application.GRPCServer.MustRun(ctx) }()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	// Stopping server
	application.GRPCServer.Stop(ctx)
}
