package middleware

import (
	"errors"
	"socket/internal/core/ports"
	"socket/pkg/apperror"
	"socket/pkg/util"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	jwtUtils *util.JWTUtils
	userRepo ports.UserRepository
}

func NewAuthMiddleware(jwtUtils *util.JWTUtils, userRepo ports.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtUtils: jwtUtils,
		userRepo: userRepo,
	}
}

func (a *AuthMiddleware) Auth(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")

	if authHeader == "" {
		return apperror.UnauthorizedError(errors.New("request without authorization header"), "Authorization header is required")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return apperror.UnauthorizedError(errors.New("invalid authorization header"), "Authorization header is invalid")
	}

	token := authHeader[7:]
	claims, err := a.jwtUtils.DecodeJWT(token)
	if err != nil {
		return apperror.UnauthorizedError(err, "Invalid token")
	}

	username := claims.Username

	user, err := a.userRepo.GetUserByUsername(username)
	if err != nil {
		return apperror.UnauthorizedError(err, "User not found")
	}

	ctx.Locals("username", username)
	ctx.Locals("user", user)

	return ctx.Next()
}

func (a *AuthMiddleware) Websocket(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return apperror.UpgradeRequiredError(errors.New("your request is not intended for a websocket upgrade"), "your request is not intended for a websocket upgrade")
	}

	token := c.Query("token")
	claims, err := a.jwtUtils.DecodeJWT(token)
	if err != nil {
		return apperror.UnauthorizedError(err, "Invalid token")
	}

	username := claims.Username

	user, err := a.userRepo.GetUserByUsername(username)
	if err != nil {
		return apperror.UnauthorizedError(err, "User not found")
	}

	c.Locals("username", username)
	c.Locals("user", user)

	return c.Next()
}
