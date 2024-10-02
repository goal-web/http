package http

type ServeClosed struct {
}

func (c *ServeClosed) Event() string {
	return "HTTP_SERVE_CLOSED"
}
