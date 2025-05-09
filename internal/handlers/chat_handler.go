package handlers

import (
	"errors"
	"slices"
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"socket/internal/dto"
	"socket/internal/hub"
	"socket/pkg/apperror"
	"socket/pkg/util"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	service      ports.ChatService
	queueService ports.MessageQueueService
	hub          *hub.Hub
}

func NewChatHandler(service ports.ChatService, queueService ports.MessageQueueService, hub *hub.Hub) ports.ChatHandler {
	return &ChatHandler{service: service, queueService: queueService, hub: hub}
}

func (h *ChatHandler) broadcastSideBar(username string) error {
	payload := make(map[string]any)
	queues, err := h.queueService.Get(username)
	if err != nil {
		return err
	}
	var queueDTO []dto.QueueResponse
	for _, o := range queues {
		queueDTO = append(queueDTO, o.ToDTO())
	}
	payload["queue"] = queueDTO
	h.hub.Broadcast <- domain.HubMessage{Type: "sidebar_update", Payload: payload, To: []domain.User{{Username: username}}}
	return nil
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

	payload := make(map[string]any)
	payload["chatID"] = chat.ID
	payload["joiner"] = username
	h.hub.Broadcast <- domain.HubMessage{Type: "new_user_group", Payload: payload}
	err = h.queueService.Create(username, chatID)
	if err != nil {
		return err
	}
	err = h.broadcastSideBar(username)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(dto.Success(chat.ToDTO()))
}

func (h *ChatHandler) CreateDirectChat(c *fiber.Ctx) error {
	chat, err := h.createChat(c, false)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(dto.Success(chat.ToDTO()))
}

func (h *ChatHandler) CreateGroupChat(c *fiber.Ctx) error {
	chat, err := h.createChat(c, true)
	if err != nil {
		return err
	}

	username := c.Locals("username").(string)
	payload := make(map[string]any)
	payload["chat"] = chat
	payload["creator"] = username
	h.hub.Broadcast <- domain.HubMessage{Type: "new_group", Payload: payload}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(chat.ToDTO()))
}

func (h *ChatHandler) createChat(c *fiber.Ctx, isGroup bool) (*domain.Chat, error) {
	req := new(dto.CreateChatRequest)
	if err := c.BodyParser(req); err != nil {
		return nil, apperror.BadRequestError(err, "Invalid input")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, apperror.UnprocessableEntityError(err, "Validation failed")
	}

	// Check if user creating chat including themselves or not
	username := c.Locals("username").(string)
	if !slices.Contains(req.Usernames, username) {
		return nil, apperror.UnprocessableEntityError(errors.New("validation failed"), "You can not create chat without you in it")
	}

	if !isGroup {
		if len(req.Usernames) != 2 {
			return nil, apperror.UnprocessableEntityError(errors.New("validation failed"), "You can not create direct message chat with less or more than 2 users")
		}
		if req.Usernames[0] == req.Usernames[1] {
			return nil, apperror.UnprocessableEntityError(errors.New("validation failed"), "You can not create direct message chat with yourself")
		}
	}

	chat, err := h.service.CreateChat(req.Name, req.Usernames, isGroup)
	if err != nil {
		return nil, err
	}

	for _, name := range req.Usernames {
		err = h.queueService.Create(name, chat.ID)
		if err != nil {
			return nil, err
		}
		err = h.broadcastSideBar(name)
		if err != nil {
			return nil, err
		}
	}

	return chat, nil
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
		if !res[i].IsGroup {
			if res[i].Members[0].Username == username {
				res[i].Name = res[i].Members[1].Username
			} else {
				res[i].Name = res[i].Members[0].Username
			}
		}
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

func (h *ChatHandler) HavePrivateChat(c *fiber.Ctx) error {
	username1 := c.Params("username1")
	username2 := c.Params("username2")
	status, id, err := h.service.HavePrivateChat(username1, username2)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(map[string]any{"status": status, "id": id}))
}
