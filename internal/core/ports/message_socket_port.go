package ports

import "github.com/gofiber/contrib/websocket"

type MessageSocketHandler interface {
	InitConnection(c *websocket.Conn)
}
