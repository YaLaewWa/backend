package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"socket/internal/config"
	"socket/internal/database"
	"socket/internal/server"
	"syscall"
)

func main() {
	config := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	db, err := database.New(config.DB)
	if err != nil {
		log.Fatalf("Failed to init  DB err: %v", err)
	}
	s := server.NewServer(config.Server, db)
	s.Start(ctx, stop)
}
