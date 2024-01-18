package main

import (
	"log"
	"fmt"
	"time"
	"github.com/tidwall/buntdb"
	"github.com/Shopify/sarama"
	"sync"
	"encoding/json"
)

var db *buntdb.DB
var err error
var b_config buntdb.Config
var k_config *sarama.Config
var k_client sarama.SyncProducer

const TOPIC_REQ string = "req"
const TOPIC_RES string = "res"

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

func init() {
	// 作为REQ的服务端，初始化db客户端操作db
	db, err = buntdb.Open("data.db")
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close() // TODO:只有在遇到中断信号时才能close

	if err := db.ReadConfig(&b_config); err != nil{
		log.Fatal(err)
	}

	// 作为RES队列的客户端
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

// kafka consumer

func main() {
	fmt.Println("监听req队列中...")
	BlockConsume()
	// UnblockConsume()
}

func SendMessage(topic string, value string) {
	t := time.Now().UnixMilli()
	msg := &sarama.ProducerMessage{}
	msg.Topic = TOPIC_RES
	msg.Value = sarama.StringEncoder(value)

	// 发送消息
	pid, offset, err := k_client.SendMessage(msg)
	if err != nil {
		fmt.Println("send res msg failed, err:", err)
		return
	}
	fmt.Printf("res partition_id:%v offset:%v\n", pid, offset)
	fmt.Println("db server send res cost ms:", time.Now().UnixMilli()-t) // 回写操作队列时间不稳定，100并发耗时在4-57ms之间
}

func BlockConsume() {
	consumer, err := sarama.NewConsumer([]string{"host.docker.internal:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(TOPIC_REQ) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(TOPIC_REQ, int32(partition), sarama.OffsetNewest) // OffsetOldest表示取过去的消息，OffsetNewest取最新消息也就是还没产生的消息随产生随消费
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				go func() {
					fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
					req := REQ{}
					err := json.Unmarshal(msg.Value, &req)
					//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
					if err != nil {
						fmt.Println(err)
					}
					if req.Op == "set" {
						t := time.Now().UnixMilli()
						err = db.Update(func(tx *buntdb.Tx) error {
							_, _, err := tx.Set(req.Key, req.Value, nil)
							return err
						})
						fmt.Println("db op cost ms:", time.Now().UnixMilli()-t) // db操作时间不稳定，100并发时耗时在0-17ms之间
						// TODO:判断逻辑，如果发生错误要回滚并通知请求方
						
						// 操作响应队列：通过唯一id回写状态、原因、返回值等数据，响应队列的消费者负责保持tcp连接并响应最原始的客户端请求。
						res := RES{
							UserId: req.UserId,
							MsgId: req.MsgId,
							Status: 200,
							Reason: "ok",
							Data: "{}",
						}
						json, _ := json.Marshal(res)
						SendMessage(TOPIC_RES, string(json))
					}
				}()
			}
		}(pc)
	}
	select{} //阻塞进程
}

func UnblockConsume() {
	var wg sync.WaitGroup
	consumer, err := sarama.NewConsumer([]string{"host.docker.internal:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(TOPIC_REQ) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println("分区列表", partitionList)
	for partition := range partitionList { // 遍历所有的分区
		wg.Add(1)
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(TOPIC_REQ, int32(partition), sarama.OffsetNewest) // OffsetOldest表示取过去的消息，OffsetNewest取最新消息也就是还没产生的消息随产生随消费
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		// 异步从每个分区消费信息
		go func(pc sarama.PartitionConsumer, wg *sync.WaitGroup) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
			}
			defer pc.AsyncClose()
			// wg.Done()
		}(pc, &wg)

		defer wg.Done()
	}
	wg.Wait()
	consumer.Close()
}

/*
输出：
分区列表 [0]

假设先启动producer生产了一条数据，然后才启动consumer并陷入阻塞，然后再次启动producer生产新的数据

如果是取OffsetOldest，此时会显示出过去消息。然后再随着生产者继续生产而继续消费。
Partition:0 Offset:0 Key:[] Value:this is a test log
Partition:0 Offset:1 Key:[] Value:this is a test log

如果是取OffsetNewest，此时不会显示过去消息，而是随着生产者的生产才开始消费
Partition:0 Offset:2 Key:[] Value:this is a test log

*/