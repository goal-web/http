package sse

import (
	"bufio"
	"github.com/goal-web/application"
	"github.com/goal-web/contracts"
	"github.com/goal-web/http"
	"github.com/goal-web/supports/logs"
	"time"
)

func New(path string, controller contracts.SseController) (string, any) {
	factory := application.Get("sse.factory").(contracts.SseFactory)
	factory.Register(path, NewSse())

	return path, func(request *http.Request, serializer contracts.Serializer, sseFactory contracts.SseFactory, config contracts.Config) any {
		sse := sseFactory.Sse(path)
		var fd = sse.GetFd()
		if err := controller.OnConnect(request, fd); err != nil {
			logs.WithError(err).WithFields(request.Fields()).WithField("fd", fd).Debug("sse.NewRouter: OnConnect failed")
			return nil
		}
		httpConfig := config.Get("http").(http.Config)

		request.Request.SetContentType("text/event-stream")
		request.Request.Response.Header.Set("Cache-Control", "no-cache")
		request.Request.Response.Header.Set("Connection", "keep-alive")
		request.Request.Response.Header.Set("Transfer-Encoding", "chunked")

		// 添加 CORS 头部，允许所有来源的请求
		request.Request.Response.Header.Set("Access-Control-Allow-Origin", "*")
		request.Request.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		request.Request.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type,Accept")

		var (
			messageChan = make(chan any)
			closeChan   = make(chan bool)
			conn        = NewConnection(messageChan, closeChan, fd)
		)
		sse.Add(conn)

		request.Request.SetBodyStreamWriter(func(w *bufio.Writer) {

			// Create a ticker to send heartbeat messages
			ticker := time.NewTicker(httpConfig.SseHeartBeat)

			defer func() {
				controller.OnClose(fd)
				close(messageChan)
				closeChan = nil
				messageChan = nil
				ticker.Stop()
				request.Request.SetConnectionClose()
			}()

			go func() {
				for {
					<-ticker.C
					_, err := w.WriteString("data: heartbeat\n\n")
					if err != nil {
						logs.WithError(err).WithField("fd", fd).Error("sse heartbeat: response.Write failed")
						_ = sse.Close(fd)
						return
					}
					err = w.Flush()
					if err != nil {
						logs.WithError(err).WithField("fd", fd).Error("sse heartbeat: response.Flush failed")
						_ = sse.Close(fd)
						return
					}
				}
			}()

			for {
				select {
				case message := <-messageChan:
					_, err := w.WriteString("data: " + handleMessage(message, serializer) + "\n\n")
					if err != nil {
						logs.WithError(err).
							WithField("message", message).WithField("fd", fd).
							Error("sse.NewRouter: response.Write failed")
						return
					}
					err = w.Flush()
					if err != nil {
						logs.WithError(err).
							WithField("message", message).WithField("fd", fd).
							Error("sse.NewRouter: response.Write failed")
						return
					}
					//response todo
				case <-closeChan:
					return
				// connection is closed then defer will be executed
				case <-request.Request.Done():
					_ = sse.Close(fd)
					return
				}
			}
		})
		return nil
	}
}

func handleMessage(msg any, serializer contracts.Serializer) string {
	switch v := msg.(type) {
	case []byte:
		return string(v)
	case string:
		return v
	default:
		return serializer.Serialize(msg)
	}
}
