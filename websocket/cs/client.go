// client.go
package main
 
import (
    "log"
    "os"
    "os/signal"
    "time"
 
    "github.com/gorilla/websocket"
)
 
var done chan interface{}
var interrupt chan os.Signal
 
func receiveHandler(connection *websocket.Conn) {
    defer close(done)
    for {
        _, msg, err := connection.ReadMessage()
        if err != nil {
            log.Println("Error in receive:", err)
            return
        }
        log.Printf("Received: %s\n", msg)
    }
}
 
func main() {
    done = make(chan interface{}) // Channel to indicate that the receiverHandler is done
    interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully
 
    signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT
 
    socketUrl := "ws://localhost:8080" + "/ws"
    conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
    if err != nil {
        log.Fatal("Error connecting to Websocket Server:", err)
    }
    defer conn.Close()
    go receiveHandler(conn)
 
    // Our main loop for the client
    // We send our relevant packets here
    for {
        select {
        case <-time.After(time.Duration(1) * time.Millisecond * 1000):
            // Send an echo packet every second
            err := conn.WriteMessage(websocket.TextMessage, []byte("Hello from GolangDocs!"))
            if err != nil {
                log.Println("Error during writing to websocket:", err)
                return
            }
 
        case <-interrupt:
            // We received a SIGINT (Ctrl + C). Terminate gracefully...
            log.Println("Received SIGINT interrupt signal. Closing all pending connections")
 
            // Close our websocket connection
            err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
            if err != nil {
                log.Println("Error during closing websocket:", err)
                return
            }
 
            select {
            case <-done:
                // 没有给channel写入数据时无法读取数据。但是当channel关闭时，就能从channel读取数据了，是零值。
                log.Println("Receiver Channel Closed! Exiting....")
            case <-time.After(time.Duration(1) * time.Second):
                log.Println("Timeout in closing receiving channel. Exiting....")
            }
            return // 如果for外还有要处理的，用break
        }
	}
}