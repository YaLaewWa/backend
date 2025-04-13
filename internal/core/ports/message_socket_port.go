package ports

import (
	"socket/internal/core/domain"

	"github.com/gofiber/contrib/websocket"
)

type MessageRepository interface {
	Create(msg *domain.Message) error
	GetAll() ([]domain.Message, error)
}

type MessageService interface {
	Create(msg *domain.Message) error
	GetAll() ([]domain.Message, error)
}

type MessageSocketHandler interface {
	InitConnection(c *websocket.Conn)
}
