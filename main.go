package main

import (
	"bufio"
	"os"
	"secure_chat/connection"
)

func main() {
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
