package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MONGODB_URI = "mongodb://admin:admin@172.21.0.2:27017/?maxPoolSize=10&minPoolSize=2&maxConnecting=2&w=mojority"

func main() {
	fmt.Println(time.Now().UnixNano() / 1e6)
	uri := MONGODB_URI
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("global").Collection("user_account")
	uid := "13732795"

	//var result bson.M
	type UserAccount struct {
		ID	primitive.ObjectID	`bson:"_id"`
		Uid	string	`bson:"uid, omitempty"`
		IsTourist	uint8	`bson:"is_tourist"`
		CreatedAt	primitive.DateTime	`bson:"created_at,omitempty"`
	}
	var result UserAccount
	fmt.Println(time.Now().UnixNano() / 1e6)
	err = coll.FindOne(context.TODO(), bson.D{{"uid", uid}}).Decode(&result)
	fmt.Println(time.Now().UnixNano() / 1e6)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the uid %s\n", uid)
		return
	}
	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	fmt.Println(time.Now().UnixNano() / 1e6)
	fmt.Printf("%s\n", jsonData)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
	fmt.Println(time.Now().UnixNano() / 1e6)
}
