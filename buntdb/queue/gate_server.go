package main

import (
	"fmt"
	"time"
	// "strconv"
	"math/rand"
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/Shopify/sarama"
)


var k_config *sarama.Config
var k_client sarama.SyncProducer
var err error
var consumer sarama.Consumer
var partitionList []int32
var locks = make(map[UMsg](chan RES))

var user_id int64 = 987654321
var server_id int64 = 1000

const TOPIC_REQ string = "req"
const TOPIC_RES string = "res"

func init() {
	// kafka连接初始化
	k_config = sarama.NewConfig() // 默认配置
	k_config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	k_config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	k_config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	k_client, err = sarama.NewSyncProducer([]string{"host.docker.internal:9092"}, k_config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	// defer k_client.Close()  // TODO:只有在遇到中断信号时才能close
}

func main() {
	// RES队列消费者，异步阻塞
	// 注意：消费者必须在发送Req消息前就调用。也就是先启动消费，再启动生产者。否则消费者还没启动，消息读取模式是最新，可能会被并发的其他消费者消费掉刚生产的消息，消费者启动后才发现：failed to start consumer for partition
	fmt.Println("监听res队列中...")
	go BlockConsume()

	fmt.Println("监听http请求中...")
	r := gin.Default()
	r.GET("/id/:id", func(c *gin.Context) {
		id := c.Param("id")
		// 生成消息id
		msg_id := GetMsgId()
		// 发送 REQ 消息到队列
		Set(id, msg_id)

		// 发出REQ消息后陷入阻塞，等待RES消息的返回
		res := <-locks[UMsg{
			UserId: user_id,
			MsgId: msg_id,
		}]
		// TODO: 需要对阻塞计时，超时做出特定的http响应
		// http响应
		c.String(http.StatusOK, res.Data) // 或者将res作为json响应
	})

	r.Run(":8888")
}

type REQ struct {
	Op string `json:"op"`
	UserId int64 `json:"user_id"`
	MsgId int64 `json:"msg_id"`
	Key string `json:"key"`
	Value string `json:"value"`
}

type RES struct {
	UserId int64 `json:"user_id"`
	MsgId int64 `json:"msg_id"`
	Status int64 `json:"status"`
	Reason string `json:"reason"`
	Data string `json:"data"`
}

type UMsg struct {
	UserId int64
	MsgId int64
}

func GetMsgId() (int64){
	// 消息的唯一性由两部分决定：Id = 玩家角色id（不定长）；MsgId（64位） = 符号位0（第1位） - 毫秒时间（41位仿雪花）- 服务器id（12位最大4095） - 随机数(10位最大1023)
	// Id是唯一的，这里只讨论MsgId
	now := time.Now().UnixMilli() // 微秒：1705499944917313，毫秒：1705500048195（二进制41位），一年大约31104000秒，当前秒时间戳1705500964（转成二进制占31位）
	var rand_int int64 = int64(rand.Intn(1024)) // 不包含1024
	return now << 22 + server_id << 10 + rand_int
}

func Set(id string, msg_id int64) {
	// 如果用string
	// msg_id := strconv.FormatInt(msg_id_int, 10)
	req := REQ{
		Op: "set",
		UserId: user_id,
		MsgId: msg_id,
		Key: id,
		Value: id,
	}
	json, _ := json.Marshal(req)

	// 发送REQ前，初始化RES锁，否则发送后理解读取Channel会报错未初始化的空指针。
	// locks锁是map，key是唯一的角色+消息id，value是channel，RES队列的消息者会向channel写入RES的返回值
	locks[UMsg{
		UserId: user_id,
		MsgId: msg_id,
	}] = make(chan RES)

	fmt.Println("发送消息id到REQ队列", msg_id)
	SendMessage(TOPIC_REQ, string(json))
	return
}

func SendMessage(topic string, value string) {
	t := time.Now().UnixMilli()
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(value)

	// 发送消息
	pid, offset, err := k_client.SendMessage(msg)
	if err != nil {
		fmt.Println("send req msg failed, err:", err)
		return
	}
	fmt.Printf("req partition_id:%v offset:%v\n", pid, offset)
	fmt.Println("gate server send req cost ms:", time.Now().UnixMilli()-t) // 写操作队列时间不稳定，100并发耗时在4-57ms之间
}

// 消费res消息，select写法
func BlockConsume() {
	consumer, err = sarama.NewConsumer([]string{"host.docker.internal:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err = consumer.Partitions(TOPIC_RES) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}

	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(TOPIC_RES, int32(partition), sarama.OffsetNewest) // OffsetOldest表示取过去的消息，OffsetNewest取最新消息也就是还没产生的消息随产生随消费
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			go func() {
				for msg := range pc.Messages() {
					fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
					res := RES{}
					err := json.Unmarshal(msg.Value, &res)
					//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
					if err != nil {
						fmt.Println(err)
					}
					locks[UMsg{
						UserId: res.UserId,
						MsgId: res.MsgId,
					}]<-res

				}
			}()
		}(pc)
	}
	select{}
}

/*
容器内测试
go run db_server.go

go run gate_server.go

curl 127.0.0.1:8888/id/1

db_server打印：
监听req队列中...
Partition:0 Offset:829 Key:[] Value:{"op":"set","user_id":987654321,"msg_id":7153649552423100557,"key":"1","value":"1"}
db op cost ms: 1
res partition_id:0 offset:829
db server send res cost ms: 18

gate_server打印：
监听res队列中...
监听http请求中...
发送消息id到REQ队列 7153649552423100557
req partition_id:0 offset:829
gate server send req cost ms: 8
Partition:0 Offset:829 Key:[] Value:{"user_id":987654321,"msg_id":7153649552423100557,"status":200,"reason":"ok","data":"{}"}
[GIN] 2024/01/18 - 15:29:21 | 200 |   40.754243ms |      172.22.0.1 | GET      "/id/1"

结论：并发请求写入队列不是太好的方案，因为写入队列会耗费60%以上的时间，并且可能出现错误消费或消费超时的问题，只能用在日志这种低并发场景下。
*/