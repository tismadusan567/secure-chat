package connection

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"net"
	"os"
)

var clientUID string
var privKey *rsa.PrivateKey
var pubKey *rsa.PublicKey

func StartClient() {
	privKey, pubKey = GenerateKeyPair(2048)
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
	message := Message{Header: JOIN, PublicKey: *pubKey}
	return dial(address, port, message)
}

// GetUsers get csv of ids of all clients in the network
func GetUsers(address, port string) *Message {
	message := Message{Header: USERS, Payload: clientUID}
	return dial(address, port, message)
}

func dial(address, port string, message Message) *Message {
	// dial
	conn, err := net.Dial("tcp4", address+":"+port)
	if err != nil {
		fmt.Println("dial error")
		return nil
	}
	defer conn.Close()

	// send msg
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
