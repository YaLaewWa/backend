package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"socket/internal/config"
	"socket/internal/database"
	"socket/internal/server"
	"socket/pkg/util"
	"syscall"
)

// @title chadChat API
// @version 1.0
// @description This is a swagger to show all RestAPI of chadChat project
// @contact.name API Support
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @description Bearer token authentication
func main() {
	config := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	db, err := database.New(config.DB)
	if err != nil {
		log.Fatalf("Failed to init  DB err: %v", err)
	}
	jwt := util.NewJWTUtils(&config.Jwt)
	s := server.NewServer(config.Server, db, jwt)
	s.Start(ctx, stop)
}
