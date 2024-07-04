package main

import (
	"flag"
	"log"

	"letmeknowio.timmo.dev/client/websocket"
)


func main() {
	flag.Parse()
	log.SetFlags(0)

	websocket.Setup()
}
