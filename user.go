package main

import (
	"encoding/json"
	"fmt"
	"image"
	"net"
	"net/http"
)

type User struct {
	Attribute uint32 `json:"attribute"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	ID        uint64 `json:"id"`
	Avatar    image.Image
	Rooms     map[uint64]bool `json:"rooms"`
}

type CreateUser struct {
	U User `json:"data"`
}

type UpdateUser struct {
	U User `json:"data"`
}

type ReadUser struct {
	Data struct {
		ID uint64 `json:"id"`
	} `json:"data"`
}

type DeleteUser struct {
	Data struct {
		ID uint64 `json:"id"`
	} `json:"data"`
}

type LoginUser struct {
	Data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"data"`
}

type LogoutUser struct {
	Data struct {
		ID uint64 `json:"id"`
	} `json:"data"`
}

func (u User) Create() DefinedAction {
	return &CreateUser{}
}

func (action *CreateUser) GetFromJSON(data []byte) {
	err := json.Unmarshal(data, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (action *CreateUser) Process(db *DB, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	db.AddUser(action, conn, w, req)
}

func (u User) Update() DefinedAction {
	return &UpdateUser{}
}

func (action *UpdateUser) GetFromJSON(data []byte) {
	err := json.Unmarshal(data, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (action *UpdateUser) Process(db *DB, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	db.UpdateUser(action, conn, w, req)
}

func (u User) Read() DefinedAction {
	return &ReadUser{}
}

func (action *ReadUser) GetFromJSON(data []byte) {
	err := json.Unmarshal(data, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (action *ReadUser) Process(db *DB, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	db.ReadUser(action, conn, w, req)
}

func (u User) Delete() DefinedAction {
	return &DeleteUser{}
}

func (action *DeleteUser) GetFromJSON(data []byte) {
	err := json.Unmarshal(data, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (action *DeleteUser) Process(db *DB, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	db.DeleteUser(action, conn, w, req)
}

func (u User) Print() {
	fmt.Printf("ID:%d Name:%s Attribute:%d Email:%s Password:%s", u.ID, u.Name, u.Attribute, u.Email, u.Password)
}
func (u User) GetId() uint64 {
	return u.ID
}

func (u User) Login() DefinedAction {
	return &LoginUser{}
}

func (action *LoginUser) GetFromJSON(data []byte) {
	err := json.Unmarshal(data, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (action *LoginUser) Process(db *DB, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	db.LoginUser(action, conn, w, req)
}

func (u User) Logout() DefinedAction {
	return &LogoutUser{}
}

func (action *LogoutUser) GetFromJSON(data []byte) {
	err := json.Unmarshal(data, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (action *LogoutUser) Process(db *DB, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	db.LogoutUser(action, conn, w, req)
}
