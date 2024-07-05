package window

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	typeNotification "github.com/timmo001/letmeknow/types/notification"
)

// TODO: Move window to top left corner (w/ margin) and make it always on top

var appInstance fyne.App
var window fyne.Window

func Setup(n typeNotification.Notification) {
	log.Printf("Setting up GUI with notification: %+v\n", n.Display())

	appInstance = app.New()
	window = appInstance.NewWindow(*n.Title + " - LetMeKnow")

	window.Resize(fyne.NewSize(270, 40))
	// window.SetFixedSize(true)

	layout := container.New(layout.NewVBoxLayout())

	if n.Title != nil {
		textTitle := canvas.NewText(*n.Title, color.White)
		textTitle.TextStyle.Bold = true
		textTitle.TextSize = 22
		layout.Objects = append(layout.Objects, textTitle)
	}
	if n.Subtitle != nil {
		textSubtitle := canvas.NewText(*n.Subtitle, color.White)
		textSubtitle.TextStyle.Italic = true
		textSubtitle.TextSize = 15
		layout.Objects = append(layout.Objects, textSubtitle)
	}
	if n.Content != nil {
		textContent := canvas.NewText(*n.Content, color.White)
		textContent.TextSize = 13
		layout.Objects = append(layout.Objects, textContent)
	}

	if n.Image != nil {
		uri, err := storage.ParseURI(n.Image.URL)
		if err != nil {
			log.Println("Error parsing URI:", err)
		} else {
			canvasImage := canvas.NewImageFromURI(uri)
			canvasImage.Resize(fyne.NewSize(n.Image.Width, n.Image.Height))
			canvasImage.FillMode = canvas.ImageFillOriginal
			canvasImage.ScaleMode = canvas.ImageScaleSmooth
			layout.Objects = append(layout.Objects, canvasImage)
		}
	}

	window.SetContent(layout)
}

func Run() {
	log.Println("Starting GUI...")

	window.Show()
	window.RequestFocus()

	appInstance.Run()

	fmt.Println("Exited GUI")
}
