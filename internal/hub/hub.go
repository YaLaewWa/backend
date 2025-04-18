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
			h.sendOnlineUsers(User.Username)
			h.broadcastUser(User.Username, true)
		case username := <-h.Unregister:
			delete(h.Clients, username)
			h.broadcastUser(username, false)
		case msg := <-h.Broadcast:
			for id := range h.Clients {
				h.Clients[id] <- msg
			}
		}
	}
}

// send list of online users to user with username
func (h *Hub) sendOnlineUsers(username string) {
	onlineUsers := []string{}
	for username := range h.Clients {
		onlineUsers = append(onlineUsers, username)
	}

	msg := map[string]interface{}{
		"type":         "online_users",
		"online_users": onlineUsers,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("error: ", err)
	} else {
		h.Clients[username] <- data
	}
}

// broadcast to all user that a user with username register or unregister from the hub
func (h *Hub) broadcastUser(username string, register bool) {
	var msgType string
	if register {
		msgType = "user_login"
	} else {
		msgType = "user_logout"
	}

	msg := map[string]interface{}{
		"type":     msgType,
		"username": username,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("error: ", err)
	} else {
		for username := range h.Clients {
			h.Clients[username] <- data
		}
	}

}
