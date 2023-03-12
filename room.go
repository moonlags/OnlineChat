package main

import (
	"encoding/json"
	"fmt"
	"image"
	"net"
)

type Room struct {
	Attribute uint32            `json:"attribute"`
	Name      string            `json:"name"`
	Messages  []Message         `json:"messages"`
	ID        uint64            `json:"id"`
	Users     map[uint64]uint32 `json:"users"`
	BgImg     image.Image
	Icon      image.Image
}

func (r Room) FindIndex(id uint64) int {
	for i, p := range r.Messages {
		if p.GetId() == id {
			return i
		}
	}
	return -1
}

type CreateRoom struct {
	R Room `json:"data"`
}

type UpdateRoom struct {
	R Room `json:"data"`
}

type ReadRoom struct {
	Data struct {
		ID uint64 `json:"id"`
	} `json:"data"`
}

type DeleteRoom struct {
	Data struct {
		ID uint64 `json:"id"`
	} `json:"data"`
}
type LoginRoom struct{}

type LogoutRoom struct{}

func (r Room) Create() DefinedAction {
	return &CreateRoom{}
}

func (action *CreateRoom) GetFromJSON(data []byte) {
	err := json.Unmarshal(data, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (action *CreateRoom) Process(db *DB, conn net.Conn) {
	(*db).AddRoom(action, conn)
}

func (r Room) Update() DefinedAction {
	return &UpdateRoom{}
}

func (action *UpdateRoom) GetFromJSON(data []byte) {
	err := json.Unmarshal(data, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (action *UpdateRoom) Process(db *DB, conn net.Conn) {
	(*db).UpdateRoom(action, conn)
}

func (r Room) Read() DefinedAction {
	return &ReadRoom{}
}

func (action *ReadRoom) GetFromJSON(data []byte) {
	err := json.Unmarshal(data, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (action *ReadRoom) Process(db *DB, conn net.Conn) {
	(*db).ReadRoom(action, conn)
}

func (r Room) Delete() DefinedAction {
	return &DeleteRoom{}
}

func (action *DeleteRoom) GetFromJSON(data []byte) {
	err := json.Unmarshal(data, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (action *DeleteRoom) Process(db *DB, conn net.Conn) {
	(*db).DeleteRoom(action, conn)
}

func (r Room) Print() {
	fmt.Printf("ID:%d Name:%s Attribute:%d", r.ID, r.Name, r.Attribute)
}
func (r Room) GetId() uint64 {
	return r.ID
}

func (r Room) Login() DefinedAction {
	return &LoginRoom{}
}

func (action *LoginRoom) GetFromJSON(data []byte) {}

func (action *LoginRoom) Process(db *DB, conn net.Conn) { fmt.Println("Room cant Login") }
func (r Room) Logout() DefinedAction {
	return &LogoutRoom{}
}

func (action *LogoutRoom) GetFromJSON(data []byte) {}

func (action *LogoutRoom) Process(db *DB, conn net.Conn) { fmt.Println("Room cant Logout") }
