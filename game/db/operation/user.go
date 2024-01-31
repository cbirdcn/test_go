package operation

import (
	"context"
	"math/rand"
	"test/game/db/model"
	"test/game/pb"

	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// 提供bson结构
	// 管理连接等
	// 各种选项参数
	// "go.mongodb.org/mongo-driver/bson/primitive" // bson转换主键
)

const DB_GLOBAL = "global"
const COLL_USER = "user"

var err error

// 请求方法前需要先切换db和集合
func UseCollection(client **mongo.Client, coll **mongo.Collection) {
	*coll = (*client).Database(DB_GLOBAL).Collection(COLL_USER)
}

func AddUser(ctx context.Context, coll *mongo.Collection, req *pb.AddUserRequest) (uint64, error) {
	// TODO: transaction
	// TODO: 判断参数合法
	user := model.User{}
	err = GetUserByName(ctx, coll, req.Name, &user)
	if err == nil || err == mongo.ErrNoDocuments {
		var data = model.User{ // 传入变量地址
			Uid:  makeNewUid(),
			Name: req.Name,
		}
		_, err = coll.InsertOne(ctx, &data)
		if err == nil {
			return data.Uid, err
		} else {
			glog.Errorf("Error during InsertOne: %v", err)
		}
	} else {
		glog.Errorf("Error during GetUserByName: %v", err)
	}
	return 0, err
}

func GetUserByName(ctx context.Context, coll *mongo.Collection, name string, user *model.User) error {
	filter := bson.D{{"name", name}}
	return coll.FindOne(ctx, filter).Decode(user)
}

// TODO:临时
func makeNewUid() uint64 {
	return uint64(rand.Intn(100000000))
}

// start = time.Now().UnixMicro()
// insertRes, err := AddOneData(ctx, coll, &CollUserAccount{ // 传入变量地址
// 	Uid: uid,
// 	Name: name,
// })
// fmt.Println("spend micro", time.Now().UnixMicro() - start)
// if err != nil {
// 	fmt.Println(err)
// }
// // 返回类型为*mongo.InsertOneResult，如果要返回主键就是insertRes.InsertedID类型是interface{}。注意不是InsertedId
// // 如果要把"_id"用到mongo中，应该用bson.D{{"_id",res.InsertedID}}包裹起来，或转成类型primitive.ObjectID
// fmt.Println(insertRes)
// fmt.Println(insertRes.InsertedID)
// // fmt.Println(insertRes.InsertedID.(primitive.ObjectID).Hex()) // 不同版本的mongo-driver处理主键ID的方式不同，过去将bsonID转成字符串：interface{}->底层类型bsonID->十六进制字符串
// fmt.Println(insertRes.InsertedID.(string))

// func AddOneData(ctx context.Context, coll *mongo.Collection, documents interface{}) (*mongo.InsertOneResult, error){ // 注意structData传入时必须用指针，返回值是*mongo.InsertOneResult，如果需要res.InsertedID就是interface{}类型
// 	return coll.InsertOne(ctx, documents)
// }
