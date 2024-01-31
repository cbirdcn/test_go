# link连接服务

## 简介

基于`gorilla/websocket`提供websocket长连接和消息转发的服务。

接收wsClient的请求，并将消息体原样通过rpc协议转发到logic service，阻塞等待logic返回后再将返回作为响应消息原样返回给wsClient

所以在wsServer启动前，需要启动logicService

注意：rpc调用过程没有使用protobuf，传输内容是字符串

目前只实现了基础功能

## 说明

wsClient是客户端demo，每秒向wsServer发送pb类型的消息，当遇到中断或1秒内无法得到响应则退出

client_readme是对wsClient代码的解释说明
