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

func sendMessage(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		// Parse JSON
		var data map[string]interface{}
		err = json.Unmarshal(message, &data)
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

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

func Setup() {
	http.HandleFunc("/send", sendMessage)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
