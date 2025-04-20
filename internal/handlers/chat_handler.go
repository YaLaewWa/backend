package handlers

import (
	"socket/internal/core/ports"
	"socket/internal/dto"
	"socket/pkg/apperror"
	"socket/pkg/util"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ChatHandler struct {
	service ports.ChatService
}

func NewChatHandler(service ports.ChatService) ports.ChatHandler {
	return &ChatHandler{service: service}
}

func (h *ChatHandler) AddUserToChat(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (h *ChatHandler) CreateDirectChat(c *fiber.Ctx) error {
	req := new(dto.CreateDirectChatRequest)
	if err := c.BodyParser(req); err != nil {
		return apperror.BadRequestError(err, "Invalid input")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return apperror.UnprocessableEntityError(err, "Validation failed")
	}

	chat, err := h.service.CreateDirectChat(req.User1, req.User2)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(chat.ToDTO()))
}

func (h *ChatHandler) CreateGroupChat(c *fiber.Ctx) error {
	req := new(dto.CreateGroupChatRequest)
	if err := c.BodyParser(req); err != nil {
		return apperror.BadRequestError(err, "Invalid input")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return apperror.UnprocessableEntityError(err, "Validation failed")
	}

	chat, err := h.service.CreateGroupChat(req.Name, req.UserIDs)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(chat.ToDTO()))
}

func (h *ChatHandler) GetChats(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	page, limit := util.PaginationQuery(c)

	chats, totalPages, totalRows, err := h.service.GetChatsByUserID(userID, limit, page)
	if err != nil {
		return err
	}

	res := make([]dto.ChatResponse, len(chats))
	for i, chat := range chats {
		res[i] = chat.ToDTO()
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessPagination(res, page, totalPages, limit, totalRows))
}

func (h *ChatHandler) GetChatMembers(c *fiber.Ctx) error {
	panic("unimplemented")
}
