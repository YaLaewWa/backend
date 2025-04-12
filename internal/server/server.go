package server

import (
	"context"
	"fmt"
	"log"
	"socket/internal/database"

	"socket/pkg/apperror"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Config struct {
	Name string `env:"NAME"`
	Port int    `env:"PORT"`
	Env  string `env:"ENV"`
}

type Server struct {
	app        *fiber.App
	config     Config
	db         *database.Database
	pgDB       *gorm.DB
	repository *Repository
	service    *Service
	handler    *Handler
}

func NewServer(config Config, pgDB *gorm.DB) *Server {

	app := fiber.New(fiber.Config{
		AppName:               config.Name,
		CaseSensitive:         true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
		ErrorHandler:          apperror.ErrorHandler,
	})

	db := database.NewDatabase()

	return &Server{
		app:    app,
		config: config,
		db:     db,
		pgDB:   pgDB,
	}
}

func (s *Server) Start(ctx context.Context, stop context.CancelFunc) {
	//init services
	s.initRoutes()
	s.initRepository()
	s.initService()
	s.initHandler()

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
