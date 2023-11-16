package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 检查每个amqp调用的返回值
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// 连接到 RabbitMQ 服务器
	conn, err := amqp.Dial("amqp://admin:admin@host.docker.internal:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 创建一个通道，用于完成任务的大多数API都驻留在通道中
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 要发送，我们必须声明一个要发送到的队列
	// 声明一个队列是幂等的——只有当它不存在时才会创建它。
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	// 将消息发布到声明的队列
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),	// 消息内容是一个字节数组
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
