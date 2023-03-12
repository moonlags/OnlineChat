package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/go-sql-driver/mysql"
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
}

type DefinedAction interface {
	GetFromJSON([]byte)
	Process(db *DB, conn net.Conn)
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

	q := `SELECT Attribute,Name,Email,Password,ID,Rooms,LoggedIn FROM users ORDER BY ID DESC`
	rows, err := db.datab.Query(q)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var data User
		var dataRooms []byte
		var a []uint8
		if err := rows.Scan(&data.Attribute, &data.Name, &data.Email, &data.Password, &data.ID, &dataRooms, &a); err != nil {
			panic(err)
		}
		if a[0] == 1 {
			data.LoggedIn = true
		} else {
			data.LoggedIn = false
		}
		if data.ID >= FreeId {
			FreeId = data.ID + 1
		}
		err := json.Unmarshal(dataRooms, &data.Rooms)
		if err != nil {
			panic(err)
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
		db.Rooms = append(db.Rooms, data)
	}

	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
		defer conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())
		go handleConnection(conn)
	}

}
func handleConnection(conn net.Conn) {
	buf := make([]byte, 1000)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		db.UseAction(buf[:n], conn)
		fmt.Println("\nDB after action:")
		for _, p := range db.Users {
			p.Print()
		}
		fmt.Println()
	}
}