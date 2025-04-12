package handlers

import (
	"log"
	"socket/internal/core/ports"
	"socket/internal/database"
	"socket/internal/hub"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type MessageSocketHandler struct {
	hub *hub.Hub
	db  *database.Database //TODO change to some service in the future when out database is not a single slice ;)
}

func NewMessageSocketHandler(hub *hub.Hub, db *database.Database) ports.MessageSocketHandler {
	return &MessageSocketHandler{hub, db}
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

	for index := range h.db.Message { // send old message to the new user
		err := c.WriteMessage(1, h.db.Message[index])
		if err != nil {
			log.Println("read:", err)
			break
		}
	}

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
		h.db.Message = append(h.db.Message, msg)
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
