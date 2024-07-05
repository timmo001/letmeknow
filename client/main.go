package main

import (
	"flag"
	"log"

	websocketClient "github.com/timmo001/letmeknow/client/websocket"
	typeNotification "github.com/timmo001/letmeknow/types/notification"

	"github.com/timmo001/letmeknow/client/window"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	log.Println("Starting client...")

	go websocketClient.Run()

	title := "Title"
	subtitle := "Subtitle"
	content := "Content"
	image := typeNotification.Image{
		Height: 480,
		Width:  270,
		URL:    "https://placehold.co/480x270",
	}
	window.Setup(typeNotification.Notification{
		Title:    &title,
		Subtitle: &subtitle,
		Content:  &content,
		Image:    &image,
	})
	window.Run()
}
