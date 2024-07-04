package main

import (
	"flag"
	"letmeknowio/client/websocket"
	"log"
)


func main() {
	flag.Parse()
	log.SetFlags(0)

	websocket.Setup()
}
