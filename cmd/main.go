package main

import (
	"context"
	"os"
	"os/signal"
	"socket/internal/config"
	"socket/internal/server"
	"syscall"
)

func main() {
	config := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	s := server.NewServer(config.Server)
	s.Start(ctx, stop)
}
