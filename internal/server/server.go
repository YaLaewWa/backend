package server

import (
	"context"
	"fmt"
	"log"
	"socket/internal/database"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Name string `env:"NAME"`
	Port int    `env:"PORT"`
	Env  string `env:"ENV"`
}

type Server struct {
	app    *fiber.App
	config Config
	db     *database.Database
}

func NewServer(config Config) *Server {

	app := fiber.New(fiber.Config{
		AppName:               config.Name,
		CaseSensitive:         true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
	})

	db := database.NewDatabase()

	return &Server{
		app:    app,
		config: config,
		db:     db,
	}
}

func (s *Server) Start(ctx context.Context, stop context.CancelFunc) {
	//init services
	s.initRoutes()

	// start server
	go func() {
		if err := s.app.Listen(fmt.Sprintf(":%d", s.config.Port)); err != nil {
			log.Panicf("Failed to start server: %v\n", err)
			stop()
		}
	}()

	// shutdown server at the end
	defer func() {
		if err := s.app.ShutdownWithContext(ctx); err != nil {
			log.Printf("Failed to shutdown server: %v\n", err)
		}
		log.Println("Server stopped")
	}()

	<-ctx.Done()

	log.Println("Server is shutting down...")
}
