package notification

import (
	types "letmeknow.timmo.dev/types/notification"

	"github.com/martinlindhe/notify"
)

func Notify(n types.Notification) {
	if n.Alert {
		notify.Alert(n.Name, n.Title, n.Message, n.Icon)
	} else {
		notify.Notify(n.Name, n.Title, n.Message, n.Icon)
	}
}
