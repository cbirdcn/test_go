package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"sync"
)

// kafka consumer

func main() {
	// BlockConsume()
	UnblockConsume()
}

func BlockConsume() {
	consumer, err := sarama.NewConsumer([]string{"host.docker.internal:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions("web_log") // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println("分区列表", partitionList)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest) // OffsetOldest表示取过去的消息，OffsetNewest取最新消息也就是还没产生的消息随产生随消费
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
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
	partitionList, err := consumer.Partitions("web_log") // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println("分区列表", partitionList)
	for partition := range partitionList { // 遍历所有的分区
		wg.Add(1)
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest) // OffsetOldest表示取过去的消息，OffsetNewest取最新消息也就是还没产生的消息随产生随消费
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