package main

import (
	"test/protobuf/pb" // 导入生成的pb包
	"fmt"
	"net/rpc"
	"net/http"
	"errors"
	"github.com/golang/protobuf/proto" // 序列化用
	"time"
	"google.golang.org/protobuf/types/known/anypb"
    "google.golang.org/protobuf/types/known/wrapperspb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"encoding/gob"
)
// net/rpc包提供了rpc协议的实现。但是也有5个固定的编码要求：

// 作用：声明一个只在rpcserver和rpcclient使用的类型
// 要求1：类型是导出的
type StuService struct {}

// 作用：在类型上提供可供rpcclient访问的方法，访问时提供的参数为`"Type.FuncName"`
// 要求2：方法是导出的
// 要求3：方法的参数(argType T1, replyType *T2)，均为导出/内置类型
// 要求4：方法的第二个参数是指针类型
// 要求5：方法的返回值是error
func (s *StuService) GetStudent(req *pb.GetStudentRequest, res *pb.GetStudentResponse) error {
	fmt.Println("Got Rpc Request:", req)
	if req.Id == 1 {
		var err error
		res.Student, err = GetRawData()
		return err
	}
	return errors.New("Student not found")
}

func main() {
	rpc.RegisterName("StuService", new(StuService)) // 注册一个本地类型的指针到rpc服务列表中，并赋予别名
	rpc.HandleHTTP() // net/rpc协议是借助http实现的，所以需要启动http server服务
	if err := http.ListenAndServe(":1234", nil); err != nil {
		fmt.Println("Error serving: ", err)
	}
}

func GetRawData() (*pb.Student, error){
	scores := make([]*pb.Score, 0, 2) // 使用pb生成代码中的类型
	scores = append(scores, &pb.Score{
		Subject: pb.Subject_SUBJECT_MATH, // 注意pb会对枚举改名为Type_VARNAME。虽然会加前缀，但由于proto文件要求枚举变量名是全局唯一的VARNAME，所以VARNAME还是要带着`SUBJECT`
		Score: 90, // 因为Go的Struct导出变量需要大写，所以会将成员名大写
	})
	scores = append(scores, &pb.Score{
		Subject: pb.Subject_SUBJECT_ENGLISH,
		Score: 55.5,
	})

	homeworks := make([]*pb.Homework, 0, 1)
	homeworks = append(homeworks, &pb.Homework{
		Id: 1,
		Url: "http://www.baidu.com/homework",
	})

	address := &pb.Address{
		Address: "",
		Province: pb.Province_PROVINCE_BEIJING,
		City: pb.City_CITY_BEIJING,
	}

	usual_performances := make(map[string]float64, 2)
	// 读取枚举的自定义选项custom_name的值：
	// 枚举PerformanceType有自定义选项custom_name，声明在extend google.protobuf.EnumValueOptions中
	// 在pb.go中，定义了一个数组var file_main_proto_enumTypes = make([]protoimpl.EnumInfo, 4)，以及方法PerformanceType.Descriptor() -> protoreflect.EnumDescriptor
	// protoreflect是一个Go实现的动态操作protobuf消息的接口。它提供描述符descriptor描述了proto源文件中定义的类型的结构、检查和操作消息内容的接口
	// GetExtension用于从proto.MessageV1或proto.MessageV2提取指定的扩展项的值，返回interface{}
	t := pb.PerformanceType_PERFORMANCE_TYPE_ATTENDANCE // 
	tt, err := proto.GetExtension(proto.MessageV1(t.Descriptor().Values().ByNumber(t.Number()).Options()), pb.E_CustomName)
	usual_performances[*(tt.(*string))] = 85 // interface{}的底层类型是字符串指针
	t = pb.PerformanceType_PERFORMANCE_TYPE_MIDTERM
	tt, err = proto.GetExtension(proto.MessageV1(t.Descriptor().Values().ByNumber(t.Number()).Options()), pb.E_CustomName)
	usual_performances[*(tt.(*string))] = 70
	if err != nil {
		return nil, err
	}

	// proto中的google.protobuf.Any类型
	strWrapper := wrapperspb.String("good student")
    special, _ := anypb.New(strWrapper)
	
	// proto中的google.protobuf.Timestamp类型
	now := time.Now()
	last_updated := timestamppb.New(now)

	gob.Register(&pb.Student_DefaultRemark{})
	raw_data := &pb.Student{
		Id: 1,
		Name: "Andy",
		Gender: pb.Gender_GENDER_MALE,
		Scores: scores,
		PassStatus: false,
		Homeworks: homeworks,
		Address: address,
		// 提示：rpc: gob error encoding body: gob: type not registered for interface: pb.Student_DefaultRemark
		// 不传的项目将会被忽略
		// Remark: &pb.Student_DefaultRemark{
		// 	DefaultRemark: pb.DefaultRemark_REMARK_GOOD,
		// },
		Special: special,
		UsualPerformances: usual_performances,
		LastUpdated: last_updated,
	}
	return raw_data, nil
}