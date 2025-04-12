package hub

import "github.com/google/uuid"

type RegisterPayload struct {
	Channel chan []byte
	ID      uuid.UUID
}

type Hub struct {
	Clients    map[uuid.UUID]chan []byte
	Register   chan *RegisterPayload
	Unregister chan uuid.UUID
	Broadcast  chan []byte
	Message    [][]byte //TODO delete in the future when schema for message is ready
}

func NewHub() *Hub {
	var message [][]byte
	return &Hub{
		Clients:    make(map[uuid.UUID]chan []byte),
		Register:   make(chan *RegisterPayload),
		Unregister: make(chan uuid.UUID),
		Broadcast:  make(chan []byte),
		Message:    message,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case UserID := <-h.Register:
			h.Clients[UserID.ID] = UserID.Channel
			for index := range h.Message {
				h.Clients[UserID.ID] <- h.Message[index]
			}
		case ID := <-h.Unregister:
			delete(h.Clients, ID)
		case msg := <-h.Broadcast:
			h.Message = append(h.Message, msg)
			for id := range h.Clients {
				h.Clients[id] <- msg
			}
		}
	}
}
