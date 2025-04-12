package main

import (
	"context"
	"fmt"
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
		fmt.Println("Error to init db")
	}
	s := server.NewServer(config.Server, db)
	s.Start(ctx, stop)
}
