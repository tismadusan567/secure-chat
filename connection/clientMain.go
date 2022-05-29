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
	recConn    net.Conn
	friendList map[int]rsa.PublicKey
}

var UserParams userParams

func initClient() {
	UserParams.privKey, UserParams.pubKey = GenerateKeyPair(2048)
	if Join(GServerAddr, GServerPort).Header == INVALIDRESP {
		panic("connection failed")
	}
	UserParams.friendList = make(map[int]rsa.PublicKey)
}

func listenMessages() {
	for {
		message, err := ReceiveMessage(UserParams.recConn)
		if err != nil {
			fmt.Println("listener handleconnection ReceiveMessage")
			return
		}
		cipherBytes := []byte(message.Payload)
		messageBytes := DecryptWithPrivateKey(cipherBytes, UserParams.privKey)
		msg := string(messageBytes)
		fmt.Println(msg)
	}
}

func StartClient() {
	initClient()
	go listenMessages()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "user":
			GetUsers()
		case "con":
			scanner.Scan()
			otherID := scanner.Text()
			Connect(otherID)
		case "send":
			scanner.Scan()
			otherID := scanner.Text()
			scanner.Scan()
			msg := scanner.Text()
			Send(otherID, msg)
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
	UserParams.clientUID = transceive(message).UID

	// create receiving connection
	UserParams.recConn, err = net.Dial("tcp4", address+":"+port)
	if err != nil {
		fmt.Println("dial error")
		return nil
	}
	message2 := Message{
		Header: ESTABLISH,
		UID:    UserParams.clientUID,
	}

	err = SendMessage(UserParams.recConn, message2)
	if err != nil {
		fmt.Println("send message error")
		return nil
	}

	// receive msg
	response, err := ReceiveMessage(UserParams.recConn)
	if err != nil {
		fmt.Println("recv message error")
		return nil
	}
	return &response
}

// GetUsers get csv of ids of all clients in the network
func GetUsers() {
	message := Message{
		Header: USERS,
		UID:    UserParams.clientUID,
	}
	fmt.Println(*transceive(message))
}

func Connect(otherID string) {
	idNum, err := strconv.Atoi(otherID)
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
		Payload: otherID,
	}
	otherPubKey := transceive(message).PublicKey
	UserParams.friendList[idNum] = otherPubKey
	fmt.Println(otherPubKey)
}

func Send(otherID string, msg string) {
	idNum, err := strconv.Atoi(otherID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if idNum == UserParams.clientUID {
		fmt.Println("cannot send to oneself")
		return
	}
	if _, ok := UserParams.friendList[idNum]; !ok {
		fmt.Println("user not in friend list")
		return
	}

	msgByteArr := []byte(msg)
	otherPubKey := UserParams.friendList[idNum]
	cipherText := string(EncryptWithPublicKey(msgByteArr, &otherPubKey))
	message := Message{
		Header:    MESSAGE,
		UID:       UserParams.clientUID,
		Payload:   cipherText,
		PublicKey: rsa.PublicKey{E: idNum}, // EVIL
	}
	response := transceive(message)
	fmt.Println(response.Payload)
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
