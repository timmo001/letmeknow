package types

type Error struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type Message struct {
	Message string `json:"message"`
}
