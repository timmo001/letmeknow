package main

import (
	"flag"
	"log"
	"net/http"

	"timmo.dev/letmeknowio/server/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	log.Println("Starting server on", *addr)

	http.HandleFunc("/websocket", websocket.WebSocket)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
