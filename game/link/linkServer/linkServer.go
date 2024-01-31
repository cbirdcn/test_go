package linkServer

import (
	"flag"
	"net/http"
	"net/rpc"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

var linkServerAddr = ":8888"
var logicServerAddr = "localhost:1234"

var upgrader = websocket.Upgrader{} // use default options
var rpcClient *rpc.Client
var err error

func connectLogicService() {
	// 通过rpc包的连接到HTTP服务器，返回一个支持rpc调用的client
	rpcClient, err = rpc.DialHTTP("tcp", logicServerAddr)
	if err != nil {
		glog.Fatalf("Error during connect logic rpc: %v", err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Errorf("Error during connection upgradation: %v", err)
		return
	}
	defer conn.Close()

	// simple event loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			glog.Errorf("Error during connection upgradation: %v", err)
			// 待定：从连接读取失败，将退出循环，中断连接
			break
		}

		// 将字节数组转成string后，原样转发给logic_rpc
		var req = string(message)
		var res string
		// client提交三个参数："Type.FuncName"、req值或指针、res指针，到rpc服务器，响应将被写入到res中，并返回error
		if err := rpcClient.Call("LogicService.RunDeliver", &req, &res); err != nil {
			glog.Errorf("Error during call rpc: %v", err)
		} else {
			glog.Infof("Req: %s", req)
			// 从logic_rpc拿到响应后原样返回给wsClient。注意，logic必须具备构造wsResponse的能力。
			err = conn.WriteMessage(messageType, []byte(res))
			if err != nil {
				glog.Errorf("Error during message writing: %v", err)
				break
			}
			glog.Info("Res: binary")
		}

	}
}

func main() {
	flag.Parse()
	defer glog.Flush()

	glog.Info("link service start...")
	glog.Info("connect logic service...")

	connectLogicService()
	glog.Info("linked to logic service...")
	http.HandleFunc("/ws", wsHandler) // 注册一个处理器，管理连接请求

	glog.Fatalf("serving error: %v", http.ListenAndServe(linkServerAddr, nil))
	glog.Info("link service close...")
}
