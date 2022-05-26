package connection

type Request string

// enum for requests
const (
	JOIN    = "join"
	USERS   = "users"
	CONNECT = "connect"
	MESSAGE = "message"
	TEST    = "test"
)
