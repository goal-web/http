package http

import "time"

type Config struct {
	Address           string
	Host              string
	Port              string
	GlobalMiddlewares []any

	StaticDirectories map[string]string
	SseHeartBeat      time.Duration
}
