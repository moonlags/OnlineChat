package main

import (
	"encoding/json"
	"fmt"
	"image"
	"net"
	"time"
)

type User struct {
	Obj struct {
		Attribute uint32 `json:"attribute"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		ID        uint64 `json:"id"`
		//Avatar    image.Image
		Rooms    map[uint64]bool `json:"rooms"`
		LoggedIn bool
	} `json:"obj"`
}

type Room struct {
	Obj struct {
		Attribute uint32            `json:"attribute"`
		Name      string            `json:"name"`
		Messages  []Message         `json:"messages"`
		ID        uint64            `json:"id"`
		Users     map[uint64]uint32 `json:"users"`
		//BgImg     image.Image
		//Icon      image.Image
	} `json:"obj"`
}
type Content struct {
	Text  string `json:"text"`
	Image image.Image
	//Audio
	//File bytes.Buffer
}

type Message struct {
	Obj struct {
		Attribute uint32  `json:"attribute"`
		Cont      Content `json:"content"`
		SendTime  time.Time
		Author    uint64 `json:"author"`
		Reply     uint64 `json:"reply"`
		Room      uint64 `json:"room"`
		ID        uint64 `json:"id"`
	} `json:"obj"`
}

type Action struct {
	Action string `json:"action"`
	Object string `json:"object"`
	Data   struct {
		ID       uint64            `json:"id"`
		Name     string            `json:"name"`
		Email    string            `json:"email"`
		Password string            `json:"password"`
		Users    map[uint64]uint32 `json:"users"`
		Room     uint64            `json:"room"`
		Author   uint64            `json:"author"`
		Reply    uint64            `json:"reply"`
		Cont     Content           `json:"content"`
	} `json:"data"`
}

type Output struct {
	Action  string `json:"action"`
	Object  string `json:"object"`
	Success bool   `json:"success"`
	Status  string `json:"status"`
}

func FindIndexRoom(id uint64, Rooms []Room) int {
	for i, p := range Rooms {
		if p.Obj.ID == id {
			return i
		}
	}
	return -1
}

func main() {
	var MainUser User
	var Rooms []Room
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	go ReadConnection(conn, &MainUser, &Rooms)
	sendToConnection(conn, &MainUser)
}

func sendToConnection(conn net.Conn, MainUser *User) {
	var d Action
	var s string
	for {
		fmt.Println(*MainUser)
		fmt.Print("Action: ")
		fmt.Scan(&s)
		switch s {
		case "Login":
			if (*MainUser).Obj.ID != 0 {
				fmt.Println("Error: You already in!")
				continue
			}
			d.Action = "login"
			d.Object = "user"
			fmt.Print("Email: ")
			fmt.Scan(&d.Data.Email)
			fmt.Print("Password: ")
			fmt.Scan(&d.Data.Password)
			text, err := json.Marshal(d)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = conn.Write(text)
			if err != nil {
				fmt.Println(err)
				return
			}
		case "Register":
			if (*MainUser).Obj.ID != 0 {
				fmt.Println("Error: You already in!")
				continue
			}
			d.Action = "create"
			d.Object = "user"
			fmt.Print("Name: ")
			fmt.Scan(&d.Data.Name)
			fmt.Print("Email: ")
			fmt.Scan(&d.Data.Email)
			fmt.Print("Password: ")
			fmt.Scan(&d.Data.Password)
			text, err := json.Marshal(d)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = conn.Write(text)
			if err != nil {
				fmt.Println(err)
				return
			}
		case "Logout":
			if (*MainUser).Obj.ID == 0 {
				fmt.Println("Error: You are not in!")
				continue
			}
			d.Action = "logout"
			d.Object = "user"
			d.Data.ID = (*MainUser).Obj.ID
			text, err := json.Marshal(d)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = conn.Write(text)
			if err != nil {
				fmt.Println(err)
				return
			}
		case "CreateRoom":
			if (*MainUser).Obj.ID == 0 {
				fmt.Println("Error: You are not in!")
				continue
			}
			d.Action = "create"
			d.Object = "room"
			fmt.Print("Name: ")
			fmt.Scan(&d.Data.Name)
			d.Data.Users = make(map[uint64]uint32)
			d.Data.Users[(*MainUser).Obj.ID] = 0
			text, err := json.Marshal(d)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = conn.Write(text)
			if err != nil {
				fmt.Println(err)
				return
			}
		case "JoinRoom":
			if (*MainUser).Obj.ID == 0 {
				fmt.Println("Error: You are not in!")
				continue
			}
			d.Action = "update"
			d.Object = "room"
			fmt.Print("Name: ")
			fmt.Scan(&d.Data.Name)
			fmt.Print("ID: ")
			fmt.Scan(&d.Data.ID)
			var i uint64
			for i = range (*MainUser).Obj.Rooms {
				if i == d.Data.ID {
					fmt.Println("Error: You already in this room!")
					break
				}
			}
			if i == d.Data.ID {
				continue
			}
			d.Data.Users = make(map[uint64]uint32)
			d.Data.Users[(*MainUser).Obj.ID] = 0
			text, err := json.Marshal(d)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = conn.Write(text)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func ReadConnection(conn net.Conn, MainUser *User, Rooms *[]Room) {
	buf := make([]byte, 1000)
	var data Output
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println(string(buf[:n]))
		err = json.Unmarshal(buf[:n], &data)
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println(data)
		fmt.Println(string(buf[:n]))
		if data.Success && data.Status == "" {
			if (data.Action == "login" || data.Action == "create" || data.Action == "update") && data.Object == "user" {
				err := json.Unmarshal(buf[:n], MainUser)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			if data.Action == "logout" && data.Object == "user" {
				*MainUser = User{}
			}
			if (data.Action == "create" || data.Action == "update") && data.Object == "room" {
				var r Room
				err := json.Unmarshal(buf[:n], &r)
				if err != nil {
					fmt.Println(err)
					return
				}
				(*Rooms) = append(*Rooms, r)
				//fmt.Println(*Rooms)
				temp := MainUser.Obj.Rooms
				MainUser.Obj.Rooms = make(map[uint64]bool)
				(*MainUser).Obj.Rooms[r.Obj.ID] = true
				for i, v := range temp {
					(*MainUser).Obj.Rooms[i] = v
				}
			}
		} else {
			fmt.Println("\nError: ", data.Status)
		}
	}
}
