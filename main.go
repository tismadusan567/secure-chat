package main

import (
	"bufio"
	"fmt"
	"os"
	"secure_chat/connection"
)

func main() {
	fmt.Println("Enter s for server, c for client...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	fmt.Println("---------------------------------------")
	switch input {
	case "s":
		connection.StartServer()
	case "c":
		connection.StartClient()
	}
}
