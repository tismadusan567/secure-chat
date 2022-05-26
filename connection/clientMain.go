package connection

import "fmt"

func StartClient() {
	online := Join("localhost", "4420")
	fmt.Printf("Online: %v\n\n", online)
}
