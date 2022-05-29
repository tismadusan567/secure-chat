package connection

import (
	"crypto/rsa"
	"encoding/gob"
	"net"
)

type Message struct {
	Header    Request
	PublicKey rsa.PublicKey
	UID       int
	Payload   string
}

func SendMessage(conn net.Conn, message Message) error {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(&message)
	if err != nil {
		return err
	}
	//fmt.Println("Sending " + message.Payload)
	return nil
}

func ReceiveMessage(conn net.Conn) (Message, error) {
	decoder := gob.NewDecoder(conn)
	response := Message{}
	err := decoder.Decode(&response)
	if err != nil {
		return Message{}, err
	}
	//fmt.Println("Recieved" + response)
	return response, nil
}
