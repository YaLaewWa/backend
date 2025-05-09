package hub

import (
	"encoding/json"
	"log"
	"socket/internal/core/domain"
	"socket/internal/dto"
	"sync"

	"github.com/google/uuid"
)

type RegisterPayload struct {
	Channel  chan []byte
	Username string
}

type Hub struct {
	Clients     map[string]chan []byte
	Register    chan *RegisterPayload
	Unregister  chan string
	Broadcast   chan domain.HubMessage
	ClientMutex sync.Mutex
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
			h.ClientMutex.Lock()
			h.Clients[User.Username] = User.Channel
			h.ClientMutex.Unlock()
			h.sendOnlineUsers(User.Username)
			h.broadcastUser(User.Username)
		case username := <-h.Unregister:
			h.ClientMutex.Lock()
			delete(h.Clients, username)
			h.ClientMutex.Unlock()
			h.broadcastUser(username)
		case msg := <-h.Broadcast:
			if msg.Type == "message" {
				h.ClientMutex.Lock()
				for _, member := range msg.To {
					if _, ok := h.Clients[member.Username]; ok {
						message := msg.Payload.(domain.Message)
						data, err := json.Marshal(message.ToDTO())
						if err != nil {
							log.Println("error: ", err)
						} else {
							h.Clients[member.Username] <- data
						}
					}
				}
				h.ClientMutex.Unlock()
			} else if msg.Type == "new_group" {
				h.ClientMutex.Lock()
				for username := range h.Clients {
					payload := msg.Payload.(map[string]any)
					group := payload["chat"].(*domain.Chat)
					creator := payload["creator"].(string)
					data, err := json.Marshal(group.ToSocketDTO(username == creator))
					if err != nil {
						log.Println("error:", err)
					} else {
						h.Clients[username] <- data
					}
				}
				h.ClientMutex.Unlock()
			} else if msg.Type == "new_user_group" {
				h.ClientMutex.Lock()
				for username := range h.Clients {
					payload := msg.Payload.(map[string]any)
					chatID := payload["chatID"].(uuid.UUID)
					joiner := payload["joiner"].(string)
					data, err := json.Marshal(dto.GetJoinSocketDTO(chatID, joiner))
					if err != nil {
						log.Println("error:", err)
					} else {
						h.Clients[username] <- data
					}
				}
				h.ClientMutex.Unlock()
			} else if msg.Type == "sidebar_update" {
				h.ClientMutex.Lock()
				data, err := json.Marshal(dto.QueueSocket{Type: msg.Type, Payload: msg.Payload})
				if err != nil {
					log.Println("error:", err)
				} else {
					h.Clients[msg.To[0].Username] <- data
				}
				h.ClientMutex.Unlock()
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

	msg := map[string]any{
		"type":    "online_users",
		"payload": onlineUsers,
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

	msg := map[string]any{
		"type":    msgType,
		"payload": username,
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
