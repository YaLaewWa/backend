package hub

import (
	"encoding/json"
	"log"
	"socket/internal/core/domain"
	"sync"
)

type RegisterPayload struct {
	Channel  chan []byte
	Username string
}

type Hub struct {
	Clients    map[string]chan []byte
	Register   chan *RegisterPayload
	Unregister chan string
	Broadcast  chan domain.HubMessage
	Mutex      *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]chan []byte),
		Register:   make(chan *RegisterPayload, 256),
		Unregister: make(chan string, 256),
		Broadcast:  make(chan domain.HubMessage, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case User := <-h.Register:
			h.Clients[User.Username] = User.Channel
			h.sendOnlineUsers(User.Username)
			h.broadcastUser(User.Username)
		case username := <-h.Unregister:
			delete(h.Clients, username)
			h.broadcastUser(username)
		case msg := <-h.Broadcast:
			for _, member := range msg.To {
				if _, ok := h.Clients[member.Username]; ok {
					data, err := json.Marshal(msg.Message.ToDTO())
					if err != nil {
						log.Println("errror: ", err)
					} else {
						h.Clients[member.Username] <- data
					}
				}
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
func (h *Hub) broadcastUser(username string) {
	var msgType string
	_, ok := h.Clients[username]
	if ok {
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
		for u := range h.Clients {
			if u != username {
				h.Clients[u] <- data
			}
		}
	}

}
