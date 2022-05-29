package connection

type Request string

const GServerAddr = "localhost"
const GServerPort = "4420"

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
