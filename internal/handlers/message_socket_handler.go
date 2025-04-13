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
	"github.com/google/uuid"
)

type MessageSocketHandler struct {
	hub     *hub.Hub
	service ports.MessageService
}

func NewMessageSocketHandler(hub *hub.Hub, service ports.MessageService) ports.MessageSocketHandler {
	return &MessageSocketHandler{hub: hub, service: service}
}

func (h MessageSocketHandler) InitConnection(c *websocket.Conn) {
	id := uuid.New()                     //TODO: change to user id in the future
	hubChannel := make(chan []byte, 256) //buffer up to 256 strings
	closeConnection := make(chan bool)
	payload := &hub.RegisterPayload{
		Channel: hubChannel,
		ID:      id,
	}
	h.hub.Register <- payload //register new connection to hub

	go h.readPump(c, id, closeConnection)
	go h.writePump(c, hubChannel)
	for {
		if <-closeConnection {
			break
		}
	}
}

func (h MessageSocketHandler) readPump(c *websocket.Conn, id uuid.UUID, close chan bool) {

	defer func() { //Porperly close connection
		h.hub.Unregister <- id
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
		if err = h.service.Create(message); err != nil {
			// not sure how to actually handle this case
			log.Println("error: ", err)
		}

		h.hub.Broadcast <- msg
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
