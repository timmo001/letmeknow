package main

import (
	"flag"
	"log"

	"letmeknow.timmo.dev/client/websocket"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	websocket.Run()
}
