module github.com/timmo001/letmeknow/server

go 1.22.4

replace github.com/timmo001/letmeknow/types v0.0.0 => ../types

require (
	github.com/gorilla/websocket v1.5.3
	github.com/martinlindhe/notify v0.0.0-20181008203735-20632c9a275a
	github.com/timmo001/letmeknow/types v0.0.0
)

require (
	github.com/deckarep/gosx-notifier v0.0.0-20180201035817-e127226297fb // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	gopkg.in/toast.v1 v1.0.0-20180812000517-0a84660828b2 // indirect
)
