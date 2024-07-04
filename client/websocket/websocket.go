package websocket

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"letmeknowio.timmo.dev/types"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func Setup() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/websocket"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Read:", err)
				return
			}
			log.Printf("Recv: %s", message)
		}
	}()

	// Send ping every 5 seconds
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	// Register for the websocket
	requestRegister := types.RequestRegister{
		Type:  "register",
		UserID: "test",
	}
	requestRegisterJSON, err := json.Marshal(requestRegister)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}
	err = c.WriteMessage(websocket.TextMessage, requestRegisterJSON)
	if err != nil {
		log.Println("Error writing register message:", err)
		return
	}

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.PingMessage, []byte(t.String()))
			if err != nil {
				log.Println("Error writing ping message:", err)
				return
			}
			log.Println("Ping")
		case <-interrupt:
			log.Println("Interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

}
