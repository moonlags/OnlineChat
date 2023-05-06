package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"nhooyr.io/websocket"
)

type Output struct {
	Action  string        `json:"action"`
	Object  string        `json:"object"`
	Success bool          `json:"success"`
	Status  string        `json:"status"`
	Jwt     string        `json:"jwt"`
	Obj     GeneralObject `json:"obj"`
}

var privkey = rand.Intn(999999999999)
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

func (db *DB) UseAction(text []byte, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {

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
	if act.ObjName != "user" || act.Action != "login" && act.Action != "create" {
		fmt.Println(act.Jwt)
		token, err := jwt.Parse(act.Jwt, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(fmt.Sprint(privkey)), nil
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["name"], claims["email"], claims["password"])
		} else {
			fmt.Println(err)
		}
	}
	toDo.GetFromJSON(text)
	//fmt.Println(sessions)
	toDo.Process(db, w, c, req)
}

func (db *DB) AddRoom(action *CreateRoom, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	var data Output
	action.R.ID = FreeId
	FreeId++
	db.Rooms = append(db.Rooms, action.R)
	for i := range action.R.Users {
		db.Users[db.FindIndexUser(i)].Rooms[action.R.ID] = true
	}
	data.Action, data.Object, data.Success, data.Status, data.Obj = "create", "room", true, "", db.Rooms[db.FindIndexRoom(action.R.ID)]
	text, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	err = c.Write(ctx, 1, text)
	if err != nil {
		panic(err)
	}

	//-----------------------------------NOT DONE!-------------------------------------------
	q := `INSERT INTO rooms (Attribute,Name,Messages,ID,Users) VALUES (?,?,?,?,?)`
	q2 := `UPDATE users SET Rooms=? WHERE ID=?`

	temp := db.Rooms[db.FindIndexRoom(action.R.ID)]
	text, err = json.Marshal(temp.Messages)
	if err != nil {
		panic(err)
	}
	text2, err := json.Marshal(temp.Users)
	if err != nil {
		panic(err)
	}
	var text3 []byte
	var i uint64
	for i = range action.R.Users {
		text3, err = json.Marshal(db.Users[db.FindIndexUser(i)].Rooms)
		if err != nil {
			panic(err)
		}
	}
	_, err = db.datab.Query(q, temp.Attribute, temp.Name, text, temp.ID, text2)
	if err != nil {
		panic(err)
	}
	_, err = db.datab.Query(q2, text3, i)
	if err != nil {
		panic(err)
	}
}
func (db *DB) UpdateRoom(action *UpdateRoom, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
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

		ctx := context.Background()
		err = c.Write(ctx, 1, text)
		if err != nil {
			panic(err)
		}
		q := `UPDATE rooms SET Attribute=?, Name=?,Messages=?,Users=? WHERE ID=?;UPDATE users SET Rooms=? WHERE ID=?`
		text, err = json.Marshal(db.Rooms[db.FindIndexRoom(action.R.ID)].Messages)
		if err != nil {
			panic(err)
		}
		text2, err := json.Marshal(db.Rooms[db.FindIndexRoom(action.R.ID)].Users)
		if err != nil {
			panic(err)
		}
		var text3 []byte
		var i uint64
		for i = range action.R.Users {
			text3, err = json.Marshal(db.Users[db.FindIndexUser(i)].Rooms)
			if err != nil {
				panic(err)
			}
		}
		_, err = db.datab.Query(q, temp.Attribute, temp.Name, text, text2, temp.ID, text3, i)
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
	ctx := context.Background()
	err = c.Write(ctx, 1, text)
	if err != nil {
		panic(err)
	}
}
func (db *DB) ReadRoom(action *ReadRoom, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	var data Output
	if db.FindIndexRoom(action.Data.ID) != -1 {
		data.Action, data.Object, data.Success, data.Status, data.Obj = "read", "room", true, "", db.Rooms[db.FindIndexRoom(action.Data.ID)]
		text, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		ctx := context.Background()
		err = c.Write(ctx, 1, text)
		if err != nil {
			panic(err)
		}
	}
	data.Action, data.Object, data.Success, data.Status = "read", "room", false, "Room not found!"
	text, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := context.Background()
	err = c.Write(ctx, 1, text)
	if err != nil {
		panic(err)
	}
}
func (db *DB) DeleteRoom(action *DeleteRoom, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	index := db.FindIndexRoom(action.Data.ID)
	db.Rooms = append(db.Rooms[:index], db.Rooms[index+1:]...)
}
func (db *DB) AddMessage(action *CreateMessage, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
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

		ctx := context.Background()
		err = c.Write(ctx, 1, text)
		if err != nil {
			panic(err)
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

	ctx := context.Background()
	err = c.Write(ctx, 1, text)
	if err != nil {
		panic(err)
	}
}
func (db *DB) UpdateMessage(action *UpdateMessage, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	db.Rooms[db.FindIndexRoom(action.M.Room)].Messages[db.Rooms[db.FindIndexRoom(action.M.Room)].FindIndex(action.M.ID)] = action.M
}
func (db *DB) ReadMessage(action *ReadMessage, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	db.Rooms[db.FindIndexRoom(action.Data.Room)].Messages[db.Rooms[db.FindIndexRoom(action.Data.Room)].FindIndex(action.Data.ID)].Print()
}
func (db *DB) DeleteMessage(action *DeleteMessage, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	index := db.Rooms[db.FindIndexRoom(action.Data.Room)].FindIndex(action.Data.ID)
	db.Rooms[db.FindIndexRoom(action.Data.Room)].Messages = append(db.Rooms[db.FindIndexRoom(action.Data.Room)].Messages[:index], db.Rooms[db.FindIndexRoom(action.Data.Room)].Messages[index+1:]...)
}
func (db *DB) AddUser(action *CreateUser, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
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
			ctx := context.Background()
			err = c.Write(ctx, 1, text)
			if err != nil {
				panic(err)
			}
			return
		}
	}
	action.U.ID = FreeId
	FreeId++
	hash, err := argon2id.CreateHash(action.U.Password, argon2id.DefaultParams)
	if err != nil {
		panic(err)
	}
	action.U.Password = hash
	db.Users = append(db.Users, action.U)
	db.Users[db.FindIndexUser(action.U.ID)].Rooms = make(map[uint64]bool)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":     action.U.Name,
		"email":    action.U.Email,
		"password": action.U.Password,
	})
	tokenString, err := token.SignedString([]byte(fmt.Sprint(privkey)))
	if err != nil {
		panic(err)
	}
	fmt.Println(tokenString)
	hmacs[action.U.ID] = tokenString

	data.Action, data.Object, data.Success, data.Status, data.Jwt, data.Obj = "create", "user", true, "", hmacs[action.U.ID], db.Users[db.FindIndexUser(action.U.ID)]
	text, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = c.Write(ctx, 1, text)
	if err != nil {
		panic(err)
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
func (db *DB) UpdateUser(action *UpdateUser, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	db.Users[db.FindIndexUser(action.U.ID)] = action.U
}
func (db *DB) ReadUser(action *ReadUser, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	db.Users[db.FindIndexUser(action.Data.ID)].Print()
}
func (db *DB) DeleteUser(action *DeleteUser, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	index := db.FindIndexUser(action.Data.ID)
	db.Users = append(db.Users[:index], db.Users[index+1:]...)
}

func (db *DB) LoginUser(action *LoginUser, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	var data Output
	for _, u := range db.Users {
		fmt.Println(action.Data.Email, action.Data.Password)
		match, err := argon2id.ComparePasswordAndHash(action.Data.Password, u.Password)
		if err != nil {
			panic(err)
		}
		if u.Email == action.Data.Email && match {

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"name":     u.Name,
				"email":    u.Email,
				"password": u.Password,
			})
			tokenString, err := token.SignedString([]byte(fmt.Sprint(privkey)))
			if err != nil {
				panic(err)
			}
			hmacs[u.ID] = tokenString

			data.Action, data.Object, data.Success, data.Status, data.Jwt, data.Obj = "login", "user", true, "", hmacs[u.ID], db.Users[db.FindIndexUser(u.ID)]
			text, err := json.Marshal(data)
			if err != nil {
				panic(err)
			}

			ctx := context.Background()
			err = c.Write(ctx, 1, text)
			if err != nil {
				panic(err)
			}
			return
		}
	}
	data.Action, data.Object, data.Success, data.Status = "login", "user", false, "User not found!"
	text, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = c.Write(ctx, 1, text)
	if err != nil {
		panic(err)
	}
}

func (db *DB) LogoutUser(action *LogoutUser, w http.ResponseWriter, c *websocket.Conn, req *http.Request) {
	var data Output
	index := db.FindIndexUser(action.Data.ID)
	if index == -1 {
		data.Action, data.Object, data.Success, data.Status = "logout", "user", false, "User Not Found, wrong ID!"
		text, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set("Chatsessionid", req.Header.Get("Chatsessionid"))
		ctx := context.Background()
		err = c.Write(ctx, 1, text)
		if err != nil {
			panic(err)
		}
		return
	}
	data.Action, data.Object, data.Success, data.Status = "logout", "user", true, ""
	text, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	err = c.Write(ctx, 1, text)
	if err != nil {
		panic(err)
	}
}
