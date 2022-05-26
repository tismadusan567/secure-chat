package connection

import (
	"fmt"
	"net"
)

func Join(address, port string) bool {
	conn, err := net.Dial("tcp4", address+":"+port)
	if err != nil {
		fmt.Println("dial error")
		return false
	}
	defer conn.Close()

	fmt.Println(conn.LocalAddr().String())
	message := NewMessage(JOIN, conn.LocalAddr().String())
	err = SendMessage(conn, message)
	if err != nil {
		fmt.Println("send message error")
		return false
	}

	rcmsg, err := RecieveMessage(conn)
	if err != nil {
		fmt.Println("recv message error")
		return false
	}
	fmt.Println(rcmsg)
	return true
}
