package connection

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
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
	fmt.Println("Joined server with id:", UserParams.clientUID)
	fmt.Println("Commands:")
	fmt.Println("users - get list of all currently online users")
	fmt.Println("conn [u] - connect with user u")
	fmt.Println("send [u] [msg] - send msg to user u")
	PrintPadding()
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
		fmt.Println("Received message from user", message.UID)
		fmt.Println(msg)
		PrintPadding()
	}
}

func StartClient() {
	initClient()
	go listenMessages()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input := scanner.Text()
		split := strings.Split(input, " ")
		switch split[0] {
		case "users":
			GetUsers()
		case "conn":
			if len(split) == 2 {
				otherID := split[1]
				Connect(otherID)
			} else {
				fmt.Println("input error")
				PrintPadding()
			}
		case "send":
			if len(split) == 3 {
				otherID := split[1]
				msg := split[2]
				Send(otherID, msg)
			} else {
				fmt.Println("input error")
				PrintPadding()
			}
		}
	}
}

// Join get unique user ID
func Join(address, port string) *Message {
	// establish send connection to server
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

	// establish receiving connection
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
		PrintPadding()
		return nil
	}

	// receive msg
	response, err := ReceiveMessage(UserParams.recConn)
	if err != nil {
		fmt.Println("recv message error")
		PrintPadding()
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
	response := transceive(message)
	split := strings.Split(response.Payload, ",")
	fmt.Println("List of online users:")
	for i := range split {
		fmt.Printf("%s ", split[i])
	}
	fmt.Printf("\n")
	PrintPadding()
}

func Connect(otherID string) {
	idNum, err := strconv.Atoi(otherID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if idNum == UserParams.clientUID {
		fmt.Println("cannot connect to oneself")
		PrintPadding()
		return
	}

	message := Message{
		Header:  CONNECT,
		UID:     UserParams.clientUID,
		Payload: otherID,
	}
	response := transceive(message)
	if response.Header == INVALIDRESP {
		fmt.Println(response.Payload)
		PrintPadding()
		return
	}
	fmt.Println("Connection successful")
	PrintPadding()
	otherPubKey := response.PublicKey
	UserParams.friendList[idNum] = otherPubKey
}

func Send(otherID string, msg string) {
	idNum, err := strconv.Atoi(otherID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if idNum == UserParams.clientUID {
		fmt.Println("cannot send to oneself")
		PrintPadding()
		return
	}
	if _, ok := UserParams.friendList[idNum]; !ok {
		fmt.Println("user not in friend list")
		PrintPadding()
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
	PrintPadding()
}

func transceive(message Message) *Message {
	// send msg
	err := SendMessage(UserParams.conn, message)
	if err != nil {
		fmt.Println("send message error")
		PrintPadding()
		return nil
	}

	// receive msg
	response, err := ReceiveMessage(UserParams.conn)
	if err != nil {
		fmt.Println("recv message error")
		PrintPadding()
		return nil
	}
	return &response
}
