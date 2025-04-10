package database

import (
	"github.com/gofiber/contrib/websocket"
)

type Database struct {
	Clients map[*websocket.Conn]bool
	Message []string
}

func NewDatabase() *Database { //Just a mock dont kill me
	clients := make(map[*websocket.Conn]bool)
	var message []string
	db := &Database{clients, message}
	return db
}
