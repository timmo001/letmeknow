package websocket

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"timmo.dev/letmeknowio/server/types"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var connectedClients []*websocket.Conn

func sendMessage(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	// Add client to connectedClients
	connectedClients = append(connectedClients, c)

	for {
		mt, messageIn, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", messageIn)

		// Parse JSON
		var request map[string]interface{}
		err = json.Unmarshal(messageIn, &request)
		if err != nil {
			log.Println("Error parsing JSON:", err)
			// Send error message
			resp := types.ResponseError{Message: "Error parsing JSON", Error: err.Error()}
			message, err := json.Marshal(resp)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
				break
			}
			c.WriteMessage(mt, message)
			break
		}

		// Validate JSON is of type Message
		if _, ok := request["message"]; !ok {
			log.Println("Error: JSON is not of type Message")
			// Send error message
			resp := types.ResponseError{Message: "Error: JSON is not of type Message", Error: ""}
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
			Targets: request["targets"].([]string),
		}

		// TODO: Do something with the Message

		// Send message to all clients
		for _, client := range connectedClients {
			// Only send message to clients that are requested
			if len(message.Targets) > 0 {
				found := false
				for _, target := range message.Targets {
					if target == client.RemoteAddr().String() {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			} else {
				// Don't send message to the client that sent it, if sending to all clients
				if client == c {
					continue
				}
			}

			err = client.WriteMessage(mt, messageIn)
			if err != nil {
				log.Println("Error writing message to client:", err)
				break
			}
		}

		// Send success message
		resp := types.ResponseSuccess{Succeeded: true, Message: "Message sent"}
		messageOut, err := json.Marshal(resp)
		if err != nil {
			log.Println("Error marshalling JSON:", err)
			break
		}

		// Send success message to client
		err = c.WriteMessage(mt, messageOut)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}

	// Remove client from connectedClients
	for i, client := range connectedClients {
		if client == c {
			connectedClients = append(connectedClients[:i], connectedClients[i+1:]...)
			break
		}
	}
}

func Setup() {
	http.HandleFunc("/send", sendMessage)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
