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

// Register godoc
// @Summary Register a user
// @Description Register a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequestBody true "Description of the authentication body"
// @Success 201 {object} dto.SuccessResponse[dto.UserResponse] "register successfully"
// @Failure 400 {object} dto.ErrorResponse "your request is invalid or your request body is incorrect or cannot save user"
// @Failure 409 {object} dto.ErrorResponse "username already exists"
// @Failure 500 {object} dto.ErrorResponse "cannot use this password"
// @Router /auth/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	user := new(dto.RegisterRequestBody)
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

// Login godoc
// @Summary Login to the system
// @Description Login to the system
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.LoginRequestBody true "Description of the authentication body"
// @Success 200 {object} dto.SuccessResponse[dto.UserWithTokenResponse] "successfully login"
// @Failure 400 {object} dto.ErrorResponse "your request is invalid or your request body is incorrect or cannot save user"
// @Failure 500 {object} dto.ErrorResponse "cannot use this password"
// @Router /auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	body := new(dto.LoginRequestBody)
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
