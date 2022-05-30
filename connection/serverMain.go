package connection

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type serverParams struct {
	ID int
}

var ServerParams serverParams

func initServer() {
	ServerParams.ID = 1
}

func StartServer() {
	initServer()
	ln, err := net.Listen("tcp4", ":"+GServerPort)
	if err != nil {
		fmt.Printf("startListening error: %v", err)
		return
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("connection accept error: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		message, err := ReceiveMessage(conn)
		if err != nil {
			fmt.Println("connection closed")
			return
		}
		handleMessage(message, conn)
	}
}

func handleMessage(message Message, conn net.Conn) {
	switch message.Header {
	case JOIN:
		handleJoin(message, conn)
	case ESTABLISH:
		handleEstablish(message, conn)
	case USERS:
		handleUsers(message, conn)
	case CONNECT:
		handleConnect(message, conn)
	case MESSAGE:
		handleSend(message, conn)
	default:
		handleInvalidRequest(conn, "invalid request")
	}
}

func getSenderAddress(conn net.Conn) (string, string) {
	addr := conn.RemoteAddr().String()
	split := strings.Split(addr, ":")
	ip := split[0]
	port := split[1]
	return ip, port
}

func handleJoin(message Message, conn net.Conn) {
	ip, port := getSenderAddress(conn)
	NewUser(ip, port, message.PublicKey)
	response := Message{Header: JOINRESP, UID: ServerParams.ID - 1}
	err := SendMessage(conn, response)
	if err != nil {
		fmt.Println("server join response error")
		return
	}
}

func handleEstablish(message Message, conn net.Conn) {
	user := GetUser(message.UID)
	user.Connection = conn
	response := Message{
		Header: ESTABLISHRESP,
	}
	err := SendMessage(conn, response)
	if err != nil {
		fmt.Println("server establish response error")
		return
	}
}

func handleUsers(message Message, conn net.Conn) {
	user := GetUser(message.UID)
	if user == nil {
		handleInvalidRequest(conn, "User not in network")
		return
	}
	ids := GetUserIDs()
	var sb strings.Builder
	for i := range ids {
		sb.WriteRune(rune(ids[i] + '0'))
		sb.WriteByte(',')
	}
	response := Message{Header: USERSRESP, Payload: sb.String()}
	err := SendMessage(conn, response)
	if err != nil {
		fmt.Println("server users response error")
		return
	}
}

func handleConnect(message Message, conn net.Conn) {
	user := GetUser(message.UID)
	if user == nil {
		handleInvalidRequest(conn, "user not in network")
		return
	}

	id, err := strconv.Atoi(message.Payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	otherUser := GetUser(id)
	if otherUser == nil {
		handleInvalidRequest(conn, "other user not in network")
		return
	}
	pubKey := otherUser.PublicKey
	response := Message{Header: CONNECTRESP, PublicKey: pubKey}
	err = SendMessage(conn, response)
	if err != nil {
		fmt.Println("invalid response error")
		return
	}
}

func handleSend(message Message, conn net.Conn) {
	user := GetUser(message.UID)
	if user == nil {
		handleInvalidRequest(conn, "user not in network")
		return
	}
	id := message.PublicKey.E // EVIL
	otherUser := GetUser(id)
	if otherUser == nil {
		handleInvalidRequest(conn, "other user not in network")
		return
	}
	forwardMessage := Message{
		Header:  FORWARD,
		UID:     message.UID,
		Payload: message.Payload,
	}
	err := SendMessage(otherUser.Connection, forwardMessage)
	if err != nil {
		handleInvalidRequest(conn, "message not sent")
		return
	}
	response := Message{
		Header:  MESSAGERESP,
		Payload: "Message sent",
	}
	err = SendMessage(conn, response)
	if err != nil {
		fmt.Println("invalid response error")
		return
	}
}

func handleInvalidRequest(conn net.Conn, s string) {
	response := Message{Header: INVALIDRESP, Payload: s}
	err := SendMessage(conn, response)
	if err != nil {
		fmt.Println("invalid response error")
	}
}
