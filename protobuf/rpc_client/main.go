package main

import (
	"test/protobuf/pb" // 导入生成的pb包
	"fmt"
	"net/rpc"
	"time"
)

func main() {
	start_dial := time.Now().UnixMicro()
	client, _ := rpc.DialHTTP("tcp", "localhost:1234") // 通过rpc包的连接到HTTP服务器，返回一个支持rpc调用的client变量
	start_req := time.Now().UnixMicro()
	var req pb.GetStudentRequest // 声明传入的请求和要接收响应的变量
	var res pb.GetStudentResponse
	req.Id = 1
	if err := client.Call("StuService.GetStudent", &req, &res); err != nil { // client提交三个参数："Type.FuncName"、req值或指针、res指针，到服务器，响应将被写入到res中，并返回error
		fmt.Println("Failed to call StuService.GetStudent.", err)
	} else {
		end_res := time.Now().UnixMicro()
		fmt.Println("rpc spend microseconds from dial: ", end_res - start_dial) // 2392
		fmt.Println("rpc spend microseconds from req: ", end_res - start_req) // 1266
		fmt.Println("Got Rpc Response:", res.Student)
		fmt.Println("Got Rpc Response not exist item:", res.Student.Remark) // 访问未传递的item为nil而不是报错
		// 解码其他item的过程看marshal
	}
	// 输出：
	// Got Rpc Response: id:1  name:"Andy"  gender:GENDER_MALE  scores:{score:90}  scores:{subject:SUBJECT_ENGLISH  score:65.5}  pass_status:true  homeworks:{id:1  url:"http://www.baidu.com/homework"}  address:{province:PROVINCE_BEIJING  city:CITY_CAPITAL}  special:{type_url:"type.googleapis.com/google.protobuf.StringValue"  value:"\n\x0cgood student"}  usual_performances:{key:"ATTENDANCE"  value:85}  usual_performances:{key:"MIDTERM"  value:70}  last_updated:{seconds:1704262802  nanos:361134067}
	// 分析：1. scores中第一个元素没有subject
	// 2. special的值是{type_url:"type.googleapis.com/google.protobuf.StringValue"  value:"\n\x0cgood student"}
	// 3. last_updated的值是{seconds:1704262802  nanos:361134067}

	req.Id = 100
	if err := client.Call("StuService.GetStudent", &req, &res); err != nil {
		fmt.Println("Failed to call StuService.GetStudent.", err)
	} else {
		fmt.Println("Got Rpc Response:", res.Student)
	}

}