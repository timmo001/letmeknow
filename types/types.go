package types

import "github.com/gorilla/websocket"

type Client struct {
	Connection *websocket.Conn
	UserID     *string
}

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
