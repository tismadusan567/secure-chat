package connection

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"net"
	"os"
	"strconv"
)

type userParams struct {
	clientUID  int
	privKey    *rsa.PrivateKey
	pubKey     *rsa.PublicKey
	conn       net.Conn
	friendList map[int]rsa.PublicKey
}

var UserParams userParams

func initClient() {
	UserParams.privKey, UserParams.pubKey = GenerateKeyPair(2048)
	UserParams.clientUID = Join(GServerAddr, GServerPort).UID
	UserParams.friendList = make(map[int]rsa.PublicKey)
}

func StartClient() {
	initClient()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "us":
			GetUsers()
		case "con":
			scanner.Scan()
			input := scanner.Text()
			Connect(input)
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

	message := Message{
		Header:    JOIN,
		PublicKey: *UserParams.pubKey,
	}
	return transceive(message)
}

// GetUsers get csv of ids of all clients in the network
func GetUsers() {
	message := Message{
		Header: USERS,
		UID:    UserParams.clientUID,
	}
	fmt.Println(*transceive(message))
}

func Connect(id string) {
	idNum, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	if idNum == UserParams.clientUID {
		fmt.Println("cannot connect to oneself")
		return
	}

	message := Message{
		Header:  CONNECT,
		UID:     UserParams.clientUID,
		Payload: id,
	}
	otherPubKey := transceive(message).PublicKey
	UserParams.friendList[idNum] = otherPubKey
	fmt.Println(otherPubKey)
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
