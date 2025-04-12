package handlers

import (
	"socket/internal/core/ports"
	"socket/internal/dto"
	"socket/pkg/apperror"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service ports.UserService
}

func NewUserHandler(service ports.UserService) ports.UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	user := new(dto.AuthBody)
	err := c.BodyParser(&user)
	if err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}
	err = h.service.Register(user.UserName, user.Password)
	if err != nil {
		return err
	}
	res := dto.Success(dto.UserResponse{Username: user.UserName})
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	body := new(dto.AuthBody)
	err := c.BodyParser(&body)
	if err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}
	user, jwt, err := h.service.Login(body.UserName, body.Password)
	if err != nil {
		return err
	}
	res := dto.Success(dto.UserWithTokenResponse{User: user.ToDTO(), Token: jwt})
	return c.Status(fiber.StatusOK).JSON(res)
}
