package ports

import (
	"socket/internal/core/domain"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MessageRepository interface {
	Create(msg *domain.Message) error
	GetByChatID(chatID uuid.UUID, limit int, page int) ([]domain.Message, int, int, error)
}

type MessageService interface {
	Create(msg *domain.Message) error
	GetByChatID(chatID uuid.UUID, limit int, page int) ([]domain.Message, int, int, error)
}

type MessageSocketHandler interface {
	InitConnection(c *websocket.Conn)
	GetByChatID(c *fiber.Ctx) error
}
