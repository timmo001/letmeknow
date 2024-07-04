package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	types "letmeknowio.timmo.dev/types/websocket"
)

// TODO: Add user authentication, so only authenticated users can send messages
// TODO: Check if user is allowed to send messages to other users


var upgrader = websocket.Upgrader{} // use default options

type ConnectedClients []types.Client

var connectedClients ConnectedClients

func (cc ConnectedClients) Display() []string {
	var clients []string
	for _, client := range cc {
		clients = append(clients, client.Display())
	}
	return clients
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	// Add client to connectedClients
	connectedClients = append(connectedClients,
		types.Client{
			Connection: c,
		},
	)
	log.Println("Client connected:", c.RemoteAddr())

	log.Println("Connected clients:", connectedClients.Display())

	for {
		mt, messageIn, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("Recv: %s", messageIn)

		// Parse JSON
		var request map[string]interface{}
		err = json.Unmarshal(messageIn, &request)
		if err != nil {
			log.Println("Error parsing JSON:", err)
			// Send error message
			errString := err.Error()
			resp := types.ResponseError{
				Message: "Error parsing JSON",
				Error:   &errString,
			}
			message, err := json.Marshal(resp)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
				break
			}
			c.WriteMessage(mt, message)
			break
		}

		// Validate JSON contains type
		if _, ok := request["type"]; !ok {
			log.Println("Error: JSON does not contain type")
			// Send error message
			resp := types.ResponseError{Message: "Error: JSON does not contain type"}
			message, err := json.Marshal(resp)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
				break
			}
			c.WriteMessage(mt, message)
			break
		}

		// If type is not "register" or "message", send error message
		if request["type"] != "register" && request["type"] != "message" {
			log.Println("Error: JSON type is not 'register' or 'message'")
			// Send error message
			resp := types.ResponseError{Message: "Error: JSON type is not 'register' or 'message'"}
			message, err := json.Marshal(resp)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
				break
			}
			c.WriteMessage(mt, message)
			break
		}

		// If type is "register", register client
		if request["type"] == "register" {
			// Validate JSON contains userID
			if _, ok := request["userID"]; !ok {
				log.Println("Error: JSON does not contain userID")
				// Send error message
				resp := types.ResponseError{Message: "Error: JSON does not contain userID"}
				message, err := json.Marshal(resp)
				if err != nil {
					log.Println("Error marshalling JSON:", err)
					break
				}
				c.WriteMessage(mt, message)
				break
			}

			// Convert request to ClientRegistration
			clientRegistration := types.RequestRegister{
				UserID: request["userID"].(string),
			}

			// Set userID for client
			alreadyRegistered := false
			for i, client := range connectedClients {
				if client.Connection == c {
					// Check if userID is already registered
					if client.UserID != nil {
						log.Println("Error: Client already registered with userID:", *client.UserID)
						alreadyRegistered = true
						break
					}

					connectedClients[i].UserID = &clientRegistration.UserID
					break
				}
			}

			// Send success message
			var resp types.ResponseSuccess
			if alreadyRegistered {
				resp = types.ResponseSuccess{
					Succeeded: false,
					Message:   "Client already registered",
				}
			} else {
				resp = types.ResponseSuccess{
					Succeeded: true,
					Message:   "Client registered",
				}
			}
			message, err := json.Marshal(resp)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
				break
			}
			c.WriteMessage(mt, message)
			continue
		}

		// Request type is "message"

		// Check if client is registered with a userID
		var clientRegistered bool = false
		for _, client := range connectedClients {
			if client.Connection == c && client.UserID != nil {
				clientRegistered = true
				break
			}
		}
		if !clientRegistered {
			log.Println("Error: Client not registered")
			// Send error message
			resp := types.ResponseError{Message: "Error: Client not registered"}
			message, err := json.Marshal(resp)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
				break
			}
			c.WriteMessage(mt, message)
			continue
		}

		// Validate JSON contains message
		if _, ok := request["message"]; !ok {
			log.Println("Error: JSON is not of type Message")
			// Send error message
			resp := types.ResponseError{Message: "Error: JSON is not of type Message"}
			message, err := json.Marshal(resp)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
				break
			}
			c.WriteMessage(mt, message)
			break
		}

		// Convert request to Message
		message := types.RequestMessage{
			Message: request["message"].(string),
			Targets: []string{},
		}

		if _, ok := request["targets"]; ok {
			targets := request["targets"].([]interface{})
			for _, target := range targets {
				message.Targets = append(message.Targets, target.(string))
			}
		}

		// Prepare message to send to clients
		messageOut, err := json.Marshal(message)
		if err != nil {
			log.Println("Error marshalling JSON:", err)
			break
		}

		// Send message to all clients
		for _, client := range connectedClients {
			// Only send message to clients that are requested
			if message.Targets != nil && len(message.Targets) > 0 {
				found := false
				for _, target := range message.Targets {
					if target == *client.UserID {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			} else {
				// Don't send message to the client that sent it, if sending to all clients
				if client.Connection == c {
					continue
				}
			}

			err = client.Connection.WriteMessage(mt, messageOut)
			if err != nil {
				log.Println("Error writing message to client:", err)
				break
			}
		}

		// Send success message
		resp := types.ResponseSuccess{Succeeded: true, Message: "Message sent"}
		messageSuccess, err := json.Marshal(resp)
		if err != nil {
			log.Println("Error marshalling JSON:", err)
			break
		}

		// Send success message to client
		err = c.WriteMessage(mt, messageSuccess)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}

	log.Println("Client disconnected:", c.RemoteAddr())

	// Remove client from connectedClients
	for i, client := range connectedClients {
		if client.Connection == c {
			connectedClients = append(connectedClients[:i], connectedClients[i+1:]...)
		}
	}

	log.Println("Connected clients:", connectedClients.Display())
}
