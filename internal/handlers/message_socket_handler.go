package handlers

import (
	"log"
	"socket/internal/core/ports"
	"socket/internal/hub"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type MessageSocketHandler struct {
	hub *hub.Hub
}

func NewMessageSocketHandler(hub *hub.Hub) ports.MessageSocketHandler {
	return &MessageSocketHandler{hub}
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
