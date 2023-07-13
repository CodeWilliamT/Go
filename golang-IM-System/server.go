package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//Online user list
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//Common Message channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// listen server msg
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		this.mapLock.Lock()
		//sync msg with all users
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()

	}
}

// Notice user online
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ": " + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//test
	//fmt.Println("Established connection")
	user := NewUser(conn, this)
	user.Online()
	isLive := make(chan bool)

	//User Message Reader
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			//Trans user msg without "\n"
			msg := string(buf[:n-1])

			//Broadcast msg
			user.SendMessage(msg)
			isLive <- true
		}
	}()

	//moniter if offline
	for {
		select {
		case <-isLive:
		case <-time.After(time.Second * 60):
			user.GetMessage("You has been kicked")
			close(user.C)
			conn.Close()
			return
		}
	}
}

func (this *Server) Start() {

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	defer listener.Close()

	//start server message listener
	go this.ListenMessage()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listenr accept err:", err)
			continue
		}

		//handler conn
		go this.Handler(conn)
	}
}
