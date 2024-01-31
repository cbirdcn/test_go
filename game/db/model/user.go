package model

type User struct {
	// Id   string `bson:"_id", json:"id"`
	Uid  uint64 `bson:"uid", json:"uid"` // 注意：bson规定了，写入到mongo对应的字段名。如果明确了`_id`对应的字段就不会再自动生成`_id`
	Name string `bson:"name", json:"name"`
}
