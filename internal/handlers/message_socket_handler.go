package handlers

import (
	"log"
	"socket/internal/core/domain"
	"socket/internal/core/ports"
	"socket/internal/dto"
	"socket/internal/hub"
	"socket/pkg/util"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type MessageSocketHandler struct {
	hub     *hub.Hub
	service ports.MessageService
}

func NewMessageSocketHandler(hub *hub.Hub, service ports.MessageService) ports.MessageSocketHandler {
	return &MessageSocketHandler{hub: hub, service: service}
}

func (h MessageSocketHandler) InitConnection(c *websocket.Conn) {
	user := c.Locals("user").(domain.User)
	username := user.Username
	hubChannel := make(chan []byte, 256) //buffer up to 256 strings
	closeConnection := make(chan bool)
	payload := &hub.RegisterPayload{
		Channel:  hubChannel,
		Username: username,
	}
	h.hub.Register <- payload //register new connection to hub

	go h.readPump(c, username, closeConnection)
	go h.writePump(c, hubChannel)
	for {
		if <-closeConnection {
			break
		}
	}
}

func (h MessageSocketHandler) readPump(c *websocket.Conn, username string, close chan bool) {

	defer func() { //Porperly close connection
		h.hub.Unregister <- username
		close <- true
		c.Close()
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		message := new(domain.Message)
		message.Content = string(msg[:])
		message.Username = username
		if err = h.service.Create(message); err != nil {
			log.Println("error: ", err)
		} else {
			h.hub.Broadcast <- msg
		}
	}
}

func (h MessageSocketHandler) writePump(c *websocket.Conn, channel chan []byte) {
	defer c.Close() // ensure that client can be unsingned from one routine only hopefully ;)
	for {
		msg := <-channel
		err := c.WriteMessage(1, msg)
		if err != nil {
			break
		}
	}
}

// GetAll godoc
// @Summary Get all messages
// @Description Retrieve a list of all messages.
// @Tags Message
// @Security Bearer
// @Produce json
// @Param limit query int false "Number of messages to retrieve (default 10, max 50)"
// @Param page query int false "Page number to retrieve (default 1)"
// @Success 200 {object} dto.PaginationResponse[dto.MessageResponse] "Messages retrieved successfully"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve messages"
// @Router /messages [get]
func (h MessageSocketHandler) GetAll(c *fiber.Ctx) error {
	page, limit := util.PaginationQuery(c)

	msgs, totalPages, totalRows, err := h.service.GetAll(limit, page)
	if err != nil {
		return err
	}

	res := make([]dto.MessageResponse, len(msgs))
	for i, msg := range msgs {
		res[i] = msg.ToDTO()
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessPagination(res, page, totalPages, limit, totalRows))
}
