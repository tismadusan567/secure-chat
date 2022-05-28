package connection

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
)

type User struct {
	ID      byte
	Address string
	Port    string
}

var users []User

func NewUser(address, port string) error {
	if GetUser(address, port) != nil {
		return errors.New("user already exists")
	}
	user := User{ID: ServerParams.ID, Address: address, Port: port}
	ServerParams.ID++
	users = append(users, user)
	return nil
}

func GetUser(address, uid string) *User {
	for _, us := range users {
		if uid == string(us.ID) /* && address == us.Address */ {
			return &us
		}
	}
	return nil
}

func GetUserIDs() (ids []byte) {
	for i := range users {
		ids = append(ids, users[i].ID)
	}
	return ids
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
