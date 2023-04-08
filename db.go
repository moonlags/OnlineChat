package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"
)

type Output struct {
	Action  string        `json:"action"`
	Object  string        `json:"object"`
	Success bool          `json:"success"`
	Status  string        `json:"status"`
	Obj     GeneralObject `json:"obj"`
}

type Session struct {
	UserID  uint64
	Timeout time.Time
}

var sessions = make(map[string]Session)
var FreeId uint64 = 1

type DB struct {
	datab *sql.DB
	Users []User
	Rooms []Room
}

func (db *DB) FindIndexUser(id uint64) int {
	for i, p := range db.Users {
		if p.GetId() == id {
			return i
		}
	}
	return -1
}
func (db *DB) FindIndexRoom(id uint64) int {
	for i, p := range db.Rooms {
		if p.GetId() == id {
			return i
		}
	}
	return -1
}

func (db *DB) UseAction(text []byte, conn net.Conn, w http.ResponseWriter, req *http.Request) {

	var act Action

	err := json.Unmarshal(text, &act)
	if err != nil {
		fmt.Println(err)
		return
	}

	var obj GeneralObject
	switch act.ObjName {
	case "room":
		obj = &Room{}
	case "user":
		obj = &User{}
	case "message":
		obj = &Message{}
	default:
		fmt.Println("Object not found", act.ObjName)
		return
	}
	var toDo DefinedAction
	switch act.Action {
	case "create":
		toDo = obj.Create()
	case "update":
		toDo = obj.Update()
	case "read":
		toDo = obj.Read()
	case "delete":
		toDo = obj.Delete()
	case "login":
		toDo = obj.Login()
	case "logout":
		toDo = obj.Logout()
	default:
		fmt.Println("unknown action", act.Action)
		return
	}
	toDo.GetFromJSON(text)
	//fmt.Println(sessions)
	toDo.Process(db, conn, w, req)
}

func (db *DB) AddRoom(action *CreateRoom, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	var data Output
	action.R.ID = FreeId
	FreeId++
	db.Rooms = append(db.Rooms, action.R)
	data.Action, data.Object, data.Success, data.Status, data.Obj = "create", "room", true, "", db.Rooms[db.FindIndexRoom(action.R.ID)]
	text, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	if conn != nil {
		_, err = conn.Write(text)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if w != nil {
		for i := range db.Rooms[db.FindIndexRoom(action.R.ID)].Users {
			if sessions[req.Header.Get("ChatSessionID")].UserID == i {
				w.Header().Set("ChatSessionID", req.Header.Get("ChatSessionID"))
				io.WriteString(w, string(text))
			}
		}
	}
	//-----------------------------------NOT DONE!-------------------------------------------
	q := `INSERT INTO rooms (Attribute,Name,Messages,ID,Users) VALUES (?,?,?,?,?); UPDATE users SET Rooms=? WHERE ID=?`

	temp := db.Rooms[db.FindIndexRoom(action.R.ID)]
	text, err = json.Marshal(temp.Messages)
	if err != nil {
		panic(err)
	}
	text2, err := json.Marshal(temp.Users)
	if err != nil {
		panic(err)
	}
	_, err = db.datab.Query(q, temp.Attribute, temp.Name, text, temp.ID, text2)
	if err != nil {
		panic(err)
	}
}
func (db *DB) UpdateRoom(action *UpdateRoom, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	var data Output
	if db.Rooms[db.FindIndexRoom(action.R.ID)].Name == action.R.Name {
		temp := db.Rooms[db.FindIndexRoom(action.R.ID)]
		db.Rooms[db.FindIndexRoom(action.R.ID)] = action.R
		for i, v := range temp.Users {
			db.Rooms[db.FindIndexRoom(action.R.ID)].Users[i] = v
		}
		for i, v := range temp.Messages {
			db.Rooms[db.FindIndexRoom(action.R.ID)].Messages[i] = v
		}
		data.Action, data.Object, data.Success, data.Status, data.Obj = "update", "room", true, "", db.Rooms[db.FindIndexRoom(action.R.ID)]
		text, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		if conn != nil {
			_, err = conn.Write(text)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if w != nil {
			for i := range db.Rooms[db.FindIndexRoom(action.R.ID)].Users {
				if sessions[req.Header.Get("ChatSessionID")].UserID == i {
					w.Header().Set("ChatSessionID", req.Header.Get("ChatSessionID"))
					io.WriteString(w, string(text))
				}
			}
		}
		//----------------------NOT DONE!-------------------------------------------
		//UPDATE `chatdb`.`rooms` SET `Name`='fsd' WHERE  `Attribute`=111
		q := `UPDATE rooms SET Attribute=?, Name=?,Messages=?,Users=? WHERE ID=?`
		text, err = json.Marshal(db.Rooms[db.FindIndexRoom(action.R.ID)].Messages)
		if err != nil {
			panic(err)
		}
		text2, err := json.Marshal(db.Rooms[db.FindIndexRoom(action.R.ID)].Users)
		if err != nil {
			panic(err)
		}
		_, err = db.datab.Query(q, temp.Attribute, temp.Name, text, text2, temp.ID)
		if err != nil {
			panic(err)
		}
		return
	}
	data.Action, data.Object, data.Success, data.Status = "update", "room", false, "Room not found!"
	text, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	if conn != nil {
		_, err = conn.Write(text)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if w != nil {
		for i := range db.Rooms[db.FindIndexRoom(action.R.ID)].Users {
			if sessions[req.Header.Get("ChatSessionID")].UserID == i {
				w.Header().Set("ChatSessionID", req.Header.Get("ChatSessionID"))
				io.WriteString(w, string(text))
			}
		}
	}
}
func (db *DB) ReadRoom(action *ReadRoom, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	db.Rooms[db.FindIndexRoom(action.Data.ID)].Print()
}
func (db *DB) DeleteRoom(action *DeleteRoom, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	index := db.FindIndexRoom(action.Data.ID)
	db.Rooms = append(db.Rooms[:index], db.Rooms[index+1:]...)
}
func (db *DB) AddMessage(action *CreateMessage, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	action.M.ID = FreeId
	FreeId++
	roomIndex := db.FindIndexRoom(action.M.Room)
	db.Rooms[roomIndex].Messages = append(db.Rooms[roomIndex].Messages, action.M)
}
func (db *DB) UpdateMessage(action *UpdateMessage, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	db.Rooms[db.FindIndexRoom(action.M.Room)].Messages[db.Rooms[db.FindIndexRoom(action.M.Room)].FindIndex(action.M.ID)] = action.M
}
func (db *DB) ReadMessage(action *ReadMessage, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	db.Rooms[db.FindIndexRoom(action.Data.Room)].Messages[db.Rooms[db.FindIndexRoom(action.Data.Room)].FindIndex(action.Data.ID)].Print()
}
func (db *DB) DeleteMessage(action *DeleteMessage, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	index := db.Rooms[db.FindIndexRoom(action.Data.Room)].FindIndex(action.Data.ID)
	db.Rooms[db.FindIndexRoom(action.Data.Room)].Messages = append(db.Rooms[db.FindIndexRoom(action.Data.Room)].Messages[:index], db.Rooms[db.FindIndexRoom(action.Data.Room)].Messages[index+1:]...)
}
func (db *DB) AddUser(action *CreateUser, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	var data Output
	//fmt.Println(action)
	for _, u := range db.Users {
		if u.Email == action.U.Email {
			data.Action, data.Object, data.Success, data.Status = "create", "user", false, "User already exists!"
			text, err := json.Marshal(data)
			if err != nil {
				fmt.Println(err)
				return
			}
			if conn != nil {
				_, err = conn.Write(text)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			if w != nil {
				io.WriteString(w, string(text))
			}
			return
		}
	}
	//INSERT INTO users (Attribute,NAME,Email,PASSWORD,ID) VALUES (0,"saf","ferg","asffa",1);
	action.U.ID = FreeId
	FreeId++
	db.Users = append(db.Users, action.U)
	data.Action, data.Object, data.Success, data.Status, data.Obj = "create", "user", true, "", db.Users[db.FindIndexUser(action.U.ID)]
	text, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	if conn != nil {
		_, err = conn.Write(text)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if w != nil {
		i := fmt.Sprint(rand.Intn(9999999999))
		for _, ok := sessions[i]; ok; {
			i = fmt.Sprint(rand.Intn(9999999999))
		}
		sessions[i] = Session{action.U.ID, time.Now()}
		w.Header().Set("ChatSessionID", i)
		io.WriteString(w, string(text))
	}
	q := `INSERT INTO users (Attribute,Name,Email,Password,ID,Rooms) VALUES (?,?,?,?,?,?)`
	temp := db.Users[db.FindIndexUser(action.U.ID)]
	text, err = json.Marshal(temp.Rooms)
	if err != nil {
		panic(err)
	}
	_, err = db.datab.Query(q, temp.Attribute, temp.Name, temp.Email, temp.Password, temp.ID, text)
	if err != nil {
		panic(err)
	}
}
func (db *DB) UpdateUser(action *UpdateUser, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	db.Users[db.FindIndexUser(action.U.ID)] = action.U
}
func (db *DB) ReadUser(action *ReadUser, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	db.Users[db.FindIndexUser(action.Data.ID)].Print()
}
func (db *DB) DeleteUser(action *DeleteUser, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	index := db.FindIndexUser(action.Data.ID)
	db.Users = append(db.Users[:index], db.Users[index+1:]...)
}

func (db *DB) LoginUser(action *LoginUser, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	var data Output
	for _, u := range db.Users {
		fmt.Println(action.Data.Email, action.Data.Password)
		if u.Email == action.Data.Email && u.Password == action.Data.Password {
			data.Action, data.Object, data.Success, data.Status, data.Obj = "login", "user", true, "", db.Users[db.FindIndexUser(u.ID)]
			text, err := json.Marshal(data)
			if err != nil {
				fmt.Println(err)
				return
			}
			if conn != nil {
				_, err = conn.Write(text)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			if w != nil {
				i := fmt.Sprint(rand.Intn(9999999999))
				for _, ok := sessions[i]; ok; {
					i = fmt.Sprint(rand.Intn(9999999999))
				}
				sessions[i] = Session{data.Obj.GetId(), time.Now()}
				w.Header().Set("ChatSessionID", i)
				io.WriteString(w, string(text))
			}
			return
		}
	}
	data.Action, data.Object, data.Success, data.Status = "login", "user", false, "User not found!"
	text, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	if conn != nil {
		_, err = conn.Write(text)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if w != nil {
		io.WriteString(w, string(text))
	}
}

func (db *DB) LogoutUser(action *LogoutUser, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	var data Output
	index := db.FindIndexUser(action.Data.ID)
	if index == -1 {
		data.Action, data.Object, data.Success, data.Status = "logout", "user", false, "User Not Found, wrong ID!"
		text, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		if conn != nil {
			_, err = conn.Write(text)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if w != nil && sessions[req.Header.Get("Chatsessionid")].UserID == action.Data.ID {
			w.Header().Set("Chatsessionid", req.Header.Get("Chatsessionid"))
			io.WriteString(w, string(text))
		}
		return
	}
	data.Action, data.Object, data.Success, data.Status = "logout", "user", true, ""
	text, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	if conn != nil {
		_, err = conn.Write(text)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println(sessions[req.Header.Get("Chatsessionid")].UserID)
	if w != nil && sessions[req.Header.Get("Chatsessionid")].UserID == action.Data.ID {
		delete(sessions, req.Header.Get("Chatsessionid"))
		w.Header().Set("Chatsessionid", req.Header.Get("Chatsessionid"))
		io.WriteString(w, string(text))
	}
}
