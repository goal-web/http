package http

type Config struct {
	Address string
	Host    string
	Port    string

	StaticDirectories map[string]string
}
