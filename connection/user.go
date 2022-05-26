package connection

import (
	"encoding/gob"
	"fmt"
	"os"
)

type User struct {
	ID      int
	Address string
	Port    string
}

var users []User

func newUser(address, port string) *User {
	user := User{ID: ServerParams.ID, Address: address, Port: port}
	ServerParams.ID++
	return &user
}

func addUser(user *User) {
	users = append(users, *user)
}

func GetUser(address, port string) *User {
	for _, us := range users {
		if address == us.Address {
			return &us
		}
	}
	u := newUser(address, port)
	addUser(u)
	return u
}

func SaveUsers() error {
	f, err := os.Create("users.bin")
	if err != nil {
		fmt.Println("Save users - create err")
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err = enc.Encode(users); err != nil {
		fmt.Printf("Save users - write err : %v\n", err)
		return err
	}
	return nil
}

func ReadUsers() error {
	f, err := os.Open("users.bin")
	if err != nil {
		fmt.Printf("Read users - open err: %v\n", err)
		return err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	users = nil

	if err = dec.Decode(&users); err != nil {
		fmt.Println("Read users - read err")
		return err
	}
	return nil
}

func PrintUsers() {
	fmt.Println(users)
}
