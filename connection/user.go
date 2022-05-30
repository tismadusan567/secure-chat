package connection

import (
	"crypto/rsa"
	"fmt"
	"net"
)

type User struct {
	ID         int
	Address    string
	Port       string
	PublicKey  rsa.PublicKey
	Connection net.Conn
}

var users []User

func NewUser(address, port string, pubKey rsa.PublicKey) {
	user := User{
		ID:        ServerParams.ID,
		Address:   address,
		Port:      port,
		PublicKey: pubKey,
	}
	ServerParams.ID++
	users = append(users, user)
}

func GetUser(uid int) *User {
	for i, us := range users {
		if uid == us.ID {
			return &users[i]
		}
	}
	return nil
}

func RemoveIndex(s []User, index int) []User {
	return append(s[:index], s[index+1:]...)
}

func DeleteUserByConn(conn net.Conn) {
	for i, us := range users {
		if conn == us.Connection {
			fmt.Println("User", us.ID, "disconnected")
			users = RemoveIndex(users, i)
			return
		}
	}
}

func GetUserIDs() (ids []int) {
	for i := range users {
		ids = append(ids, users[i].ID)
	}
	return ids
}
