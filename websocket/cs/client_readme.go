package main

import (
    "time"
    "fmt"
)

var done chan interface{}

func receiveHandler() {
    defer close(done)
    for {
        fmt.Println("receive message")
        time.Sleep(time.Duration(1) * time.Millisecond)
		// return
    }
}

func main() {
    done = make(chan interface{})
    go receiveHandler()
    for {
        select {
        case <-done:
			// 接收处理器已关闭Channel。原因是正常情况下没有给channel写入数据，所以无法读取数据。但是当channel关闭时，就能从channel读取数据了，是零值。
            fmt.Println("Receiver Channel Closed! Exiting...")
        case <-time.After(time.Duration(2) * time.Millisecond):
			// 达到超时时间
            fmt.Println("Timeout! Exiting...")
        }
        break // 退出
    }
    fmt.Println("end")
}


/*
注释receiveHandler的for中的return，输出：

receive message
receive message
Timeout! Exiting...
end

或

receive message
receive message
receive message
Timeout! Exiting...
end

由于goroutine可能在倒计时器之前启动，所以可能是2-3条消息
*/

/*
如果放开receiveHandler的for中的return注释，输出：

receive message
Receiver Channel Closed! Exiting...
end

由于receiveHandler在倒计时器截止前就return了，也就是channel被关闭了。select中从已关闭channel读取到了零值，不再阻塞
*/