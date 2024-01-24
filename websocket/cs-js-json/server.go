package main

import (
    "fmt"
    "net/http"
	"encoding/json"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	if err != nil {
		// panic(err)
		fmt.Println(err)
		return
	}
	// defer conn.Close()

	for {
		// Read message
		msgType, reqByte, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the message
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(reqByte))

		resString, err := makeResponse(string(reqByte), r)
		if err != nil {
			// TODO: 特定返回值
		}

		// Write message back
		if err = conn.WriteMessage(msgType, []byte(resString)); err != nil {
			return
		}
	}
}

type Req struct {
	R string `json:"r"` // route
	D string `json:"d"` // data:json
	U string `json:"u"` // user_id
	T string `json:"t"` // timestamp
}

type ReqParam struct {
	// TODO:对Req解码解构后的数据
}

func makeResponse(reqString string, r *http.Request) (resString string, err error){
	req := Req{}
	err = json.Unmarshal([]byte(reqString), &req)
	if err != nil {
		fmt.Println(err)
		return "", err // TODO:要给标准错误返回，比如解码失败
	}
	// TODO: req解构到reqParam，并用reqParam做后续处理
	switch req.R {
		case "/set": resString, err = handleSet(&req)
		default: resString, err = handleDefault(&req)
	}
	return
}

func handleSet(req *Req) (string, error){
	return "haha", nil
}

func handleDefault(req *Req) (string, error){
	return "", nil
}

func main() {
    http.HandleFunc("/ws", handleWebSocket)

    http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "websockets.html")
    })

    http.ListenAndServe(":8888", nil)
}

/*

go run server.go 

浏览器打开http://127.0.0.1:8888/index，确认Status: Connected
输入：
{"r":"/set", "d":"{}", "u":"1", "t":"1700000000"}
界面显示回声"haha"，查看网络耗时2-3ms
*/