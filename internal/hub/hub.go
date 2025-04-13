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
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uuid.UUID]chan []byte),
		Register:   make(chan *RegisterPayload),
		Unregister: make(chan uuid.UUID),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case UserID := <-h.Register:
			h.Clients[UserID.ID] = UserID.Channel
			// for index := range h.Message {
			// 	h.Clients[UserID.ID] <- h.Message[index]
			// }
		case ID := <-h.Unregister:
			delete(h.Clients, ID)
		case msg := <-h.Broadcast:
			for id := range h.Clients {
				h.Clients[id] <- msg
			}
		}
	}
}
