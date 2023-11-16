package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// MONGODB_URI must have a / before the query ?
const MONGODB_URI = "mongodb://admin:admin@172.21.0.3:27017/?maxPoolSize=10&minPoolSize=5&maxConnection=5&w=majority"

// CRUD、Upsert、Replace
// 索引
// 多字段
// 聚合 $replaceAll
// struct与bson
// No Document
// 插入和查询，在字符串和数值
// bulkWrite与writeModel
// distinct
// count
// mapReduce
// remove
// transaction

type Users struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Age       int32              `bson:"age" json:"age"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
}

func main() {
	// connect
	uri := MONGODB_URI
	mongoCtx := context.Background()
	client, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// defer: disconnect connection
	defer func() {
		if err = client.Disconnect(mongoCtx); err != nil {
			panic(err)
		}
	}()


	// ping
	fmt.Println("ping...")
	err = client.Ping(mongoCtx, readpref.Primary())
	if err != nil {
		fmt.Println("ping err")
		panic(err)
	}
	fmt.Println("ping ok")


	// select db and collection
	coll := client.Database("testing").Collection("users")


	// deleteOne by ObjectID
	fmt.Println()
	fmt.Println("deleteOne by ObjectID...")
	deleteOneByIDID, err := primitive.ObjectIDFromHex("5eb3d668b31de5d588f42a7a")
	if err != nil {
		fmt.Println("deleteOne by ObjectID: generate ObjectID err")
		fmt.Println(err)
	}
	deleteOneByIDFilter := bson.M{
		"_id": deleteOneByIDID,
	}
	deleteOneByIDRes, err := coll.DeleteOne(mongoCtx, deleteOneByIDFilter)
	if err != nil {
		fmt.Println("deleteOne by ObjectID err")
		fmt.Println(err)
	}
	fmt.Println(deleteOneByIDRes.DeletedCount)


	// deleteOne by filter
	fmt.Println()
	fmt.Println("deleteOne by filter...")
	deleteOneByFilterFilter := bson.M{
		"name": "Andrew",
	}
	deleteOneByFilterRes, err := coll.DeleteOne(mongoCtx, deleteOneByFilterFilter)
	if err != nil {
		fmt.Println("deleteOne by filter err")
		fmt.Println(err)
	}
	fmt.Println(deleteOneByFilterRes.DeletedCount)


	// deleteMany
	fmt.Println()
	fmt.Println("deleteMany...")
	deleteManyFilter := bson.D{}
	deleteManyRes, err := coll.DeleteMany(mongoCtx, deleteManyFilter)
	if err != nil {
		fmt.Println("deleteMany err")
		fmt.Println(err)
	}
	fmt.Println(deleteManyRes.DeletedCount)


	// drop index
	fmt.Println()
	fmt.Println("drop collection...")
	// 删除集合。Delete只是删除数据，仍然保留集合结构和索引等。
	err = coll.Drop(mongoCtx)
	if err != nil {
		fmt.Println("drop collection err")
		fmt.Println(err)
	} else {
		fmt.Println("drop collection ok")
	}


	// create index
	fmt.Println()
	fmt.Println("createMany index...")
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{"name", 1},
				{"age", 1},
				{"created_at", 1},
			},
			Options: options.Index().SetName("idx_name_age_create"),
		},{
			Keys: bson.D{
				{"name", 1},
			},
			Options: options.Index().
				SetName("idx_name").
				SetUnique(true).
				SetCollation(&options.Collation{
					Locale:          "en_US",
					CaseLevel:       false,
					CaseFirst:       "",
					Strength:        0,
					NumericOrdering: false,
					Alternate:       "",
					MaxVariable:     "",
					Normalization:   false,
					Backwards:       false,
				}),
		},{
			Keys: bson.D{
				{"age", 1},
			},
			Options: options.Index().SetName("idx_age").SetMin(0),
		},{
			Keys: bson.D{
				{"created_at", 1},
			},
			Options: options.Index().SetName("idx_create"),
		},
	}
	createIndexOpts := options.CreateIndexes().SetMaxTime(2 * time.Second)
	indexNames, err := coll.Indexes().CreateMany(mongoCtx, indexModels, createIndexOpts)
	if err != nil {
		fmt.Println("createMany index err")
		fmt.Println(err)
	}
	fmt.Println(indexNames)
	// 可以使用db.users.find({"name": "Andrew"}).explain()查看索引使用情况


	// insertOne
	fmt.Println()
	fmt.Println("insertOne...")
	insertOneData := bson.D{
		{"name", "Andrew"},
		{"age", 10},
		{"created_at", time.Now()},
	}
	insertOneRes, err := coll.InsertOne(mongoCtx, insertOneData)
	if err != nil {
		fmt.Println("InsertOne err")
		fmt.Println(err)
	}
	fmt.Println(insertOneRes.InsertedID)


	// insertMany
	fmt.Println()
	fmt.Println("insertMany...")
	insertManyData := []interface{}{
		bson.D{
			{"name", "Bob"},
			{"age", 20},
			{"created_at", time.Now()},
		},
		bson.D{
			{"name", "Cindy"},
			{"age", 30},
			{"created_at", time.Now()},
		},
	}
	insertManyRes, err := coll.InsertMany(mongoCtx, insertManyData)
	if err != nil {
		fmt.Println("InsertMany err")
		fmt.Println(err)
	}
	fmt.Println(insertManyRes.InsertedIDs)


	// findOne by bson.M
	fmt.Println()
	fmt.Println("findOne by bson.M...")
	findOneByBsonMFilter := bson.D{
		{"name", "Andrew"},
	}
	var findOneByBsonMRes bson.M
	err = coll.FindOne(mongoCtx, findOneByBsonMFilter).Decode(&findOneByBsonMRes)
	if err != nil {
		fmt.Println("findOne by bson.M decode err")
		fmt.Println(err)
	}
	findOneByBsonMResJson, err := json.Marshal(findOneByBsonMRes)
	if err != nil {
		fmt.Println("findOne by bson.M json marshal err")
		fmt.Println(err)
	}
	fmt.Println(findOneByBsonMRes)
	fmt.Println(findOneByBsonMRes["_id"])
	fmt.Println(findOneByBsonMRes["created_at"])	// millisecond, eg: 1688223289265，即2023年7月1日22:54:29
	fmt.Println(string(findOneByBsonMResJson))	// CreatedAt为带timezone的日期时间 eg: "2023-07-01T14:54:49.265Z"，解析时根据不同时区得到不同时间


	// findOne by struct
	fmt.Println()
	fmt.Println("findOne by struct...")
	var findOneByStruct Users
	findOneByStructFilter := bson.D{
		{"name", "Andrew"},
	}
	err = coll.FindOne(mongoCtx, findOneByStructFilter).Decode(&findOneByStruct)
	if err != nil {
		fmt.Println("findOne by struct decode err")
		fmt.Println(err)
	}
	findOneByStructJson, err := json.Marshal(findOneByStruct)
	if err != nil {
		fmt.Println("findOne by struct json marshal err")
		fmt.Println(err)
	}
	fmt.Println(findOneByStruct)
	fmt.Println(findOneByBsonMRes["_id"])
	fmt.Println(findOneByBsonMRes["created_at"])
	fmt.Println(string(findOneByStructJson))


	// find by bson.M
	fmt.Println()
	fmt.Println("find by bson.M...")
	// $and/$or must be an array
	// where (age >= 10 and age <= 20) or (name = "Cindy")
	findByBsonMFilter := bson.D{
		{"$or", []bson.M{
			{"age": bson.M{
				"$gte": 10,
				"$lte": 20,
			}},
			{"name": bson.M{
				"$eq": "Cindy",
			}},
		}},
	}
	findByBsonMCursor, err := coll.Find(mongoCtx, findByBsonMFilter)
	if err != nil {
		fmt.Println("find by BsonM find err")
		fmt.Println(err)
	}
	var findByBsonMRes []bson.M		// results must be a pointer to a slice
	if err = findByBsonMCursor.All(mongoCtx, &findByBsonMRes); err != nil {
		fmt.Println("find by BsonM decode err")
		fmt.Println(err)
	}
	fmt.Println(findByBsonMRes)
	for _, v := range findByBsonMRes {
		fmt.Println(v["_id"])
	}


	// updateOne by ID
	fmt.Println()
	fmt.Println("updateOne by ID...")
	updateOneByID := findOneByBsonMRes["_id"]
	//updateOneByID, err := primitive.ObjectIDFromHex("5eb3d668b31de5d588f42a7a")
	//if err != nil {
	//	fmt.Println("updateOne by ID make ObjectID err")
	//	fmt.Println(err)
	//}
	updateOneByIDUpdater := bson.D{
		{"$set", bson.M{
			"age": 11,
		}},
	}
	updateByIDRes, err := coll.UpdateByID(mongoCtx, updateOneByID, updateOneByIDUpdater)
	if err != nil {
		fmt.Println("updateOne by ID update err")
		fmt.Println(err)
	}
	fmt.Println(updateByIDRes.MatchedCount)	// 匹配
	fmt.Println(updateByIDRes.ModifiedCount)	// 修改


	// updateMany by filter
	fmt.Println()
	fmt.Println("updateMany by filter...")
	updateManyByFilter := bson.M{
		"age": bson.M{
			"$gte": 10,
			"$lte": 20,
		},
	}
	updateManyByFilterUpdater := bson.D{
		{"$mul", bson.M{
			"age": 1.2,
		}},
	}
	updateManyByFilterRes, err := coll.UpdateMany(mongoCtx, updateManyByFilter, updateManyByFilterUpdater)
	if err != nil {
		fmt.Println("updateMany by filter update err")
		fmt.Println(err)
	}
	fmt.Println(updateManyByFilterRes.MatchedCount)
	fmt.Println(updateManyByFilterRes.ModifiedCount)


	// upsertOne by filter：FindOneAndUpdate
	// FindOneAndUpdate相比UpdateOne、UpdateMany而言，可以将修改后的记录通过Decode解码。
	fmt.Println()
	fmt.Println("upsertOne by filter: FindOneAndUpdate...")
	upsertOneByFilterFilter := bson.M{
		"age": bson.M{
			"$gte": 30,
		},
	}
	// update document must contain key beginning with '$'
	upsertOneByFilterUpdater := bson.M{
		"$set": bson.M{
			"age": 40,
		},
	}
	// SetOptions:
	// Upsert: true/false.
	// ReturnDocument: Optional. When true, returns the updated document instead of the original document.Defaults to false.
	findAndUpsertOneByFilterOptions := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var findAndUpsertOneByFilterRes bson.M
	err = coll.FindOneAndUpdate(mongoCtx, upsertOneByFilterFilter, upsertOneByFilterUpdater, findAndUpsertOneByFilterOptions).Decode(&findAndUpsertOneByFilterRes)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("findAndUpsertOne by filter FindOne err")
			fmt.Println(err)
		} else {
			fmt.Println("findAndUpsertOne by filter Update err")
			fmt.Println(err)
		}
	}
	fmt.Println(findAndUpsertOneByFilterRes["age"])	// 40


	// upsertOne by filter：UpdateOne
	// UpdateOne、UpdateMany只获取操作的记录行数、ID等信息，不解码记录。
	fmt.Println()
	fmt.Println("upsertOne by filter: UpdateOne...")
	upsertOneByFilterUpdater2 := bson.M{
		"$set": bson.M{
			"age": 41,
		},
	}
	upsertOneByFilterOptions := options.Update().SetUpsert(true)
	// 二选一，区分Update还是Insert：
	// 1.filter匹配时，执行UpdateOne过程
	upsertOneByFilterRes, err := coll.UpdateOne(mongoCtx, upsertOneByFilterFilter, upsertOneByFilterUpdater2, upsertOneByFilterOptions)
	// 2.filter不匹配时，Upsert就是指Insert，res的Upsert结果也就是Insert结果
	//upsertOneByFilterFilter2 := bson.M{
	//	"age": bson.M{
	//		"$gte": 50,
	//	},
	//}
	//upsertOneByFilterRes, err := coll.UpdateOne(mongoCtx, upsertOneByFilterFilter2, upsertOneByFilterUpdater2, upsertOneByFilterOptions)
	if err != nil {
		fmt.Println("upsertOne by filter err")
		fmt.Println(err)
	}
	fmt.Println(upsertOneByFilterRes)
	fmt.Println(upsertOneByFilterRes.MatchedCount)	// 匹配，区分插入还是更新：Insert返回0, Update返回count
	fmt.Println(upsertOneByFilterRes.ModifiedCount) // 更新行数，区分插入还是更新：Insert返回0，Update返回count
	fmt.Println(upsertOneByFilterRes.UpsertedCount)	// 不存在才插入的行数，区分插入还是更新：Insert返回count, Update返回0
	fmt.Println(upsertOneByFilterRes.UpsertedID) // 不存在才插入的ID，区分插入还是更新：Insert返回新的ObjectID，Insert返回nil
	var findAfterUpsertRes bson.M
	_ = coll.FindOne(mongoCtx, upsertOneByFilterFilter).Decode(&findAfterUpsertRes)
	fmt.Println(findAfterUpsertRes)


	// replaceOne by filter
	// replace和update的区别：
	// 1. 使用 replaceOne() 只能替换整个文档，而 updateOne() 允许更新字段。
	// 由于 replaceOne() 替换了整个文档，因此旧文档中没有包含在新文档中的字段将会丢失。
	// 使用 updateOne() 可以在不丢失旧文档中的字段的情况下添加新字段。
	// 2. update有updateMany()。replace不提供replaceMany，只能在aggregation聚合中，用$replaceAll
	// { $replaceAll: { input: <expression>, find: <expression>, replacement: <expression> } }
	fmt.Println()
	fmt.Println("replaceOne by filter...")
	replaceOneByFilterFilter := bson.M{
		"name": "Cindy",
	}
	// replacement document cannot contain keys beginning with '$'
	// 1.当只填写部分field时，原field将会被删除，只留下新field
	replaceOneByFilterReplacement := bson.M{
		"age": 30,
	}
	replaceOneByFilterRes, err := coll.ReplaceOne(mongoCtx, replaceOneByFilterFilter, replaceOneByFilterReplacement)
	if err != nil {
		fmt.Println("replaceOne by filter err")
		fmt.Println(err)
	}
	fmt.Println(replaceOneByFilterRes.MatchedCount)
	fmt.Println(replaceOneByFilterRes.ModifiedCount)
	var findAfterReplaceOneRes bson.M
	_ = coll.FindOne(mongoCtx, replaceOneByFilterFilter).Decode(&findAfterReplaceOneRes)
	// 将会输出map[]。因为replacement中只有age，replace时原来的name值会被删掉，但是新的field只有age。所以通过name将无法再查出数据
	fmt.Println(findAfterReplaceOneRes)

	// 2.只有填写全部field时，才能完全覆盖原数据
	replaceOneByFilterFilter2 := bson.M{
		"age": 30,
	}
	replaceOneByFilterReplacement2 := bson.M{
		"name": "Cindy",
		"age": 30,
		"created_at": time.Now(),
	}
	_, _ = coll.ReplaceOne(mongoCtx, replaceOneByFilterFilter2, replaceOneByFilterReplacement2)
	var findAfterReplaceOneRes2 bson.M
	_ = coll.FindOne(mongoCtx, replaceOneByFilterFilter).Decode(&findAfterReplaceOneRes2)
	// 将会输出map[_id:ObjectID("64a11972b3f166c5077121c1") age:30 name:Cindy]。
	fmt.Println(findAfterReplaceOneRes2)


	// aggregate
	fmt.Println()
	fmt.Println("aggregate...")
	stage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"name", 1},
			{"age", "$age"},
			{"create", "$created_at"},
		}},
	}
	opts := options.Aggregate().SetMaxTime(2 * time.Second)
	cursor, err := coll.Aggregate(mongoCtx, mongo.Pipeline{stage}, opts)
	if err != nil {
		fmt.Println("aggregate err")
		fmt.Println(err)
	}
	var results []bson.M
	if err = cursor.All(mongoCtx, &results); err != nil {
		fmt.Println("aggregate decode err")
		fmt.Println(err)
	}
	for _,v := range results {
		fmt.Println(v)
	}


	//{"totalAge", bson.D{
	//	{"$sum", "$age"},
	//}},

}