package dbServer

import (
	"flag"
	"net/http"
	"net/rpc"
	"test/game/db/operation"
	"test/game/pb" // 导入生成的pb包

	"context" // 提供bson结构

	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/mongo"         // 管理连接等
	"go.mongodb.org/mongo-driver/mongo/options" // 各种选项参数
	// "go.mongodb.org/mongo-driver/bson/primitive" // bson转换主键
	// "github.com/golang/protobuf/proto" // 序列化用
	// "google.golang.org/protobuf/types/known/anypb"
	// "google.golang.org/protobuf/types/known/timestamppb"
	// "google.golang.org/protobuf/types/known/wrapperspb"
)

// const MONGODB_URI_REPL = "mongodb://root:123456@172.23.0.2:27017,172.23.0.4:27017,172.23.0.3:27017/?directConnection=true" // 内外集群
// const MONGODB_URI = "mongodb://root:123456@host.docker.internal:27011,host.docker.internal:27012,host.docker.internal:27013/?directConnection=true" // 外网集群
// const MONGODB_URI_WITHOUT_CREDENTIAL = "mongodb://host.docker.internal:27011/?directConnection=true"
const MONGODB_URI_SINGLE = "mongodb://172.23.0.2:27017/?directConnection=true"
const MONGODB_USERNAME = "root"
const MONGODB_PASSWORD = "123456"

var logicServerAddr = ":3333"

var (
	client *mongo.Client
	coll   *mongo.Collection
	err    error
	ctx    context.Context
)

func init() {
	ctx = context.TODO()
	credential := options.Credential{
		Username: MONGODB_USERNAME,
		Password: MONGODB_PASSWORD,
	}
	opts := options.Client().ApplyURI(MONGODB_URI_SINGLE).SetAuth(credential)
	client, err = GetMongoClient(ctx, opts) // 已声明全局变量
}

func GetMongoClient(ctx context.Context, opts *options.ClientOptions) (*mongo.Client, error) {
	if client == nil {
		return Connect(ctx, opts)
	}
	return client, nil
}

func Connect(ctx context.Context, opts *options.ClientOptions) (*mongo.Client, error) {
	glog.Info("connect mongodb...")
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGODB_URI)) // 将账号密码作为连接字符串的方式
	client, err = mongo.Connect(ctx, opts) // 将账号密码作为options的方式。client是全局变量
	if err != nil {
		glog.Errorf("Error during connect to mongodb: %v", err)
	}

	// 注意：Connect()没有实际连接到服务器，只是在本地创建一个连接对象，需要请求（比如ping）后才确认连接是否联通。
	err = client.Ping(ctx, nil)
	if err != nil {
		glog.Errorf("Error during ping mongodb: %v", err)
	}

	glog.Info("connected to mongodb...")
	// 到这里不会出现err==nil的情况，因为前面panic了，但是一般还是要返回
	return client, err
}

// net/rpc包提供了rpc协议的实现。但是也有5个固定的编码要求：

// 作用：声明一个只在rpcserver和rpcclient使用的类型
// 要求1：类型是导出的
type DbService struct{}

// 作用：在类型上提供可供rpcclient访问的方法，访问时提供的参数为`"Type.FuncName"`
// 要求2：方法是导出的
// 要求3：方法的参数(argType T1, replyType *T2)，均为导出/内置类型
// 要求4：方法的第二个参数是指针类型
// 要求5：方法的返回值是error
func (s *DbService) AddUser(req *pb.AddUserRequest, res *pb.AddUserResponse) (err error) {
	// TODO: 设计模式：在基类的init()中打印请求
	// glog.Infof("Req: %#v", req) // 注意struct中包含pb添加的state等属性，比较长，不如String()方便
	glog.Infof("Req String: %s", req.String())
	// TODO: 设计模式：提供repo
	operation.UseCollection(&client, &coll)
	uid, err := operation.AddUser(ctx, coll, req)
	if err != nil {
		glog.Errorf("Error: %v", err)
		res.Code = pb.StatusCode_STATUS_CODE_ADD_USER_FAIL
	} else {
		res.Code = pb.StatusCode_STATUS_CODE_OK
	}
	res.Uid = uid
	// TODO: 设计模式：在基类的最后打印响应
	glog.Infof("Res: %v", res.String())
	return
}

func Main() {
	flag.Parse()
	defer glog.Flush()

	glog.Info("db service start...")
	rpc.RegisterName("DbService", new(DbService)) // 注册一个本地类型的指针到rpc服务列表中，并赋予别名
	rpc.HandleHTTP()                              // net/rpc协议是借助http实现的，所以需要启动http server服务
	if err := http.ListenAndServe(logicServerAddr, nil); err != nil {
		glog.Errorf("Error serving: %v", err)
	}
}
