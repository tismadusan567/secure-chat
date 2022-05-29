package connection

import (
	"crypto/rsa"
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

func NewUser(address, port string, pubKey rsa.PublicKey) error {
	user := User{
		ID:        ServerParams.ID,
		Address:   address,
		Port:      port,
		PublicKey: pubKey,
	}
	ServerParams.ID++
	users = append(users, user)
	return nil
}

func GetUser(uid int) *User {
	for i, us := range users {
		if uid == us.ID {
			return &users[i]
		}
	}
	return nil
}

func GetUserIDs() (ids []int) {
	for i := range users {
		ids = append(ids, users[i].ID)
	}
	return ids
}
