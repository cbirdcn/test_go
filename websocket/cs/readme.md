# cs模式的websocket

websocket是一种在单个TCP连接上进行全双工通信的协议。允许服务端主动向客户端推送数据。

只需要完成一次握手，CS之间就直接可以创建持久性的连接，并进行双向数据传输。

websocket协议通过原始HTTP连接起作用。

## server

使用`net/http`对外提供http连接

在`http router`上注册路由`/ws`及处理函数`wsHandler`。

在`wsHandler`内，使用`websocket.Upgrader`将每个到来的连接从http协议升级成ws协议（内部逻辑不管）。然后陷入for死循环（event loop的简单实现）。在循环中读取请求消息conn.ReadMessage。再将响应消息写入连接conn.WriteMessage。

注意，只是在route为`/ws`上的请求才会被提升为ws协议。其他http路由不受影响。

注意，ws和tcp一样，读取和写入到连接的数据都是字节数组。

## client的go实现

客户端可以给拨号器加自定义配置，或直接用默认拨号器`websocket.DefaultDialer.Dial`连接到服务器`"ws://host:port" + "/route"`

然后并发启动一个协程`receiveHandler`接收服务器下发的消息。在协程内陷入for死循环（event loop的简单实现），不断从连接读取消息conn.ReadMessage，并利用消息比如打印。

注意，客户端也可以用其他语言。

## 回声测试

在client的main中，有一段for循环阻塞。在里面，select监听以下事件就绪：

- 时间过了1秒：`time.After(time.Duration(1) * time.Millisecond * 1000)`
- 捕获了中断信号：`os.Interrupt`
  - 客户端主动发送关闭连接的消息：`websocket.FormatCloseMessage`
  - 由于`receiveHandler`在异步接收服务器下发的消息，所以需要等待此Handler完成或达到超时时间

注意，这里使用了`已关闭channel能读取出零值`的特性，如果异步进行的`receiveHandler`耗时比倒计时器短就表示`done channel`被关闭了，达成select的channel关闭条件。否则select将阻塞到倒计时器截止才能解开。

## 参考

[GoLang中使用Gorilla Websocket](https://juejin.cn/post/6946952376825675812)
