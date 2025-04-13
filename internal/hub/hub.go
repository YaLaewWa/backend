package hub

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
		Register:   make(chan *RegisterPayload),
		Unregister: make(chan string),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case User := <-h.Register:
			h.Clients[User.Username] = User.Channel
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
