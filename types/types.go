package types

import "github.com/gorilla/websocket"

type Client struct {
	ID         string `json:"id"`
	Connection *websocket.Conn
}

type RequestMessage struct {
	Message string   `json:"message"`
	Targets []string `json:"targets"`
}

type ResponseError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type ResponseSuccess struct {
	Succeeded bool   `json:"succeeded"`
	Message   string `json:"message"`
}
