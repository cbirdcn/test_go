package main

import (
	"context"
	"encoding/json"
	"fmt"
	"errors"
	"time"
	"go.mongodb.org/mongo-driver/bson" // 提供bson结构
	"go.mongodb.org/mongo-driver/mongo" // 管理连接等
	"go.mongodb.org/mongo-driver/mongo/options" // 各种选项参数
	// "go.mongodb.org/mongo-driver/bson/primitive" // bson转换主键
)

// const MONGODB_URI = "mongodb://root:123456@host.docker.internal:27011/?maxPoolSize=10&minPoolSize=2&maxConnecting=2&w=mojority" // 单机
const MONGODB_URI = "mongodb://root:123456@172.22.0.3:27011,172.22.0.4:27012,172.22.0.5:27013/?directConnection=true" // 内外集群
// const MONGODB_URI = "mongodb://root:123456@host.docker.internal:27011,host.docker.internal:27012,host.docker.internal:27013/?directConnection=true" // 外网集群
const MONGODB_URI_WITHOUT_CREDENTIAL = "mongodb://host.docker.internal:27011/?directConnection=true"
const MONGODB_USERNAME = "root"
const MONGODB_PASSWORD = "123456"
const DB_GLOBAL = "global"
const COLL_USER_ACCOUNT = "user_account"

var (
	client *mongo.Client
	db *mongo.Database
	coll *mongo.Collection
	err error
)

func Connect(ctx context.Context, opts *options.ClientOptions) (*mongo.Client, error) {
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGODB_URI)) // 将账号密码作为连接字符串的方式
	client, err = mongo.Connect(ctx, opts) // 将账号密码作为options的方式。client是全局变量
	if err != nil {
		panic(err)
	}
	// 关闭连接。应用中可以不关闭。
	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// 注意：Connect()没有实际连接到服务器，只是在本地创建一个连接对象，需要请求（比如ping）后才确认连接是否联通。
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	// 到这里不会出现err==nil的情况，因为前面panic了，但是一般还是要返回
	return client, err
}

// 获取客户端连接，没有则初始化，有则直接使用
// GetMongoClient和Connect部分，及全局变量client可以独立成一个util.go包，导入进项目
func GetMongoClient(ctx context.Context, opts *options.ClientOptions) (*mongo.Client, error) {
	if client == nil {
		return Connect(ctx, opts)
	}
	return client, nil
}

// 类型定义（表声明）的好处是，可以选择性读写特定的列，并且避免编写列名和语句等防止出错
// 当然也可以嵌套结构体
// 在生产时，也应该独立到一个单独的model.go包中
type CollUserAccount struct {
	Uid string `bson:"_id", json:"uid"`	// 注意：bson规定了，写入到mongo对应的字段名。如果明确了`_id`对应的字段就不会再自动生成`_id`
	Name string `bson:"name", json:"name"`
}

// 为filter设定一个全新的struct，不能用CollStruct，因为没有给值的字段会被当做零值进行filter
// 但是如果很多条件怎么办？所以从条件查询的角度来说，bson.M更适合作为filter条件，只不过每次要注意一下拼写
type CollUserAccountFilterByUid struct {
	Uid string `bson:"uid"`
}

// 添加一行struct数据
func AddOneData(ctx context.Context, coll *mongo.Collection, documents interface{}) (*mongo.InsertOneResult, error){ // 注意structData传入时必须用指针，返回值是*mongo.InsertOneResult，如果需要res.InsertedID就是interface{}类型
	return coll.InsertOne(ctx, documents)
}

// 添加一行struct数据
func AddManyData(ctx context.Context, coll *mongo.Collection, documents []interface{}) (*mongo.InsertManyResult, error){ // 可以拿到所有插入的主键：res.InsertedIDs
	return coll.InsertMany(ctx, documents)
}

func UpsertOneData(ctx context.Context, coll *mongo.Collection, filter bson.D, update bson.D, opts *options.UpdateOptions) (*mongo.UpdateResult, error){
	return coll.UpdateOne(ctx, filter, update, opts) // 注意opts是指针
}

// 注意如果查询不到，就无法给指定变量赋值为解码的数据，并且返回的error就是mongo.ErrNoDocuments
// 也能添加options，比如skip，limit等，但是暂时不加，因为涉及迭代，单独处理
func FindOneDataToBson(ctx context.Context, coll *mongo.Collection, filter bson.D, result *bson.M) (error){
	return coll.FindOne(ctx, filter).Decode(result)
}

func FindOneDataToStructByBsonFilter(ctx context.Context, coll *mongo.Collection, filter bson.D, result *CollUserAccount) (error){ // FindOne()返回*mongo.SingleResult，再Decode()返回error
	return coll.FindOne(ctx, filter).Decode(result)
}

func FindOneDataToStructByStructFilter(ctx context.Context, coll *mongo.Collection, filter *CollUserAccountFilterByUid, result *CollUserAccount) (error){ // FindOne()返回*mongo.SingleResult，再Decode()返回error
	return coll.FindOne(ctx, filter).Decode(result)
}

func DeleteManyData(ctx context.Context, coll *mongo.Collection, filter bson.M) (*mongo.DeleteResult, error){
	return coll.DeleteMany(ctx, filter)
}

// 获取结果的一页数据
func FindManyDataToStructByBsonFilter(ctx context.Context, coll *mongo.Collection, filter bson.D, opts ...*options.FindOptions) (*mongo.Cursor, error){ // 注意可变参数的类型写法
	return coll.Find(ctx, filter, opts...)
}

// 迭代Find()所有结果到一个结果集中，需要提供分页迭代参数
/*
举例：
skipPage := 0 // 跳过了多少页，跳过的记录数就是skipPage*limit
total := 4 // 初始总数量，如果没给就用默认初始值，可以定为100
totalInc := 4 // 如果查完total还没查完，就让total=total+totalInc。可以理解为总量翻倍扩容，则直接让total作为inc。
limit := 2 // 每次查limit条记录，如果有一次没查到记录就停止循环
*/
const DEFAULT_ITERATOR_INITIAL_TOTAL_COUNT = 100
const DEFAULT_ITERATOR_INITIAL_TOTAL_MAX_COUNT = 100000
func IteratorFindManyDataToStructByBsonFilter(ctx context.Context, coll *mongo.Collection, filter bson.D, opts *options.FindOptions, startSkipPage, initialTotalCount, totalInc, limit int64)  ([]CollUserAccount, error){
	// 
	findManyRes := make([]CollUserAccount, 0, initialTotalCount)
	if initialTotalCount < 0 {
		return findManyRes, errors.New("iterator find initialTotalCount must >= 0")
	}
	if initialTotalCount >= DEFAULT_ITERATOR_INITIAL_TOTAL_MAX_COUNT {
		return findManyRes, errors.New("iterator find initialTotalCount over max")
	}
	if limit <= 0 {
		return findManyRes, errors.New("iterator find limit must > 0")
	}
	if initialTotalCount == 0 {
		initialTotalCount = DEFAULT_ITERATOR_INITIAL_TOTAL_COUNT
	}
	if totalInc <= 0 {
		totalInc = initialTotalCount
	}
	total := initialTotalCount
	skipPage := startSkipPage

	var cursor *mongo.Cursor
	defer func(){
		if err = cursor.Close(ctx); err != nil {
			fmt.Println(err)
		}
	}()
	
	for {
		findManyResTemp := make([]CollUserAccount, 0, limit)
		opts.SetSkip(int64(skipPage * limit))
		opts.SetLimit(int64(limit))
		cursor, err = FindManyDataToStructByBsonFilter(ctx, coll, filter, opts)
		if err != nil {
			return findManyRes, err
		}
		if err = cursor.All(ctx, &findManyResTemp); err != nil { // 即使读取的数量少于limit就结束了，也不会报错，当读取结果为空时表示查询结束
			return findManyRes, err
		}
		if len(findManyResTemp) > 0 {
			if skipPage * limit + limit >= total {
				total = total + totalInc
			}
			skipPage++
			findManyRes = append(findManyRes, findManyResTemp...)
		} else {
			break
		}
	}
	return findManyRes, err
}


func main() {

	start := time.Now().UnixMicro()
	// 创建客户端连接。
	ctx := context.TODO()
	credential := options.Credential{
		Username: MONGODB_USERNAME,
		Password: MONGODB_PASSWORD,
	}
	opts := options.Client().ApplyURI(MONGODB_URI_WITHOUT_CREDENTIAL).SetAuth(credential)
	client, err = GetMongoClient(ctx, opts) // 已声明全局变量

	// 指定db和集合
	coll = client.Database(DB_GLOBAL).Collection(COLL_USER_ACCOUNT)

	fmt.Println("spend micro", time.Now().UnixMicro() - start)

	// 数据
	uid := "13732795"
	name := "Bob"

	// 删除所有数据
	fmt.Println("----------------")
	fmt.Println("DeleteManyData")
	start = time.Now().UnixMicro()
	deleteRes, err := DeleteManyData(ctx, coll, bson.M{})
	fmt.Println("spend micro", time.Now().UnixMicro() - start)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(deleteRes.DeletedCount)


	// 直接插入一条数据，如果数据出现唯一约束错误，会返回error
	fmt.Println("----------------")
	fmt.Println("AddOneData")
	start = time.Now().UnixMicro()
	insertRes, err := AddOneData(ctx, coll, &CollUserAccount{ // 传入变量地址
		Uid: uid,
		Name: name,
	})
	fmt.Println("spend micro", time.Now().UnixMicro() - start)
	if err != nil {
		fmt.Println(err)
	}
	// 返回类型为*mongo.InsertOneResult，如果要返回主键就是insertRes.InsertedID类型是interface{}。注意不是InsertedId
	// 如果要把"_id"用到mongo中，应该用bson.D{{"_id",res.InsertedID}}包裹起来，或转成类型primitive.ObjectID
	fmt.Println(insertRes)
	fmt.Println(insertRes.InsertedID)
	// fmt.Println(insertRes.InsertedID.(primitive.ObjectID).Hex()) // 不同版本的mongo-driver处理主键ID的方式不同，过去将bsonID转成字符串：interface{}->底层类型bsonID->十六进制字符串
	fmt.Println(insertRes.InsertedID.(string))


	// 直接插入多条数据，如果数据出现唯一约束错误，会返回error
	fmt.Println("----------------")
	fmt.Println("AddManyData")
	start = time.Now().UnixMicro()
	insertManyRes, err := AddManyData(ctx, coll, []interface{}{
		&CollUserAccount{ // 传入变量地址
			Uid: uid + "a",
			Name: name,
		},
		&CollUserAccount{
			Uid: uid + "b",
			Name: name,
		},
	})
	fmt.Println("spend micro", time.Now().UnixMicro() - start)
	if err != nil {
		fmt.Println(err)
	}
	// 返回类型为*mongo.InsertOneResult，如果要返回主键就是insertRes.InsertedID类型是interface{}。注意不是InsertedId
	// 如果要把"_id"用到mongo中，应该用bson.D{{"_id",res.InsertedID}}包裹起来，或转成类型primitive.ObjectID
	fmt.Println(insertManyRes)
	fmt.Println(insertManyRes.InsertedIDs)
	for _, v := range insertManyRes.InsertedIDs {
		fmt.Println(v)
	}


	// 不存在则插入
	fmt.Println("----------------")
	fmt.Println("UpsertOneData")
	start = time.Now().UnixMicro()
	// filter := bson.D{{"_id", uid}} // 注意bson.D{}是固定格式，里面才是表达式，所以不能少写最外层的{}
	filter := bson.D{{"uid", uid}} // 注意：不提供"_id"，也没指定哪个字段是主键时，mongo创建记录时将自动生成一个"_id"字段
	update := bson.D{{"$setOnInsert", bson.D{{"uid", uid}, {"name", name}}}}
	update_opts := options.Update().SetUpsert(true) // 当match失败时，如果upsert=true，将创建新数据。注意：返回*options.UpdateOptions
	update_record, err := UpsertOneData(ctx, coll, filter, update, update_opts) // 注意：UpdateOne返回的类型和InsertOne返回的类型不同，所以返回值不能赋值给同一个变量
	// update_record, err := coll.UpdateOne(context.TODO(), filter, update, update_opts)
	fmt.Println("spend micro", time.Now().UnixMicro() - start)
	if err != nil {
		fmt.Println(err)
	} else {
		// 如果filter已匹配，则match=1,upsert=0表示匹配到所以不插入不更新。
		// 如果filter未匹配，则match=0,upsert=1表示未匹配到所以插入。
		fmt.Println(update_record.MatchedCount) // 如果发生err，将无法取得count，因为res将是值为<nil>的interface{}。报错：type interface{} has no field or method MatchedCount
		fmt.Println(update_record.UpsertedCount) // upserted行数
	}


	fmt.Println("----------------")
	fmt.Println("FindOneDataToBson")
	var result bson.M
	start = time.Now().UnixMicro()
	err = FindOneDataToBson(ctx, coll, bson.D{{"uid", uid}}, &result) // 传入地址
	// err = coll.FindOne(context.TODO(), bson.D{{"uid", uid}}).Decode(&result)
	fmt.Println("spend micro", time.Now().UnixMicro() - start)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the uid %s\n", uid)
		} else {
			fmt.Println(err)
		}
	}
	fmt.Println(result) // 如果返回err出错，result就不会被赋值


	fmt.Println("----------------")
	fmt.Println("FindOneDataToStructByBsonFilter")
	var resStruct1 = CollUserAccount{}
	start = time.Now().UnixMicro()
	err = FindOneDataToStructByBsonFilter(ctx, coll, bson.D{{"uid", uid}}, &resStruct1) // 传入地址。用bson.D作为filter的方式
	// err = coll.FindOne(context.TODO(), bson.D{{"uid", uid}}).Decode(&resStruct1)
	fmt.Println("spend micro", time.Now().UnixMicro() - start)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the uid %s\n", uid)
		} else {
			fmt.Println(err)
		}
	}
	fmt.Println(resStruct1)


	fmt.Println("----------------")
	fmt.Println("FindOneDataToStructByStructFilter")
	start = time.Now().UnixMicro()
	var resStruct2 = CollUserAccount{}
	// 用struct查询时，不能直接在Struct中添加filter作为条件。因为Struct中没添加的字段也会作为零值参与filter过程的，会出现错误。所以需要单独为filter设定一个新的struct才行
	err = FindOneDataToStructByStructFilter(ctx, coll, &CollUserAccountFilterByUid{Uid: uid}, &resStruct2) // 传入地址。用struct作为filter的方式
	// err = coll.FindOne(context.TODO(), bson.D{{"uid", uid}}).Decode(&resStruct2)
	fmt.Println("spend micro", time.Now().UnixMicro() - start)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the uid %s\n", uid)
		} else {
			fmt.Println(err)
		}
	}
	fmt.Println(resStruct2)


	// 查询多个文档，并迭代输出
	fmt.Println("----------------")
	fmt.Println("FindManyDataToStructByBsonFilter")
	// 注意skip表示跳过的数量，limit表示此次查询的总数。所以skip0limit1就是查1条。但是skip跳过时还是要一条条地去跳过，skip数量太大效率不好
	// 官方未实现深度分页功能，如果要用skip，需要自己在迭代过程中调整skip和limit值
	/*
	逻辑：skip是跳过的记录数=pageSize * page，limit是要读的数量
	//第一页
	db.collection.find().skip(0).limit(2) // 跳过0条记录，读2条
	//第二页
	db.collection.find().skip(2).limit(2) // 跳过2条记录，读2条
	//第三页
	db.collection.find().skip(4).limit(2) // 跳过4条记录，读2条
	*/
	iteratorFilter := bson.D{{"name", name}}
	startSkipPage := 0
	initialTotalCount := 4
	totalInc := 4
	limit := 2
	findManyRes := make([]CollUserAccount, 0, initialTotalCount)
	findOpts := options.Find() // 技巧：不需要传递多个*mongo.FindOptions，可以定义一个options.Find()返回*mongo.FindOptions结构体，后面的所有SetSort、SetSkip等都只是对这个FindOptions的修改，所以传入一个FindOptions就够了
	findOpts.SetSort(bson.D{{"_id", 1}}) // 不需要SetSkip和SetLimit，内部会自动拼接
	start = time.Now().UnixMicro()
	findManyRes, err = IteratorFindManyDataToStructByBsonFilter(ctx, coll, iteratorFilter, findOpts, int64(startSkipPage), int64(initialTotalCount), int64(totalInc), int64(limit))
	fmt.Println("spend micro", time.Now().UnixMicro() - start)
	if err != nil {
		fmt.Println(err)
	}
	for _, findLine := range findManyRes {
		fmt.Println(findLine)
	}

	
	fmt.Println("----------------")
	fmt.Println("json")
	start = time.Now().UnixMicro()
	jsonData, err := json.MarshalIndent(result, "", "    ")
	fmt.Println("spend micro", time.Now().UnixMicro() - start)
	fmt.Printf("%s\n", jsonData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", jsonData)

  
}


/*
第一遍运行：
spend micro 31633
----------------
DeleteManyData
spend micro 2705
0
----------------
AddOneData
spend micro 5843
&{13732795}
13732795
13732795
----------------
AddManyData
spend micro 6044
&{[13732795a 13732795b]}
[13732795a 13732795b]
13732795a
13732795b
----------------
UpsertOneData
spend micro 11571
0
1
----------------
FindOneDataToBson
spend micro 4099
map[_id:ObjectID("659fac479211060923596836") name:Bob uid:13732795]
----------------
FindOneDataToStructByBsonFilter
spend micro 3848
{659fac479211060923596836 Bob}
----------------
FindOneDataToStructByStructFilter
spend micro 2213
{659fac479211060923596836 Bob}
----------------
FindManyDataToStructByBsonFilter
spend micro 5099
{13732795 Bob}
{13732795a Bob}
{13732795b Bob}
{659fac479211060923596836 Bob}
----------------
json
spend micro 67
{
    "_id": "659fac479211060923596836",
    "name": "Bob",
    "uid": "13732795"
}
{
    "_id": "659fac479211060923596836",
    "name": "Bob",
    "uid": "13732795"
}




第二遍运行：
spend micro 36338
----------------
DeleteManyData
spend micro 4118
4
----------------
AddOneData
spend micro 5002
&{13732795}
13732795
13732795
----------------
AddManyData
spend micro 5249
&{[13732795a 13732795b]}
[13732795a 13732795b]
13732795a
13732795b
----------------
UpsertOneData
spend micro 4057
0
1
----------------
FindOneDataToBson
spend micro 3336
map[_id:ObjectID("659fac719211060923596a4e") name:Bob uid:13732795]
----------------
FindOneDataToStructByBsonFilter
spend micro 2381
{659fac719211060923596a4e Bob}
----------------
FindOneDataToStructByStructFilter
spend micro 2913
{659fac719211060923596a4e Bob}
----------------
FindManyDataToStructByBsonFilter
spend micro 9384
{13732795 Bob}
{13732795a Bob}
{13732795b Bob}
{659fac719211060923596a4e Bob}
----------------
json
spend micro 138
{
    "_id": "659fac719211060923596a4e",
    "name": "Bob",
    "uid": "13732795"
}
{
    "_id": "659fac719211060923596a4e",
    "name": "Bob",
    "uid": "13732795"
}
*/


// 官方mongo：https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo
// 官方options：https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo/options
// github：https://github.com/ypankaj007/golang-mongodb-restful-starter-kit
// 知乎：https://zhuanlan.zhihu.com/p/144308830

/*
时间计算：
在3C3G的三个容器中，部署了三副本分片集群，与Go容器在同一网络组。
连接耗时36ms，每次读写单行数据耗时1-2ms。
这要是一次性更新1000行数据，可能耗时1s！
*/

/*
特点：https://blog.51cto.com/u_12592884/2771762
MongoDB与MySQL对比优点是，基本读写性能翻倍。保证千万数据表可用，即使无索引的范围扫描，也能有500+QPS，mysql则会宕机。
也就这些优点，稳定性不如MySQL，可能出bug。driver也不好用。
*/