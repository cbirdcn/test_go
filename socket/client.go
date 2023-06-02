package main

import (
	"fmt"
	"net"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9999"
	SERVER_TYPE = "tcp"
)

func main() {
	// establish connection
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	defer connection.Close()
	if err != nil {
		panic(err)
	}

	// send some data
	_, err = connection.Write([]byte("Hello World."))
	if err != nil {
		fmt.Println("Error writing:", err.Error())
	}

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer) // 返回有效消息的字节长度
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Client received:", string(buffer[:mLen]))
}
