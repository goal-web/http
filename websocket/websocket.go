package websocket

import (
	"errors"
	"github.com/goal-web/contracts"
	"sync"
)

var (
	ConnectionDontExistsErr = errors.New("connection does not exist")
)

type WebSocket struct {
	fdMutex     sync.Mutex
	connMutex   sync.Mutex
	connections map[uint64]contracts.WebSocketConnection
	count       uint64
}

func (ws *WebSocket) Add(connect contracts.WebSocketConnection) {
	ws.connMutex.Lock()
	defer ws.connMutex.Unlock()
	ws.connections[connect.Fd()] = connect
}

func (ws *WebSocket) GetFd() uint64 {
	ws.fdMutex.Lock()
	defer ws.fdMutex.Unlock()
	ws.count++
	var fd = ws.count
	return fd
}

func (ws *WebSocket) Close(fd uint64) error {
	var conn, exists = ws.connections[fd]
	if exists {
		ws.connMutex.Lock()
		defer ws.connMutex.Unlock()
		delete(ws.connections, fd)
		return conn.Close()
	}

	return ConnectionDontExistsErr
}

func (ws *WebSocket) Send(fd uint64, message any) error {
	var conn, exists = ws.connections[fd]
	if exists {
		return conn.Send(message)
	}

	return ConnectionDontExistsErr
}
