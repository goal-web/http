package websocket

import "github.com/fasthttp/websocket"

type Config struct {
	Upgrader websocket.FastHTTPUpgrader
}
