package main

import (
	"flag"
	"os"
	"os/signal"
	"test/game/pb"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var linkServerAddr = "ws://localhost:8888"
var done chan interface{}
var interrupt chan os.Signal
var err error

func receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		glog.Infof("Receive at Micro: %v", time.Now().UnixMicro())
		if err != nil {
			glog.Errorf("Error in Receive: %v", err)
			return
		}
		// 编码后的[]byte，无法直接查看
		glog.Infof("Received Msg: %v", msg)

		// 解码
		decoded_data := &pb.AddUserResponse{}
		err = proto.Unmarshal(msg, decoded_data)
		if err != nil {
			glog.Errorf("Unmarshaling error: %v", err)
		}
		glog.Infof("Decoded Code: %v", decoded_data.GetCode())
		glog.Infof("Decoded Uid: %v", decoded_data.GetUid())
	}
}

// 注意：wsClient和wsServer作为连接和消息中转服务，传输的数据是
func main() {
	// 如果不提供参数log_dir，将默认生成到/tmp。注意重名文件能生成以不同进程ID为文件名的日志，但因为软链接只有一个`文件名.日志级别`，所以只能过滤出其中一个进程的日志数据。
	// alsologtostderr参数默认为false。
	flag.Parse()
	defer glog.Flush()

	done = make(chan interface{})    // Channel to indicate that the receiverHandler is done
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	socketUrl := linkServerAddr + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		glog.Fatalf("Error connecting to Websocket Server: %v", err)
	}
	defer conn.Close()
	go receiveHandler(conn)

	// Our main loop for the client
	// We send our relevant packets here
	for {
		select {
		case <-time.After(time.Duration(1) * time.Millisecond * 1000):
			// Send an echo packet every second
			// msg := "Hello"
			var req = pb.AddUserRequest{
				Name: "Andy",
			}
			// glog.Infof("Send Msg Struct: %#v", &req) // 注意struct中包含pb添加的state等属性，比较长，不如String()方便
			serialized_data, err := proto.Marshal(&req)
			if err != nil {
				glog.Errorf("Marshaling error: %v", err)
			}
			msg := string(serialized_data)

			glog.Infof("Send Msg String: %s", msg)
			glog.Infof("Send at Micro: %v", time.Now().UnixMicro())
			err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				glog.Errorf("Error during writing to websocket: %v", err)
				return
			}

		case <-interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			glog.Info("Received SIGINT interrupt signal. Closing all pending connections")

			// Close our websocket connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				glog.Errorf("Error during closing websocket: %v", err)
				return
			}

			select {
			case <-done:
				// 没有给channel写入数据时无法读取数据。但是当channel关闭时，就能从channel读取数据了，是零值。
				glog.Info("Receiver Channel Closed! Exiting....")
			case <-time.After(time.Duration(1) * time.Second):
				glog.Info("Timeout in closing receiving channel. Exiting....")
			}
			return // 如果for外还有要处理的，用break
		}
	}
}
