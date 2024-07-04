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
	types "letmeknowio.timmo.dev/types/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func Run() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		// Attempt to reconnect after 10 seconds
		select {
		case <-interrupt:
			log.Println("Interrupt from Run loop")
			return
		default:
			code := setup()

			if code == 0 {
				break
			}

			log.Printf("Attempting to reconnect in 10 seconds")
			// Make sleep interruptible
			select {
			case <-time.After(10 * time.Second):
				// Sleep for 10 seconds or until interrupt
			case <-interrupt:
				log.Println("Interrupt during sleep")
				return
			}
		}
	}
}

func setup() int {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/websocket"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("Dial:", err)
		return 1
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}
			log.Printf("Recv: %s", message)
		}
	}()

	// Send ping every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Register for the websocket
	requestRegister := types.RequestRegister{
		Type:   "register",
		UserID: "test",
	}
	requestRegisterJSON, err := json.Marshal(requestRegister)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return 1
	}
	err = c.WriteMessage(websocket.TextMessage, requestRegisterJSON)
	if err != nil {
		log.Println("Error writing register message:", err)
		return 1
	}

	for {
		select {
		case <-done:
			return 1
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.PingMessage, []byte(t.String()))
			if err != nil {
				log.Println("Error writing ping message:", err)
				return 1
			}
		case <-interrupt:
			log.Println("Interrupt from setup loop")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return 1
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return 0
		}
	}
}
