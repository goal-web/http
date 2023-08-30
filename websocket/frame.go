package websocket

import (
	"github.com/goal-web/contracts"
)

type Frame struct {
	contracts.WebSocketConnection
	raw        []byte
	serializer contracts.Serializer
}

func NewFrame(raw []byte, conn contracts.WebSocketConnection, serializer contracts.Serializer) contracts.WebSocketFrame {
	return &Frame{
		WebSocketConnection: conn,
		raw:                 raw,
		serializer:          serializer,
	}
}

func (frame *Frame) Connection() contracts.WebSocketConnection {
	return frame.WebSocketConnection
}

func (frame *Frame) Raw() []byte {
	return frame.raw
}

func (frame *Frame) RawString() string {
	return string(frame.raw)
}

func (frame *Frame) Parse(v any) error {
	return frame.serializer.Unserialize(frame.RawString(), v)
}
