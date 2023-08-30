package websocket

import (
	"github.com/fasthttp/websocket"
	"github.com/goal-web/contracts"
	"github.com/goal-web/http"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/logs"
)

var (
	upgrader = websocket.FastHTTPUpgrader{}
)

func New(controller contracts.WebSocketController) any {
	return func(request *http.Request, serializer contracts.Serializer, socket contracts.WebSocket, handler contracts.ExceptionHandler) error {
		var err = upgrader.Upgrade(request.Request, func(ws *websocket.Conn) {
			var fd = socket.GetFd()

			if err := controller.OnConnect(request, fd); err != nil {
				logs.WithError(err).Error("websocket.NewRouter: OnConnect failed")
			}

			var conn = NewConnection(ws, fd)
			socket.Add(conn)

			defer func() {
				controller.OnClose(fd)
				if closeErr := socket.Close(conn.Fd()); closeErr != nil {
					logs.WithError(closeErr).Error("websocket.NewRouter: Connection close failed")
				}
			}()

			for {
				// Read
				var msgType, msg, readErr = ws.ReadMessage()
				if readErr != nil {
					logs.WithError(readErr).Error("websocket.NewRouter: Failed to read message")
					break
				}

				switch msgType {
				case websocket.TextMessage, websocket.BinaryMessage:
					go handleMessage(NewFrame(msg, conn, serializer), controller, handler)
				case websocket.CloseMessage:
					break
				}
			}
		})

		if err != nil {
			logs.WithError(err).Error("websocket.NewRouter: Upgrade failed")
		}

		return err
	}
}

func handleMessage(frame contracts.WebSocketFrame, controller contracts.WebSocketController, handler contracts.ExceptionHandler) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			handler.Handle(Exception{
				Exception: exceptions.WithRecover(panicValue),
			})
		}
	}()
	controller.OnMessage(frame)
}
