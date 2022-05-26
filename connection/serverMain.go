package connection

import (
	"fmt"
	"net"
	"strings"
)

type serverParams struct {
	ID int
}

var ServerParams serverParams

func StartServer() {
	ServerParams.ID = 0
	ln, err := net.Listen("tcp4", ":4420")
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
	message, err := RecieveMessage(conn)
	if err != nil {
		panic("listener handleconnection RecieveMessage")
	}
	handleMessage(message, conn)
}

func handleMessage(message Message, conn net.Conn) {
	switch message.Header {
	case JOIN:
		handleJoin(message, conn)
	case USERS:
		fmt.Println("Handle sendfile")
	case CONNECT:
		fmt.Println("Handle requestfile")
	case MESSAGE:
		fmt.Println("Handle requestfile")
	default:
		fmt.Println("Invalid request")
	}
}

func handleJoin(message Message, conn net.Conn) {
	addr := conn.RemoteAddr().String()
	ip := strings.TrimRight(addr, ":")
	port := strings.TrimLeft(addr, ":")
	GetUser(ip, port)

	response := NewMessage(TEST, "odg")
	err := SendMessage(conn, response)
	if err != nil {
		fmt.Println(err)
		panic("listener handlePing SendMessage")
	}
}
