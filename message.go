package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

type Content struct {
	Text  string `json:"text"`
	Image image.Image
	//Audio
	File bytes.Buffer
}

type Message struct {
	Attribute uint32  `json:"attribute"`
	Cont      Content `json:"content"`
	SendTime  time.Time
	Author    uint64 `json:"author"`
	Reply     uint64 `json:"reply"`
	Room      uint64 `json:"room"`
	ID        uint64 `json:"id"`
}

type CreateMessage struct {
	M Message `json:"data"`
}

type UpdateMessage struct {
	M Message `json:"data"`
}

type ReadMessage struct {
	Data struct {
		ID   uint64 `json:"id"`
		Room uint64 `json:"room"`
	} `json:"data"`
}

type DeleteMessage struct {
	Data struct {
		ID   uint64 `json:"id"`
		Room uint64 `json:"room"`
	} `json:"data"`
}
type LoginMessage struct{}

type LogoutMessage struct{}

func (m Message) Create() DefinedAction {
	return &CreateMessage{}
}

func (action *CreateMessage) GetFromJSON(data []byte) {
	err := json.Unmarshal(data, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (action *CreateMessage) Process(db *DB, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	(*db).AddMessage(action, w, c, req)
}

func (m Message) Update() DefinedAction {
	return &UpdateMessage{}
}

func (action *UpdateMessage) GetFromJSON(data []byte) {
	err := json.Unmarshal(data, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (action *UpdateMessage) Process(db *DB, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	(*db).UpdateMessage(action, w, c, req)
}

func (m Message) Read() DefinedAction {
	return &ReadMessage{}
}

func (action *ReadMessage) GetFromJSON(data []byte) {
	err := json.Unmarshal(data, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (action *ReadMessage) Process(db *DB, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	(*db).ReadMessage(action, w, c, req)
}

func (m Message) Delete() DefinedAction {
	return &DeleteMessage{}
}

func (action *DeleteMessage) GetFromJSON(data []byte) {
	err := json.Unmarshal(data, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (action *DeleteMessage) Process(db *DB, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	(*db).DeleteMessage(action, w, c, req)
}

func (m Message) Print() {
	fmt.Printf(" text:%s Attribute:%d Author:%d Reply:%d Room:%d", m.Cont.Text, m.Attribute, m.Author, m.Reply, m.Room)
}
func (m Message) GetId() uint64 {
	return m.ID
}

func (m Message) Login() DefinedAction {
	return &LoginMessage{}
}

func (action *LoginMessage) GetFromJSON(data []byte) {}

func (action *LoginMessage) Process(db *DB, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	fmt.Println("Message cant Login")
}
func (m Message) Logout() DefinedAction {
	return &LogoutMessage{}
}

func (action *LogoutMessage) GetFromJSON(data []byte) {}

func (action *LogoutMessage) Process(db *DB, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	fmt.Println("Message cant Logout")
}
