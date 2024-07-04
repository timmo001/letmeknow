package main

import (
	"flag"
	"log"

	websocketClient "github.com/timmo001/letmeknow/client/websocket"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	log.Println("Starting client...")

	websocketClient.Run()
}
