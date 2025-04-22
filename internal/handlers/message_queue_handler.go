package handlers

import (
	"socket/internal/core/ports"
	"socket/internal/dto"

	"github.com/gofiber/fiber/v2"
)

type MessageQueueHandler struct {
	service ports.MessageQueueService
}

func NewMessageQueueHandler(service ports.MessageQueueService) ports.MessageQueueHandler {
	return &MessageQueueHandler{service: service}
}

func (h *MessageQueueHandler) Get(c *fiber.Ctx) error {
	username := c.Locals("username").(string)
	queue, err := h.service.Get(username)
	if err != nil {
		return err
	}
	res := make([]dto.QueueResponse, len(queue))
	for i, q := range queue {
		res[i] = q.ToDTO()
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(res))
}
