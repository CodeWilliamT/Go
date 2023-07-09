package main

import (
	"fmt"
	"net"
)

type Server struct{
	Ip string
	Port int
}


func NewServer(ip string,port int) *Server{
	server := &Server{
		Ip: ip,
		Port: port,
	}
	return server
}

func (this *Server) Handler(conn net.Conn){
	fmt.Println("Established connection")
}

func (this *Server) Start(){

	listener, err := net.Listen("tcp",fmt.Sprintf("%s:%d",this.Ip,this.Port))
	if err!= nil {
		fmt.Println("net.Listen err:",err)
		return
	}

	defer listener.Close()

	for{
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listenr accept err:",err)
			continue
		}

		//handler conn
		go this.Handler(conn)
	}
}