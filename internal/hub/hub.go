package hub

import (
	"encoding/json"
	"log"
)

type RegisterPayload struct {
	Channel  chan []byte
	Username string
}

type Hub struct {
	Clients    map[string]chan []byte
	Register   chan *RegisterPayload
	Unregister chan string
	Broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]chan []byte),
		Register:   make(chan *RegisterPayload, 256),
		Unregister: make(chan string, 256),
		Broadcast:  make(chan []byte, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case User := <-h.Register:
			h.Clients[User.Username] = User.Channel
			h.broadcastOnlineUsers()
		case ID := <-h.Unregister:
			delete(h.Clients, ID)
			h.broadcastOnlineUsers()
		case msg := <-h.Broadcast:
			for id := range h.Clients {
				h.Clients[id] <- msg
			}
		}
	}
}

func (h *Hub) broadcastOnlineUsers() {
	onlineUsers := make([]string, len(h.Clients))
	for username := range h.Clients {
		onlineUsers = append(onlineUsers, username)
	}

	// not too sure about this yet but
	// our json may have to look like {"type":"?", "content":?}
	// because we can have many kind of stuff to send over socket
	// not too sure what to call each field yet
	msg := map[string]interface{}{
		"type":         "online_users",
		"online_users": onlineUsers,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("error: ", err)
	} else {
		for id := range h.Clients {
			h.Clients[id] <- data
		}
	}

}
