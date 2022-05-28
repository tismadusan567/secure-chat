package connection

import (
	"fmt"
	"net"
	"strings"
)

type serverParams struct {
	ADDR string
	PORT string
	ID   byte
}

var ServerParams serverParams

func StartServer() {
	ln, err := net.Listen("tcp4", ":"+ServerParams.PORT)
	if err != nil {
		fmt.Errorf("StartListening error: %v", err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Errorf("greska pri prihvatanju konekcije: %v", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	message, err := ReceiveMessage(conn)
	if err != nil {
		panic("listener handleconnection ReceiveMessage")
	}
	handleMessage(message, conn)
}

func handleMessage(message Message, conn net.Conn) {
	switch message.Header {
	case JOIN:
		handleJoin(message, conn)
	case USERS:
		handleUsers(message, conn)
	case CONNECT:
		fmt.Println("Handle requestfile")
	case MESSAGE:
		fmt.Println("Handle requestfile")
	default:
		handleInvalidRequest(conn, "Invalid request")
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
	er := NewUser(ip, port)
	if er != nil {
		fmt.Println(er)
	}

	response := NewMessage(JOINRESP, string(ServerParams.ID-1))
	err := SendMessage(conn, response)
	if err != nil {
		fmt.Println("Server join response error")
		return
	}
}

func handleUsers(message Message, conn net.Conn) {
	ip, _ := getSenderAddress(conn)
	user := GetUser(ip, message.Payload)
	if user == nil {
		handleInvalidRequest(conn, "User not in network")
		return
	}
	ids := GetUserIDs()
	var sb strings.Builder
	for i := range ids {
		sb.WriteByte(ids[i])
		sb.WriteByte(',')
	}
	response := NewMessage(USERSRESP, sb.String())
	err := SendMessage(conn, response)
	if err != nil {
		fmt.Println("Server users response error")
		return
	}
}

func handleInvalidRequest(conn net.Conn, s string) {
	response := NewMessage(INVALIDRESP, s)
	err := SendMessage(conn, response)
	if err != nil {
		fmt.Println("Invalid response error")
	}
}
