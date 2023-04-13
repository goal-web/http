package sse

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
)

type Connection struct {
	fd        uint64
	msgPipe   chan any
	closePipe chan bool
}

func NewConnection(pipe chan any, closePipe chan bool, fd uint64) contracts.SseConnection {
	return &Connection{
		fd:        fd,
		msgPipe:   pipe,
		closePipe: closePipe,
	}
}

func (conn *Connection) Fd() uint64 {
	return conn.fd
}

func (conn *Connection) Close() error {
	conn.closePipe <- true
	return nil
}

func (conn *Connection) Send(msg any) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = Exception{exceptions.WithRecover(v)}
		}
	}()
	conn.msgPipe <- msg
	return
}
