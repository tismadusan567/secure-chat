package connection

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var clientUID string

func StartClient() {
	clientUID = Join(ServerParams.ADDR, ServerParams.PORT).Payload
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "users":
			fmt.Println(*GetUsers(ServerParams.ADDR, ServerParams.PORT))
		}
	}
}

// Join get unique user ID
func Join(address, port string) *Message {
	return dial(address, port, JOIN, "")
}

// GetUsers get csv of ids of all clients in the network
func GetUsers(address, port string) *Message {
	fmt.Printf("Online: %v\n\n", clientUID)
	return dial(address, port, USERS, clientUID)
}

func dial(address, port string, header Request, payload string) *Message {
	// dial
	conn, err := net.Dial("tcp4", address+":"+port)
	if err != nil {
		fmt.Println("dial error")
		return nil
	}
	defer conn.Close()

	// send msg
	fmt.Println(conn.LocalAddr().String())
	message := NewMessage(header, payload)
	err = SendMessage(conn, message)
	if err != nil {
		fmt.Println("send message error")
		return nil
	}

	// receive msg
	response, err := ReceiveMessage(conn)
	if err != nil {
		fmt.Println("recv message error")
		return nil
	}
	return &response
}
