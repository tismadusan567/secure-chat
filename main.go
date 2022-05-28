package main

import (
	"bufio"
	"os"
	"secure_chat/connection"
)

func main() {
	connection.ServerParams.ADDR = "localhost"
	connection.ServerParams.PORT = "4420"
	connection.ServerParams.ID = '1'
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	switch input {
	case "s":
		connection.StartServer()
	case "c":
		connection.StartClient()
	}
}
