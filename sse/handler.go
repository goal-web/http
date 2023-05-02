package sse

import (
	"bytes"
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/http"
	"github.com/goal-web/supports/logs"
)

func New(controller contracts.SseController) any {
	return func(request *http.Request, serializer contracts.Serializer, sse contracts.Sse) error {
		var fd = sse.GetFd()
		if err := controller.OnConnect(request, fd); err != nil {
			logs.WithError(err).WithFields(request.Fields()).WithField("fd", fd).Debug("sse.NewRouter: OnConnect failed")
			return err
		}

		var response = request.Request().Response
		response.Header.Set("Content-Type", "text/event-stream")
		response.Header.Set("Cache-Control", "no-cache")
		response.Header.Set("Connection", "keep-alive")
		response.Header.Set("Access-Control-Allow-Origin", "*")

		var (
			messageChan = make(chan any)
			closeChan   = make(chan bool)
			conn        = NewConnection(messageChan, closeChan, fd)
		)
		sse.Add(conn)

		defer func() {
			controller.OnClose(fd)
			close(messageChan)
			close(closeChan)
			closeChan = nil
			messageChan = nil
		}()

		for {
			select {
			case message := <-messageChan:
				var err = response.Write(bytes.NewBuffer([]byte(fmt.Sprintf("%s\n", handleMessage(message, serializer)))))
				//var _, err = fmt.Fprintf(response, "%s\n", handleMessage(message, serializer))
				if err != nil {
					logs.WithError(err).
						WithField("message", message).WithField("fd", fd).
						Error("sse.NewRouter: response.Write failed")
					return nil
				}
				//response todo

			case <-closeChan:
				return nil

			// connection is closed then defer will be executed
			case <-request.Request().Context().Done():
				_ = sse.Close(fd)
				return nil
			}
		}
	}
}

func handleMessage(msg any, serializer contracts.Serializer) []byte {
	switch v := msg.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	default:
		return []byte(serializer.Serialize(msg))
	}
}
