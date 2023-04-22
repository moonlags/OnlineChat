package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Output struct {
	Action  string        `json:"action"`
	Object  string        `json:"object"`
	Success bool          `json:"success"`
	Status  string        `json:"status"`
	Obj     GeneralObject `json:"obj"`
}

var privkey, _ = rsa.GenerateKey(rand.Reader, 256)
var hmacs = make(map[uint64]string)

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
	str := ""
	var data Output
	action.R.ID = FreeId
	FreeId++
	db.Rooms = append(db.Rooms, action.R)
	for i := range action.R.Users {
		db.Users[db.FindIndexUser(i)].Rooms[action.R.ID] = true
		//rooms, _ := json.Marshal(db.Users[db.FindIndexUser(i)].Rooms)
		//str = str + "UPDATE user SET Rooms=" + string(rooms) + " WHERE ID=" + fmt.Sprint(i) + ";"
	}
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

	}

	//-----------------------------------NOT DONE!-------------------------------------------
	q := `INSERT INTO rooms (Attribute,Name,Messages,ID,Users) VALUES (?,?,?,?,?);` + str

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
			db.Users[db.FindIndexUser(i)].Rooms[action.R.ID] = true
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

	}
}
func (db *DB) ReadRoom(action *ReadRoom, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	var data Output
	if db.FindIndexRoom(action.Data.ID) != -1 {
		data.Action, data.Object, data.Success, data.Status, data.Obj = "read", "room", true, "", db.Rooms[db.FindIndexRoom(action.Data.ID)]
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

		}
	}
	data.Action, data.Object, data.Success, data.Status = "read", "room", false, "Room not found!"
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

	}
}
func (db *DB) DeleteRoom(action *DeleteRoom, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	index := db.FindIndexRoom(action.Data.ID)
	db.Rooms = append(db.Rooms[:index], db.Rooms[index+1:]...)
}
func (db *DB) AddMessage(action *CreateMessage, conn net.Conn, w http.ResponseWriter, req *http.Request) {
	action.M.ID = FreeId
	FreeId++
	var data Output
	if db.FindIndexRoom(action.M.Room) != -1 {
		db.Rooms[db.FindIndexRoom(action.M.Room)].Messages = append(db.Rooms[db.FindIndexRoom(action.M.Room)].Messages, action.M)
		data.Action, data.Object, data.Success, data.Status = "create", "message", true, ""
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

		}
		q := `UPDATE rooms SET Messages=? WHERE ID=?`
		text, err = json.Marshal(db.Rooms[db.FindIndexRoom(action.M.Room)].Messages)
		if err != nil {
			panic(err)
		}
		_, err = db.datab.Query(q, text, action.M.Room)
		if err != nil {
			panic(err)
		}
		return
	}
	data.Action, data.Object, data.Success, data.Status = "create", "message", false, "Room not found!"
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

	}
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
	db.Users[db.FindIndexUser(action.U.ID)].Rooms = make(map[uint64]bool)
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
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"name":     action.U.Name,
			"email":    action.U.Email,
			"password": action.U.Password,
		})
		tokenString, err := token.SignedString([]byte(fmt.Sprint(privkey)))
		if err != nil {
			fmt.Println(err)
			return
		}
		i := strings.Split(tokenString, ".")
		hmacs[action.U.ID] = i[len(i)-1]
		w.Header().Set("jwt", tokenString)
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
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"name":     u.Name,
					"email":    u.Email,
					"password": u.Password,
				})
				tokenString, err := token.SignedString(privkey)
				if err != nil {
					fmt.Println(err)
					return
				}
				i := strings.Split(tokenString, ".")
				hmacs[u.ID] = i[len(i)-1]
				w.Header().Set("jwt", tokenString)
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
		if w != nil {
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
	fmt.Println(req.Header.Get("jwt"))
	if w != nil {
		res := strings.Split(req.Header.Get("jwt"), ".")
		h := hmac.New(sha256.New, []byte(fmt.Sprint(privkey)))
		h.Write([]byte(res[0] + "." + res[1]))
		if !hmac.Equal(h.Sum(nil), []byte(res[len(res)-1])) {
			fmt.Println("Invalid hmac!", privkey)
			return
		}
		w.Header().Set("jwt", req.Header.Get("jwt"))
		io.WriteString(w, string(text))
	}
}
