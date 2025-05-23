package handlers

import (
	"encoding/json"
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
	message ports.MessageService
	chat    ports.ChatService
	queue   ports.MessageQueueService
}

func NewMessageSocketHandler(hub *hub.Hub, message ports.MessageService, chat ports.ChatService, queue ports.MessageQueueService) ports.MessageSocketHandler {
	return &MessageSocketHandler{hub: hub, message: message, chat: chat, queue: queue}
}

func (h *MessageSocketHandler) broadcastSideBar(username string) error {
	payload := make(map[string]any)
	queues, err := h.queue.Get(username)
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

func (h MessageSocketHandler) InitConnection(c *websocket.Conn) {
	user := c.Locals("user").(*domain.User)
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

		var input dto.MessageRequest
		err = json.Unmarshal(msg, &input)
		if err != nil {
			log.Println("error: ", err)
			continue
		}

		if input.Type == "message" {
			payload := input.Payload
			isMember, err := h.chat.IsUserInChat(payload.ChatID, username)
			if err != nil {
				log.Println("error: ", err)
				continue
			} else if !isMember {
				log.Println(username, "tried to send message to chat they aren't in.")
				continue
			}

			savedMsg, err := h.message.Create(username, payload.ChatID, payload.Content)
			if err != nil {
				log.Println("error: ", err)
				continue
			}

			members, _, _, err := h.chat.GetChatMembers(payload.ChatID, -1, 0)
			if err != nil {
				log.Println("error: ", err)
			} else {
				hubMsg := domain.HubMessage{Type: "message", Payload: *savedMsg, To: members}
				h.hub.ClientMutex.Lock()
				for _, member := range hubMsg.To {
					if _, ok := h.hub.Clients[member.Username]; !ok {
						err = h.queue.ReceiveMessage(member.Username, payload.ChatID)
						if err != nil {
							log.Println(err)
							continue
						}
					}
				}
				h.hub.ClientMutex.Unlock()
				h.hub.Broadcast <- hubMsg
			}
		} else if input.Type == "read_chat" {
			payload := input.Payload
			err = h.queue.ReadMessage(c.Locals("username").(string), payload.ChatID)
			if err != nil {
				log.Println(err)
				continue
			}
			err = h.broadcastSideBar(username)
			if err != nil {
				log.Println(err)
				continue
			}
		} else if input.Type == "ignore" {
			payload := input.Payload
			err = h.queue.ReceiveMessage(c.Locals("username").(string), payload.ChatID)
			if err != nil {
				log.Println(err)
				continue
			}
			err = h.broadcastSideBar(username)
			if err != nil {
				log.Println(err)
				continue
			}
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
func (h MessageSocketHandler) GetByChatID(c *fiber.Ctx) error {
	id := c.Params("id")
	chatID, err := util.ParseIdParam(id)
	if err != nil {
		return err
	}

	page, limit := util.PaginationQuery(c)

	username := c.Locals("username").(string)

	msgs, totalPages, totalRows, err := h.message.GetByChatID(chatID, limit, page, username)
	if err != nil {
		return err
	}

	res := make([]dto.MessageResponse, len(msgs))
	for i, msg := range msgs {
		res[i] = msg.ToDTO()
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessPagination(res, page, totalPages, limit, totalRows))
}
