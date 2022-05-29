package connection

type Request string

const GServerAddr = "localhost"
const GServerPort = "4420"

// enum for requests and responses
const (
	JOIN      = "join"
	ESTABLISH = "establish"
	USERS     = "users"
	CONNECT   = "connect"
	MESSAGE   = "message"
	TEST      = "test"

	JOINRESP      = "joinresp"
	ESTABLISHRESP = "establishresp"
	USERSRESP     = "usersresp"
	CONNECTRESP   = "connectresp"
	MESSAGERESP   = "messageresp"
	FORWARD       = "forward"
	INVALIDRESP   = "invalidresp"
	TESTRESP      = "testresp"
)
