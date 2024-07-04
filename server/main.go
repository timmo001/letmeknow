package main

import (
	"flag"
	"log"

	"timmo.dev/letmeknowio/server/websocket"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	websocket.Setup()
}
