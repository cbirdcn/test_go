package logicServer

import (
	"flag"
	"net/http"
	"net/rpc"
	"test/game/pb" // 导入生成的pb包

	"github.com/golang/glog"
	"google.golang.org/protobuf/proto"
)

var logicServerAddr = ":1234"
var dbServerAddr = "localhost:3333"
var rpcClient *rpc.Client
var err error

func connectDbService() {
	// 通过rpc包的连接到HTTP服务器，返回一个支持rpc调用的client
	rpcClient, err = rpc.DialHTTP("tcp", dbServerAddr)
	if err != nil {
		glog.Fatalf("Error during connect db rpc:", err)
	}
}

// net/rpc包提供了rpc协议的实现。但是也有5个固定的编码要求：

// 作用：声明一个只在rpcserver和rpcclient使用的类型
// 要求1：类型是导出的
type LogicService struct{}

// 作用：在类型上提供可供rpcclient访问的方法，访问时提供的参数为`"Type.FuncName"`
// 要求2：方法是导出的
// 要求3：方法的参数(argType T1, replyType *T2)，均为导出/内置类型
// 要求4：方法的第二个参数是指针类型
// 要求5：方法的返回值是error
func (s *LogicService) RunDeliver(req *string, res *string) error {
	// 注意：
	glog.Infof("Req: %s", *req)
	// *res = "hello:" + *req
	// *res = "callback"
	// 从其他服务获取数据
	// 处理逻辑
	// 响应

	// 对请求消息解码（proto编码过的[]byte-->pb.AddUserRequest)
	decoded_data := &pb.AddUserRequest{}
	err = proto.Unmarshal([]byte(*req), decoded_data)
	if err != nil {
		glog.Errorf("Unmarshaling error: %v", err)
	}
	glog.Infof("Decoded Name: %v", decoded_data.GetName())

	// 实际上应该再加一个请求响应结构，这里暂时借用上一层的结构
	var dbReq pb.AddUserRequest
	var dbRes pb.AddUserResponse
	dbReq.Name = decoded_data.GetName()
	// client提交三个参数："Type.FuncName"、req值或指针、res指针，到rpc服务器，响应将被写入到res中，并返回error
	glog.Infof("DbReq Struct: %#v", &dbReq)
	glog.Infof("DbReq String: %s", dbReq.String())
	if err := rpcClient.Call("DbService.AddUser", &dbReq, &dbRes); err != nil {
		glog.Errorf("Error during call rpc: %v", err)
		return err
	} else {
		glog.Infof("DbRes: %v", dbRes.String())
		// 注意：正常情况下，从db获取数据后要经过logic处理后转成响应类型AddUserResponse，这里为了方便直接将db的返回也用AddUserResponse类型了，而省略了logic处理和转换类型的环节。
		// 另外，在这之前都可以用String()打印。最后经过proto.Marshal编码后就是二进制字节数组，无法打印查看。
		serialized_data, err := proto.Marshal(&dbRes)
		if err != nil {
			glog.Errorf("Marshaling error: %v", err)
		}
		*res = string(serialized_data)
		glog.Info("transfer binary")
		// fmt.Println("Send Rpc Response:", *res)
	}
	return nil
}

func main() {
	flag.Parse()
	defer glog.Flush()

	glog.Info("logic service start...")
	glog.Info("connect db service...")

	connectDbService()
	glog.Info("linked to db service...")

	rpc.RegisterName("LogicService", new(LogicService)) // 注册一个本地类型的指针到rpc服务列表中，并赋予别名
	rpc.HandleHTTP()                                    // net/rpc协议是借助http实现的，所以需要启动http server服务
	if err := http.ListenAndServe(logicServerAddr, nil); err != nil {
		glog.Errorf("Error serving:", err)
	}
}
