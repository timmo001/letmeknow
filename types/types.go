package types

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	Connection *websocket.Conn
	UserID     *string
}

func (c *Client) Display() string {
	if c.UserID == nil {
		return fmt.Sprintf("Client(Connection=%s, UserID=nil)", c.Connection.RemoteAddr())
	}
	return fmt.Sprintf("Client(Connection=%s, UserID=%s)", c.Connection.RemoteAddr(), *c.UserID)
}

// TODO: Make sure registration contains password/token
type ClientRegistration struct {
	UserID string `json:"userID"`
}

type RequestMessage struct {
	Message string   `json:"message"`
	Targets []string `json:"targets"`
}

type ResponseError struct {
	Message string  `json:"message"`
	Error   *string `json:"error"`
}

type ResponseSuccess struct {
	Succeeded bool   `json:"succeeded"`
	Message   string `json:"message"`
}
