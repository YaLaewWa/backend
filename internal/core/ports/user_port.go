package ports

import (
	"socket/internal/core/domain"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
)

type UserHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type UserService interface {
	Register(userName, password string) error
	Login(userName, password string) (*domain.User, string, error)
}

type UserRepository interface {
	Create(user *domain.User) error
	GetUserByID(userID uuid.UUID) (*domain.User, error)
	GetUserByUsername(userName string) (*domain.User, error)
}
