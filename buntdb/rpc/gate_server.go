package main

import (
	"test/buntdb/rpc/pb" // 导入生成的pb包
	"fmt"
	"net/rpc"
	"time"
	"net/http"
	"math/rand"
	"github.com/gin-gonic/gin"
)

var client *rpc.Client
var user_id int64 = 987654321
var server_id int64 = 1000

func init() {
	client, _ = rpc.DialHTTP("tcp", "localhost:1234") // 通过rpc包的连接到HTTP服务器，返回一个支持rpc调用的client变量。一次连接，重复使用。
}

func main() {
	fmt.Println("监听http请求中...")
	r := gin.Default()
	r.GET("/id/:id", func(c *gin.Context) {
		id := c.Param("id")
		msg_id := GetMsgId()
		res := CallRpc(msg_id, user_id, id)
		c.String(http.StatusOK, res.Data) // 或者作为json返回res
		// TODO: 需要对阻塞计时，超时做出特定的http响应
	})

	r.Run(":8888")
}

func GetMsgId() (int64){
	// 消息的唯一性由两部分决定：Id = 玩家角色id（不定长）；MsgId（64位） = 符号位0（第1位） - 毫秒时间（41位仿雪花）- 服务器id（12位最大4095） - 随机数(10位最大1023)
	// Id是唯一的，这里只讨论MsgId
	now := time.Now().UnixMilli() // 微秒：1705499944917313，毫秒：1705500048195（二进制41位），一年大约31104000秒，当前秒时间戳1705500964（转成二进制占31位）
	var rand_int int64 = int64(rand.Intn(1024)) // 不包含1024
	return now << 22 + server_id << 10 + rand_int
}

func CallRpc(msg_id int64, user_id int64, id string) (res pb.SetResponse){
	start := time.Now().UnixMilli()
	req := pb.SetRequest{
		MsgId: uint64(msg_id),
		UserId: uint64(user_id),
		Op: "set",
		Key: id,
		Value: id,
	}
	if err := client.Call("DBService.Set", &req, &res); err != nil { // client提交三个参数："Type.FuncName"、req值或指针、res指针，到服务器，响应将被写入到res中，并返回error
		fmt.Println("Failed to call DBService.Set.", err)
	} else {
		fmt.Println("rpc spend ms: ", time.Now().UnixMilli() - start)
		fmt.Println("Got Rpc Response:", res)
		// fmt.Println("Got Rpc Response:", res.Data)
		// fmt.Println("Got Rpc Response not exist item:", res.Emp) // 访问未传递的item为nil而不是报错
	}
	return
}

/*

容器内业务测试

go run rpc_server.go

go run gate_server.go

curl 127.0.0.1:8888/id/1

rpc_server打印：
监听rpc请求中...
Got Rpc Request: msg_id:7153655651805143597 user_id:987654321 op:"set" key:"1" value:"1"
db op cost ms: 1

gate_server打印：
监听http请求中...
rpc spend ms:  6
Got Rpc Response: {{{} [] [] <nil>} 0 [] 7153655651805143597 987654321 200 ok {}}
[GIN] 2024/01/18 - 15:53:35 | 200 |    8.163208ms |      172.22.0.1 | GET      "/id/1"

结论：
数据逻辑简单
需要额外编排proto文件，可以生成基于string等基础类型的结构
耦合性低，并发支持好，数据库层面的并发与服务器并发无关，需要数据库自己处理好并发访问问题
多次rpc请求之间是无序的，可以通过消息id+角色id辅助确定消息顺序。
性能好，耗时较少

*/

/*
容器内并发测试

100并发
Label	# Samples	Average	Min	Max	Std. Dev.	Error %	Throughput	Received KB/sec	Sent KB/sec	Avg. Bytes
HTTP Request	100	20	5	54	12.11	0.00%	98.81423	11.39	11.58	118
TOTAL	100	20	5	54	12.11	0.00%	98.81423	11.39	11.58	118

500并发
Label	# Samples	Average	Min	Max	Std. Dev.	Error %	Throughput	Received KB/sec	Sent KB/sec	Avg. Bytes
HTTP Request	500	1107	17	1752	415.51	0.00%	187.05574	21.56	21.92	118
TOTAL	500	1107	17	1752	415.51	0.00%	187.05574	21.56	21.92	118
*/

/*
2C8Gmacmini真机并发测试，jmeter局域网请求

100并发
Label	# Samples	Average	Min	Max	Std. Dev.	Error %	Throughput	Received KB/sec	Sent KB/sec	Avg. Bytes
HTTP Request	100	2	1	5	0.77	0.00%	100.40161	11.57	12.35	118
TOTAL	100	2	1	5	0.77	0.00%	100.40161	11.57	12.35	118

500并发
Label	# Samples	Average	Min	Max	Std. Dev.	Error %	Throughput	Received KB/sec	Sent KB/sec	Avg. Bytes
HTTP Request	500	3	2	10	1.06	0.00%	497.51244	57.33	61.22	118
TOTAL	500	3	2	10	1.06	0.00%	497.51244	57.33	61.22	118

5000并发（推荐），平均耗时12ms
Label	# Samples	Average	Min	Max	Std. Dev.	Error %	Throughput	Received KB/sec	Sent KB/sec	Avg. Bytes
HTTP Request	5000	12	2	126	17.13	0.00%	2653.92781	305.82	326.56	118
TOTAL	5000	12	2	126	17.13	0.00%	2653.92781	305.82	326.56	118

50000并发，注意错误率67.32%
Label	# Samples	Average	Min	Max	Std. Dev.	Error %	Throughput	Received KB/sec	Sent KB/sec	Avg. Bytes
HTTP Request	50000	426	1	1583	351.04	67.32%	1217.49294	1884.43	48.95	1584.9
TOTAL	50000	426	1	1583	351.04	67.32%	1217.49294	1884.43	48.95	1584.9
*/