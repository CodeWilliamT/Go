package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// New User and its listener loop
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	//start user listener
	go user.ListenMessage()
	return user
}

func (this *User) Online() {

	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "Online")

}

func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "Offline")
}

func (this *User) GetMessage(msg string) {
	this.conn.Write([]byte(msg + "\n"))
}

func (this *User) SendMessage(msg string) {
	if msg == "-userlist" {
		//list users
		for _, user := range this.server.OnlineMap {
			userinfo := "[" + user.Addr + "]" + user.Name + ":" + " Online"
			this.GetMessage(userinfo)
		}
		return
	}
	if len(msg) > 8 && msg[:8] == "-rename|" {
		//Format:rename|Tony
		// newName := strings.Split(msg, "|")[1]
		newName := msg[8:]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.GetMessage("Name has been used")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.GetMessage("Your name has been changed to " + this.Name)
		}
		return
	}
	if len(msg) > 5 && msg[:4] == "-to|" {
		targetName := strings.Split(msg, "|")[1]
		if targetName == "" {
			this.GetMessage("Format error. Please user format like \"-to|Tony|hello\"")
			return
		}
		targetUser, ok := this.server.OnlineMap[targetName]
		if !ok {
			this.GetMessage("Target User is not online")
			return
		}

		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.GetMessage("Don't send empty message")
			return
		}
		targetUser.GetMessage("[" + this.Addr + "]" + this.Name + " chat to you: \n" + content)
		return
	}
	this.server.BroadCast(this, msg)
}

// listen current user channel loop
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
