package connection

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"net"
	"os"
)

type userParams struct {
	clientUID int
	privKey   *rsa.PrivateKey
	pubKey    *rsa.PublicKey
	conn      net.Conn
}

var UserParams userParams

func initClient() {
	UserParams.privKey, UserParams.pubKey = GenerateKeyPair(2048)
	UserParams.clientUID = Join(GServerAddr, GServerPort).UID
}

func StartClient() {
	initClient()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "users":
			fmt.Println(*GetUsers())
		}
	}
}

// Join get unique user ID
func Join(address, port string) *Message {
	// dial
	conn, err := net.Dial("tcp4", address+":"+port)
	if err != nil {
		fmt.Println("dial error")
		return nil
	}
	UserParams.conn = conn

	message := Message{Header: JOIN, PublicKey: *UserParams.pubKey}
	return transceive(message)
}

// GetUsers get csv of ids of all clients in the network
func GetUsers() *Message {
	message := Message{Header: USERS, UID: UserParams.clientUID}
	return transceive(message)
}

func transceive(message Message) *Message {
	// send msg
	err := SendMessage(UserParams.conn, message)
	if err != nil {
		fmt.Println("send message error")
		return nil
	}

	// receive msg
	response, err := ReceiveMessage(UserParams.conn)
	if err != nil {
		fmt.Println("recv message error")
		return nil
	}
	return &response
}
