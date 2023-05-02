package http

type Config struct {
	Address           string
	Host              string
	Port              string
	GlobalMiddlewares []any

	StaticDirectories map[string]string
}
