package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
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
	Process(db *DB, conn net.Conn, w http.ResponseWriter, req *http.Request)
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

	go httphandle()
	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
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
			conn.Close()
			break
		}
		db.UseAction(buf[:n], conn, nil, nil)
		fmt.Println("\nDB after action:")
		for _, p := range db.Users {
			p.Print()
		}
		fmt.Println()
	}
}

func Handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Accept,X-Requested-With, Content-Type, ChatSessionID, Chatsessionid, Access-Control-Request-Method, Access-Control-Request-Headers, X-Auth-Token")
	if req.Method == "POST" {
		data, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			panic(err)
		}
		fmt.Println("Got Connection")
		db.UseAction(data, nil, w, req)
		fmt.Println(sessions)
		fmt.Println()
	}
	if req.Method == "OPTIONS" {
		fmt.Println("Got OPTIONS with header: ")
		for key, value := range req.Header {
			fmt.Println(key, value)
		}
		w.WriteHeader(204)
	}
}

func httphandle() {
	http.HandleFunc("/", Handler)
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}
