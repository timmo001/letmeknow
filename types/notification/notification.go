package notification

import "fmt"

type Image struct {
	Height float32
	Width  float32
	URL    string
}

func (i Image) Display() string {
	return fmt.Sprintf("  Height: %f\n  Width: %f\n  URL: %s", i.Height, i.Width, i.URL)
}

type Notification struct {
	Title    *string
	Subtitle *string
	Content  *string
	Image    *Image
}

func (n Notification) Display() string {
	return fmt.Sprintf("Title: %s\nSubtitle: %s\nContent: %s\nImage:\n%s", *n.Title, *n.Subtitle, *n.Content, n.Image.Display())
}
