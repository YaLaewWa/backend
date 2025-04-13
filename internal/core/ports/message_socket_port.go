package ports

import (
	"socket/internal/core/domain"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type MessageRepository interface {
	Create(msg *domain.Message) error
	GetAll(limit, page int) ([]domain.Message, int, int, error)
}

type MessageService interface {
	Create(msg *domain.Message) error
	GetAll(limit, page int) ([]domain.Message, int, int, error)
}

type MessageSocketHandler interface {
	InitConnection(c *websocket.Conn)
	GetAll(c *fiber.Ctx) error
}
