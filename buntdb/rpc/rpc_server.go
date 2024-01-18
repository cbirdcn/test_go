package main

import (
	"log"
	"fmt"
	"time"
	"github.com/tidwall/buntdb"
	"test/buntdb/rpc/pb" // 导入生成的pb包
	"net/rpc"
	"net/http"
	"errors"
)

var db *buntdb.DB
var err error
var b_config buntdb.Config

func init() {
	// 作为REQ的rpc server，初始化db客户端操作db
	db, err = buntdb.Open("data.db")
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close() // 只有在遇到中断信号时才能close

	if err := db.ReadConfig(&b_config); err != nil{
		log.Fatal(err)
	}
}

type DBService struct {}

func (s *DBService) Set(req *pb.SetRequest, res *pb.SetResponse) error {
	fmt.Println("Got Rpc Request:", req)
	if req.Op == "set" {
		t := time.Now().UnixMilli()
		err = db.Update(func(tx *buntdb.Tx) error {
			_, _, err := tx.Set(req.Key, req.Value, nil)
			return err
		})
		fmt.Println("db op cost ms:", time.Now().UnixMilli()-t) // db操作时间不稳定，100并发时耗时在0-17ms之间
		// 判断逻辑，如果发生错误要回滚。异步操作怎么通知请求方呢？加一个操作响应队列，通过msg_id回写状态、原因、返回值三项数据，新队列的消费者负责保持tcp连接并响应原始的客户端请求。
		res.MsgId = req.MsgId
		res.UserId = req.UserId
		res.Status = uint64(200)
		res.Reason = "ok"
		res.Data = "{}"
		return nil
	}

	return errors.New("Student not found")
}

func main() {
	// initConsumer()
	fmt.Println("监听rpc请求中...")
	rpc.RegisterName("DBService", new(DBService)) // 注册一个本地类型的指针到rpc服务列表中，并赋予别名
	rpc.HandleHTTP() // net/rpc协议是借助http实现的，所以需要启动http server服务
	if err := http.ListenAndServe(":1234", nil); err != nil {
		fmt.Println("Error serving: ", err)
	}
}