package sse

import (
	"errors"
	"fmt"
	"github.com/goal-web/contracts"
	"sync"
)

var (
	ConnectionDontExistsErr = errors.New("connection does not exist")
)

func NewSse() contracts.Sse {
	return &Sse{
		fdMutex:     sync.Mutex{},
		connMutex:   sync.Mutex{},
		connections: map[uint64]contracts.SseConnection{},
		count:       0,
	}
}

type Sse struct {
	fdMutex     sync.Mutex
	connMutex   sync.Mutex
	connections map[uint64]contracts.SseConnection
	count       uint64
}

func (sse *Sse) Add(connect contracts.SseConnection) {
	sse.connMutex.Lock()
	defer sse.connMutex.Unlock()
	sse.connections[connect.Fd()] = connect
}

func (sse *Sse) Count() uint64 {
	return sse.count
}

func (sse *Sse) GetFd() uint64 {
	sse.fdMutex.Lock()
	defer sse.fdMutex.Unlock()
	sse.count++
	var fd = sse.count
	return fd
}

func (sse *Sse) Close(fd uint64) error {
	var conn, exists = sse.connections[fd]
	if exists {
		sse.connMutex.Lock()
		defer sse.connMutex.Unlock()
		delete(sse.connections, fd)
		return conn.Close()
	}

	return ConnectionDontExistsErr
}

func (sse *Sse) Send(fd uint64, message any) error {
	var conn, exists = sse.connections[fd]
	if exists {
		return conn.Send(message)
	}

	return ConnectionDontExistsErr
}

func (sse *Sse) Broadcast(message any) error {
	var errorFds []uint64
	for fd, conn := range sse.connections {
		if err := conn.Send(message); err != nil {
			errorFds = append(errorFds, fd)
		}
	}
	if len(errorFds) > 0 {
		return fmt.Errorf("the following connection failed to send. [%v]", errorFds)
	}
	return nil
}
