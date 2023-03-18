package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
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
		Rooms map[uint64]bool `json:"rooms"`
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
	client := &http.Client{}

	sendToConnection(client, &MainUser, &Rooms)
}

func sendToConnection(client *http.Client, MainUser *User, Rooms *[]Room) {
	var d Action
	var s string
	var session string
	for {
		fmt.Println(*MainUser)
		fmt.Print("Action: ")
		fmt.Scan(&s)
		switch s {
		case "Login":
			if (*MainUser).Obj.ID != 0 || session != "" {
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
			var body bytes.Buffer
			body.Write(text)
			req, err := http.NewRequest("POST", "http://localhost:8080/", &body)
			if err != nil {
				panic(err)
			}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			resp.Body.Close()
			session = resp.Header.Get("ChatSessionID")
			var data2 Output
			err = json.Unmarshal(data, &data2)
			if err != nil {
				fmt.Println(err)
				return
			}
			if !data2.Success {
				fmt.Println("\nError: ", data2.Status)
				continue
			}
			err = json.Unmarshal(data, MainUser)
			if err != nil {
				fmt.Println(err)
				return
			}
		case "Register":
			if (*MainUser).Obj.ID != 0 || session != "" {
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
			var body bytes.Buffer
			body.Write(text)
			req, err := http.NewRequest("POST", "http://localhost:8080/", &body)
			if err != nil {
				panic(err)
			}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			resp.Body.Close()
			session = resp.Header.Get("ChatSessionID")
			var data2 Output
			err = json.Unmarshal(data, &data2)
			if err != nil {
				fmt.Println(err)
				return
			}
			if !data2.Success {
				fmt.Println("\nError: ", data2.Status)
				continue
			}
			err = json.Unmarshal(data, MainUser)
			if err != nil {
				fmt.Println(err)
				return
			}
		case "Logout":
			if (*MainUser).Obj.ID == 0 || session == "" {
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
			var body bytes.Buffer
			body.Write(text)
			req, err := http.NewRequest("POST", "http://localhost:8080/", &body)
			if err != nil {
				panic(err)
			}
			req.Header.Set("ChatSessionID", session)
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			resp.Body.Close()
			var data2 Output
			err = json.Unmarshal(data, &data2)
			if err != nil {
				fmt.Println(err)
				return
			}
			if !data2.Success {
				fmt.Println("\nError: ", data2.Status)
				continue
			}
			*MainUser = User{}
			session = ""
			*Rooms = make([]Room, 0)
		case "CreateRoom":
			if (*MainUser).Obj.ID == 0 || session == "" {
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
			var body bytes.Buffer
			body.Write(text)
			req, err := http.NewRequest("POST", "http://localhost:8080/", &body)
			if err != nil {
				panic(err)
			}
			req.Header.Set("ChatSessionID", session)
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			resp.Body.Close()
			var data2 Output
			err = json.Unmarshal(data, &data2)
			if err != nil {
				fmt.Println(err)
				return
			}
			if !data2.Success {
				fmt.Println("\nError: ", data2.Status)
				continue
			}
			var r Room
			err = json.Unmarshal(data, &r)
			if err != nil {
				fmt.Println(err)
				return
			}
			(*Rooms) = append(*Rooms, r)
			temp := MainUser.Obj.Rooms
			MainUser.Obj.Rooms = make(map[uint64]bool)
			(*MainUser).Obj.Rooms[r.Obj.ID] = true
			for i, v := range temp {
				(*MainUser).Obj.Rooms[i] = v
			}
		case "JoinRoom":
			if (*MainUser).Obj.ID == 0 || session == "" {
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
			var body bytes.Buffer
			body.Write(text)
			req, err := http.NewRequest("POST", "http://localhost:8080/", &body)
			if err != nil {
				panic(err)
			}
			req.Header.Set("ChatSessionID", session)
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			resp.Body.Close()
			var data2 Output
			err = json.Unmarshal(data, &data2)
			if err != nil {
				fmt.Println(err)
				return
			}
			if !data2.Success {
				fmt.Println("\nError: ", data2.Status)
				continue
			}
			var r Room
			err = json.Unmarshal(data, &r)
			if err != nil {
				fmt.Println(err)
				return
			}
			(*Rooms) = append(*Rooms, r)
			temp := MainUser.Obj.Rooms
			MainUser.Obj.Rooms = make(map[uint64]bool)
			(*MainUser).Obj.Rooms[r.Obj.ID] = true
			for i, v := range temp {
				(*MainUser).Obj.Rooms[i] = v
			}
		}
	}
}
