# 多进程的game服务器

## 简介

包含dbServer、logicServer、linkServer三个主要服务，以及一个测试程序linkClient。

dbServer提供数据库访问服务，支持rpc协议。具体实现在模块内部指定，这里仅连接mongodb

logicServer提供逻辑服务，支持rpc协议。外部接受linkServer的请求，内部将linkServer请求路由并重新rpc调用dbServer的方法，获取到数据或返回后打包响应linkServer的rpc调用。

linkServer提供长连接和消息中转服务，支持ws连接。外部接受linkClient的连接，并将字符串消息解码成pb支持的消息类型，并调用logicServer提供的rpc方法，再将返回值构建成字节数组返回给linkClient。

linkClient提供长连接的客户端测试功能，可以将pb支持的消息类型转码为proto编码过的字节数组并通过tcp发送到linkServer等待响应并解析显示返回值。

## 说明

### 执行方法

由于使用了`glog`捕获日志，所以如果要定制glog，可以在执行时添加参数。否则默认`logtostderr=false`并将日志写入`/tmp/`。

注意，重名文件运行时产生多个glog文件，但是软链接只会指向其中的一个。

正确的启动方式：进入指定项目文件夹运行

```shell
go run dbServer.go
go run logicServer.go
go run linkServer.go
go run linkClient.go
```

已废弃的单入口模式（通过参数启动对应的服务）

```shell
go run main.go db
go run main.go logic
go run main.go linkServer
go run main.go linkClient
```
