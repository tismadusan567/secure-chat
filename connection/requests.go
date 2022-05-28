package connection

type Request string

// enum for requests and responses
const (
	JOIN    = "join"
	USERS   = "users"
	CONNECT = "connect"
	MESSAGE = "message"
	TEST    = "test"

	JOINRESP    = "joinresp"
	USERSRESP   = "usersresp"
	CONNECTRESP = "connectresp"
	MESSAGERESP = "messageresp"
	INVALIDRESP = "invalidresp"
	TESTRESP    = "testresp"
)
