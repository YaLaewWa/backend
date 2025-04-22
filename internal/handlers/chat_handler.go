package handlers

import (
	"errors"
	"slices"
	"socket/internal/core/ports"
	"socket/internal/dto"
	"socket/pkg/apperror"
	"socket/pkg/util"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	service      ports.ChatService
	queueService ports.MessageQueueService
}

func NewChatHandler(service ports.ChatService, queueService ports.MessageQueueService) ports.ChatHandler {
	return &ChatHandler{service: service, queueService: queueService}
}

func (h *ChatHandler) JoinChat(c *fiber.Ctx) error {
	username := c.Locals("username").(string)
	id := c.Params("id")
	chatID, err := util.ParseIdParam(id)
	if err != nil {
		return err
	}

	chat, err := h.service.AddUserToChat(chatID, username)
	if err != nil {
		return err
	}

	err = h.queueService.Create(username, chatID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(chat.ToDTO()))
}

func (h *ChatHandler) CreateDirectChat(c *fiber.Ctx) error {
	return h.createChat(c, false)
}

func (h *ChatHandler) CreateGroupChat(c *fiber.Ctx) error {
	return h.createChat(c, true)
}

func (h *ChatHandler) createChat(c *fiber.Ctx, isGroup bool) error {
	req := new(dto.CreateChatRequest)
	if err := c.BodyParser(req); err != nil {
		return apperror.BadRequestError(err, "Invalid input")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return apperror.UnprocessableEntityError(err, "Validation failed")
	}

	// Check if user creating chat including themselves or not
	username := c.Locals("username").(string)
	if !slices.Contains(req.Usernames, username) {
		return apperror.UnprocessableEntityError(errors.New("validation failed"), "You can not create chat without you in it")
	}

	// Not sure if will have user amount check in group chat yet or not
	if !isGroup && len(req.Usernames) != 2 {
		return apperror.UnprocessableEntityError(errors.New("validation failed"), "You can not create direct message chat with less or more than 2 users")
	}

	chat, err := h.service.CreateChat(req.Name, req.Usernames, isGroup)
	if err != nil {
		return err
	}
	err = h.queueService.Create(username, chat.ID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(chat.ToDTO()))
}

func (h *ChatHandler) GetChats(c *fiber.Ctx) error {
	username := c.Locals("username").(string)
	page, limit := util.PaginationQuery(c)

	chats, totalPages, totalRows, err := h.service.GetChatsByUsername(username, limit, page)
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
	id := c.Params("id")
	chatID, err := util.ParseIdParam(id)
	if err != nil {
		return err
	}
	page, limit := util.PaginationQuery(c)

	members, totalPages, totalRows, err := h.service.GetChatMembers(chatID, limit, page)
	if err != nil {
		return err
	}

	res := make([]dto.UserResponse, len(members))
	for i, member := range members {
		res[i] = *member.ToDTO()
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessPagination(res, page, totalPages, limit, totalRows))
}

func (h *ChatHandler) GetGroupChats(c *fiber.Ctx) error {
	username := c.Locals("username").(string)
	page, limit := util.PaginationQuery(c)

	groups, totalPages, totalRows, err := h.service.GetGroupChats(username, limit, page)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessPagination(groups, page, totalPages, limit, totalRows))
}
