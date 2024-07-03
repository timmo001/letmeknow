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
		var data map[string]interface{}
		err = json.Unmarshal(messageIn, &data)
		if err != nil {
			log.Println("Error parsing JSON:", err)
			// Send error message
			resp := types.Error{Message: "Error parsing JSON", Error: err.Error()}
			message, err := json.Marshal(resp)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
				break
			}
			c.WriteMessage(mt, message)
			break
		}

		// Print JSON
		log.Printf("data: %v", data)

		// Validate JSON is of type Message
		if _, ok := data["message"]; !ok {
			log.Println("Error: JSON is not of type Message")
			// Send error message
			resp := types.Error{Message: "Error: JSON is not of type Message", Error: ""}
			message, err := json.Marshal(resp)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
				break
			}
			c.WriteMessage(mt, message)
			break
		}

		// TODO: Do something with the Message data

		// Convert JSON to bytes
		messageIn, err = json.Marshal(data)
		if err != nil {
			log.Println("Error marshalling JSON:", err)
			break
		}

		// Send message to all clients
		for _, client := range connectedClients {
			// Don't send message to the client that sent it
			if client == c {
				continue
			}

			err = client.WriteMessage(mt, messageIn)
			if err != nil {
				log.Println("Error writing message:", err)
				break
			}
		}

		// Send success message
		resp := types.Success{Succeeded: true, Message: "Message sent"}
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
