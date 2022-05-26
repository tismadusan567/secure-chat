package connection

import (
	"encoding/gob"
	"fmt"
	"net"
)

type Message struct {
	Header  Request
	Payload string
}

func NewMessage(request Request, data string) Message {
	message := Message{Header: request, Payload: data}
	return message
}

func SendMessage(conn net.Conn, message Message) error {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(&message)
	if err != nil {
		return err
	}
	fmt.Println("Sending " + message.Payload)
	return nil
}

func RecieveMessage(conn net.Conn) (Message, error) {
	decoder := gob.NewDecoder(conn)
	response := Message{}
	err := decoder.Decode(&response)
	if err != nil {
		return Message{}, err
	}
	fmt.Printf("Recieved : %+v\n", response)
	return response, nil
}
