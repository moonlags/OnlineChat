package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"nhooyr.io/websocket"
)

const (
	AttrRead uint32 = 1 << iota
	AttrWrite
	AttrDelete
	AttrDeleted
	AttrVisible
)

//var attr uint32
//attr = AttrWrite|AttrRead
//if attr & AttrWrite>0

type Action struct {
	Action  string `json:"action"`
	ObjName string `json:"object"`
	UserID  uint64 `json:"userid"`
	Jwt     string `json:"jwt"`
}

type DefinedAction interface {
	GetFromJSON([]byte)
	Process(db *DB, w http.ResponseWriter, c *websocket.Conn, req *http.Request)
}

type GeneralObject interface { // room message user
	Create() DefinedAction
	Read() DefinedAction
	Update() DefinedAction
	Delete() DefinedAction
	Login() DefinedAction
	Logout() DefinedAction
	GetId() uint64
	Print()
}

var db DB

func main() {
	fmt.Println(time.Now())

	cfg := mysql.NewConfig()
	(*cfg).Addr = "localhost"
	(*cfg).User = "root"
	(*cfg).Passwd = "masterkey"
	(*cfg).Net = "tcp"
	(*cfg).DBName = "chatdb"

	var err error
	db.datab, err = sql.Open("mysql", cfg.FormatDSN())
	//SELECT Attribute,NAME,Email,PASSWORD,ID,Rooms,LoggedIn FROM users ORDER BY ID DESC
	if err != nil {
		fmt.Println(err)
		return
	}

	q := `SELECT Attribute,Name,Email,Password,ID,Rooms FROM users ORDER BY ID DESC`
	rows, err := db.datab.Query(q)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var data User
		var dataRooms []byte
		if err := rows.Scan(&data.Attribute, &data.Name, &data.Email, &data.Password, &data.ID, &dataRooms); err != nil {
			panic(err)
		}
		if data.ID >= FreeId {
			FreeId = data.ID + 1
		}
		err := json.Unmarshal(dataRooms, &data.Rooms)
		if err != nil {
			panic(err)
		}
		if data.Rooms == nil {
			data.Rooms = make(map[uint64]bool)
		}
		db.Users = append(db.Users, data)
	}
	q = `SELECT Attribute,Name,Messages,ID,Users FROM rooms ORDER BY ID DESC`
	rows, err = db.datab.Query(q)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var data Room
		var dataMessages []byte
		var dataUsers []byte
		if err := rows.Scan(&data.Attribute, &data.Name, &dataMessages, &data.ID, &dataUsers); err != nil {
			panic(err)
		}
		if data.ID >= FreeId {
			FreeId = data.ID + 1
		}
		err := json.Unmarshal(dataMessages, &data.Messages)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(dataUsers, &data.Users)
		if err != nil {
			panic(err)
		}
		if data.Users == nil {
			data.Users = make(map[uint64]uint32)
		}
		db.Rooms = append(db.Rooms, data)
	}

	http.HandleFunc("/ws", HandlerWS)

	err = http.ListenAndServe(":8080", nil)
	panic(err)
}

func HandlerWS(w http.ResponseWriter, req *http.Request) {
	c, err := websocket.Accept(w, req, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:3000"},
	})
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "")

	ctx := context.Background()
	fmt.Println("websocket accepted")

	for {
		msgType, data, err := c.Read(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("MsgType: %d, Data: %s\n", msgType, string(data))
		db.UseAction(data, w, c, req)
	}
}
