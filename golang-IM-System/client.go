package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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

// deal server response
func (client *Client) DealResponse() {
	//alway waiting for response
	io.Copy(os.Stdout, client.conn)
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

func (client *Client) ListUsers() {
	sendMsg := "-userlist\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err: ", err)
		return
	}
}

func (client *Client) PrivateChat() {
	remoteName := ""
	chatMsg := ""

	for remoteName != "exit" {
		client.ListUsers()
		fmt.Println(">>>>Please input target username,\"exit\" to exit")
		fmt.Scanln(&remoteName)
		for chatMsg != "exit" {
			fmt.Println(">>>>Please input your message, \"exit\" to exit")
			fmt.Scanln(&chatMsg)
			if len(chatMsg) != 0 {
				sendMsg := "-to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err: ", err)
					break
				}
			}
			chatMsg = ""
		}
	}
}

func (client *Client) PublicChat() {
	chatMsg := ""

	for chatMsg != "exit" {
		fmt.Println(">>>>Please input your message, \"exit\" to exit")
		fmt.Scanln(&chatMsg)
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err: ", err)
				break
			}
		}
		chatMsg = ""
	}

}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>Please enter User Name:")
	fmt.Scanln(&client.Name)

	sendMsg := "-rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}
		switch client.flag {
		case 1:
			fmt.Println("Select Public-Chat")
			client.PublicChat()
			break
		case 2:
			fmt.Println("Select Private-Chat")
			client.PrivateChat()
			break
		case 3:
			fmt.Println("Select Rename")
			client.UpdateName()
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

	//start a thread to deal response
	go client.DealResponse()

	fmt.Println(">>>>>Succeed to connect server...")
	client.Run()
}
