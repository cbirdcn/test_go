# channel

## emptyChannel.go

用struct{}{}填充的Channel可以省内存。尤其在事件通信时，channel不发数据，就能通知其他协程。
