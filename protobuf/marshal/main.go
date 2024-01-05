package main

import (
	"github.com/golang/protobuf/proto" // 序列化用
	"test/protobuf/pb" // 导入生成的pb包
	"fmt"
	"time"
	"google.golang.org/protobuf/types/known/anypb"
    "google.golang.org/protobuf/types/known/wrapperspb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
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

	// proto中的google.protobuf.Any类型
	strWrapper := wrapperspb.String("good student")
    special, _ := anypb.New(strWrapper)
	
	// proto中的google.protobuf.Timestamp类型
	now := time.Now()
	last_updated := timestamppb.New(now)

	raw_data := &pb.Student{
		Id: 1,
		Name: "Andy",
		Gender: pb.Gender_GENDER_MALE,
		Scores: scores,
		PassStatus: false,
		Homeworks: homeworks,
		Address: address,
		Remark: &pb.Student_DefaultRemark{
			DefaultRemark: pb.DefaultRemark_REMARK_GOOD,
		},
		Special: special,
		UsualPerformances: usual_performances,
		LastUpdated: last_updated,
	}
	serialized_data, err := proto.Marshal(raw_data)
	if err != nil {
		fmt.Println("marshaling error: ", err)
	}

	decoded_data := &pb.Student{}
	err = proto.Unmarshal(serialized_data, decoded_data)
	if err != nil {
		fmt.Println("unmarshaling error: ", err)
	}

	// Now compare the variable contain the same data.
	if raw_data.GetName() != decoded_data.GetName() {
		fmt.Printf("data mismatch %q != %q\n", raw_data.GetName(), decoded_data.GetName())
	} else {
		fmt.Printf("data match %q == %q\n", raw_data.GetName(), decoded_data.GetName())
		fmt.Println("-----------------------------")
		fmt.Println(decoded_data.GetId())
		fmt.Println(decoded_data.GetName())
		fmt.Println(decoded_data.GetGender())
		fmt.Println("-----------------------------")
		fmt.Println(decoded_data.GetScores()) // 假设枚举SUBJECT_MATH=0，打印：[score:90 subject:SUBJECT_ENGLISH score:65.5]，缺少一个subject。原因是protobufV3为了减少无效数据的传输会放弃传输零值数据，比如枚举的0，bool的false，字符串的""。这是设计意图。如果造成困扰可以改用protobufV2并使用required限定。要用v3就避免在枚举中有效值使用零值。
		fmt.Println(decoded_data.GetScores()[0].Subject) // 打印：SUBJECT_MATH。虽然没有经过编码，但是因为是默认值，所以仍然可以读取。
		fmt.Println("-----------------------------")
		fmt.Println(decoded_data.GetPassStatus()) // 打印：false
		fmt.Println(decoded_data.GetHomeworks())
		fmt.Println("-----------------------------")
		fmt.Println(decoded_data.GetAddress()) // 打印：province:PROVINCE_BEIJING city:CITY_CAPITAL。address没有被编码，因为是""，和上面subject同样的逻辑。
		fmt.Println(decoded_data.GetAddress().Address) // 打印：空字符串。因为是默认值，也是可以读取的。
		fmt.Println("-----------------------------")
		fmt.Println(decoded_data.GetRemark()) // 打印：&{REMARK_GOOD}
		fmt.Println(decoded_data.GetDefaultRemark())  // 打印：REMARK_GOOD。pb中，为了实现oneof功能，实现了一个新接口isStudent_Remark，并为*Student实现了GetDefaultRemark()和GetOtherRemark()方法，可以将GetRemark()得到的isStudent_Remark类型数据转成指定的两种类型之一的数据。
		fmt.Println("-----------------------------")
		fmt.Println(decoded_data.GetSpecial()) // 打印：[type.googleapis.com/google.protobuf.StringValue]:{value:"good student"}
		fmt.Println(string(decoded_data.GetSpecial().Value)) // 不能用此方式获取数据，因为得到的是字节数组，包含了多种数据在内，不止是value
		m := wrapperspb.String("") // 编码和解码应该用同样的方式和包，由https://pkg.go.dev/google.golang.org/protobuf/types/known/wrapperspb#String提供
		_ = decoded_data.GetSpecial().UnmarshalTo(m)
		fmt.Printf("%+v\n",string(m.GetValue()))
		fmt.Println("-----------------------------")
		fmt.Println(decoded_data.GetUsualPerformances()) // 打印：map[ATTENDANCE:85 MIDTERM:70]。key为字符串类型，protobuf不允许map中放enum
		fmt.Println(decoded_data.GetUsualPerformances()["ATTENDANCE"])
		fmt.Println("-----------------------------")
		fmt.Println(decoded_data.GetLastUpdated()) // 打印：seconds:1704263229 nanos:164360459
		fmt.Println(decoded_data.GetLastUpdated().GetSeconds())
		fmt.Println(decoded_data.GetLastUpdated().GetNanos())
		fmt.Println(decoded_data.GetLastUpdated().AsTime()) // 将timestamppb格式的时间戳转成time.Time，但是缺少时区
		var cstZone = time.FixedZone("CST", 8*3600) // 东八区
		fmt.Println(decoded_data.GetLastUpdated().AsTime().In(cstZone).Format("2006-01-02 15:04:05"))
		fmt.Println(time.Now()) //

		// 其他类型的编解码和大量例子可参考：https://github.com/golang/protobuf/tree/master/ptypes
	}
}