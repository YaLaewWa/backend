package ports

import (
	"socket/internal/core/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ChatRepository interface {
	Create(name string, usernames []string, isGroup bool) (*domain.Chat, error)
	GetPaginatedChatMembers(chatID uuid.UUID, limit int, page int) ([]domain.User, int, int, error)
	GetAllChatMembers(chatID uuid.UUID) ([]domain.User, error)
	GetPaginatedChatsByUsername(username string, limit int, page int) ([]domain.Chat, int, int, error)
	GetAllChatsByUsername(username string) ([]domain.Chat, error)
	AddUserToChat(chatID uuid.UUID, username string) error
	GetByID(chatID uuid.UUID) (*domain.Chat, error)
	IsUserInChat(chatID uuid.UUID, username string) (bool, error)
}

type ChatService interface {
	CreateChat(name string, usernames []string, isGroup bool) (*domain.Chat, error)
	GetChatMembers(chatID uuid.UUID, limit int, page int) ([]domain.User, int, int, error)
	GetChatsByUsername(username string, limit int, page int) ([]domain.Chat, int, int, error)
	AddUserToChat(chatID uuid.UUID, username string) (*domain.Chat, error)
	IsUserInChat(chatID uuid.UUID, username string) (bool, error)
}

type ChatHandler interface {
	CreateDirectChat(c *fiber.Ctx) error
	CreateGroupChat(c *fiber.Ctx) error
	GetChatMembers(c *fiber.Ctx) error
	GetChats(c *fiber.Ctx) error
	JoinChat(c *fiber.Ctx) error
}
