package main

import (
	"bufio"
	"os"
	"secure_chat/connection"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input := scanner.Text()

		switch {
		case input == "server":
			connection.StartServer()
		case input == "client":
			connection.StartClient()
		}
	}
}
