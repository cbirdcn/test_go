package main

import (
	"fmt"
	"net"
	"os"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9999"
	SERVER_TYPE = "tcp"
)

func main() {
	fmt.Println("Server Running...")
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST + ":" + SERVER_PORT) // 监听本机端口
	if err != nil {
		panic(err.Error())
	}
	defer server.Close() // 延迟关闭服务器监听

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting for client...")

	// 循环处理请求
	for {
		// waiting for connection 等待客户端连接
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Client connected")
		go processClient(connection) // 开启协程，处理信息交互
	}

}

func processClient(connection net.Conn) {
	defer connection.Close()
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer) //Read reads data from the connection. And then save to buffer. Return the length of message byte.
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Server received:", string(buffer[:mLen])) // 从buffer读取本次传输的消息
	_, err = connection.Write([]byte("Thanks, Server got your msg:" + string(buffer[:mLen]))) // 向连接写入msg（发还给客户端）
}