module github.com/timmo001/letmeknow/server

go 1.22.4

replace github.com/timmo001/letmeknow/types v0.0.0 => ../types

require (
	github.com/gorilla/websocket v1.5.3
	github.com/timmo001/letmeknow/types v0.0.0
)
