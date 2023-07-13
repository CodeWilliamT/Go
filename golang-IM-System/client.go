package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	//connect to server with tcp proxy
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error", err)
		return nil
	}
	client.conn = conn
	return client
}

func (client *Client) menu() bool {

	modes := map[int]string{
		1: "Public-Chat",
		2: "Private-Chat",
		3: "Rename",
		0: "Exit",
	}

	var flag int
	for i := 0; i < 4; i++ {
		fmt.Println(i, "."+modes[i])
	}
	fmt.Scanln(&flag)
	_, ok := modes[flag]
	if ok {
		client.flag = flag
		return true
	} else {
		fmt.Println("Please input valid value")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}
		switch client.flag {
		case 1:
			fmt.Println("Select Public-Chat")
			break
		case 2:
			fmt.Println("Select Private-Chat")
			break
		case 3:
			fmt.Println("Select Rename")
			break
		case 0:
			break
		}
	}
}

var serverIp string
var serverPort int

var modes map[int]string

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "To set server IP\n(default: 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "To set server port\n(default: 8888)")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>Failed to connect server...")
		return
	}
	fmt.Println(">>>>>Succeed to connect server...")
	client.Run()
}
