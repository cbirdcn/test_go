# buntdb服务

## 方案1：rpc 服务

gate_server 作为客户端访问的服务器端，可以提供 http 或 tcp 服务，这里简单用 gin 提供 http 服务。

协议类型会影响响应的方式， http 需要同步阻塞直到有操作结果时（也就是事件就绪）才返回， tcp 可以异步地等操作成功后再向客户端主动下发操作结果。

gate_server 只是承担高并发、路由、请求参数检查和响应等工作，还可以包含部分 gs 逻辑（也可以分给 rpc server 处理），不会直接操作 DB 。

gate_server 的 gs 逻辑中如果需要操作DB时，可以作为 rpc 客户端连接到已经运行的 db rpc server ，将不同类型的请求（读、写、不同类型）按照约定好的格式通过 protobuf 提交过去并得到同步的响应或失败原因。

再将结果返回给 http 客户端。

特点：
- gate_server 能抗住高并发
- 低耦合
- 需要 rpc server 提供方控制数据访问频率，避免高并发对DB造成过大性能负担，打垮DB。
- 无法保证请求的顺序性，需要客户端辅助保证，或者添加唯一消息id

## 方案2：消息队列

gate_server 作为 REQ 队列的消息生产者和 RES 队列的消息消费者。 db_server 作为 REQ 队列的消息消费者和 RES 队列的消息生产者。他们互相持有对方的资源。

消费者需要用 goroutine 并发启动并陷入阻塞，在此期间不停地消费，所以需要提前启动。注意，消息只能被消费一次。

db_server 启动的 REQ 消费队列会持续不断地消费 gate_server 发送过来的 REQ 队列消息，同步访问DB并将结果通过 RES 队列消息发送回 gate_server。

gate_server 把 DB 操作请求 json 序列化后交给 REQ 队列。具体来说，gate_server 对 DB 发送请求消息前会生成全局唯一的消息队列操作id（雪花id+user_id组合保证唯一），并初始化一个全局map的value（channel）。发送 REQ 后通过读取map的value陷入阻塞，map的key是当前唯一的消息队列操作id，value是一个channel。这个channel什么时候被写入呢？异步地，每当 RES 消息返回并被 gate_server 消费时，都会将key和value写入到全局map中。这时发送REQ消息后陷入阻塞的流程就能解开channel阻塞，并将value返回给http 客户端。另外根据需要，有必要将map中的kv立即销毁或隔一段时间销毁最旧。

注意几点细节：
- consumer需要提前启动，并处于阻塞循环中
- REQ 和 RES 需要根据唯一id才能匹配
- map的value是channel，需要在发送前就初始化，因为发送后就要立即读取了
- 如果迟迟等不到RES消息，http server需要作出超时响应给http client

特点：
- 可以将请求序列化，消息是有序且有记录的
- 两个队列对时间、性能的消耗有点大
- 对并发支持能力更弱了
- 适合作为日志收集器，而不是并发处理器

## 优化

服务应该用tcp，至少也应该是ws

消息id应该和userid+time+rand有关

消息内容避免用json，最好是定义好的结构或protobuf、msgpack等