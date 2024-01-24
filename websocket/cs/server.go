// server.go
package main
 
import (
    "log"
    "net/http"
	"fmt"
 
    "github.com/gorilla/websocket"
)
 
var upgrader = websocket.Upgrader{} // use default options
 
func wsHandler(w http.ResponseWriter, r *http.Request) {
    // Upgrade our raw HTTP connection to a websocket based one
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("Error during connection upgradation:", err)
        return
    }
    defer conn.Close()

    // The event loop
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error during message reading:", err)
            break
        }
        // 读写数据都是字节数组，要打印需要转型string
        fmt.Println(message)
        log.Printf("Received: %s", message)
        // 回声：服务器对客户端的每个消息都原样返回
        err = conn.WriteMessage(messageType, message)
        if err != nil {
            log.Println("Error during message writing:", err)
            break
        }
    }
}
 
func home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Index Page")
}
 
func main() {
    fmt.Println("等待连接中...")
    http.HandleFunc("/ws", wsHandler) // 注册一个处理器，管理连接请求
    http.HandleFunc("/", home) // 处理业务请求
    log.Fatal(http.ListenAndServe("localhost:8080", nil))
}