package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

var (
    counter int32          //计数器
    wg      sync.WaitGroup //信号量
)

func main() {
    threadNum := 5000
    wg.Add(threadNum)
    for i := 0; i < threadNum; i++ {
        go incCounter(i)
    }
    wg.Wait()
	fmt.Printf("final counter:%d\n", counter)
}

func incCounter(index int) {
    defer wg.Done()

	// 逻辑：自旋锁中持续CAS，修改成功就break，不成功就进入下一循环（计数+1）	
    spinNum := 0 // 自旋失败计数
    for {
        // 原子操作
        old := counter
        ok := atomic.CompareAndSwapInt32(&counter, old, old+1)
        if ok {
            break
        } else {
            spinNum++
        }
    }
	if spinNum > 0 {
		fmt.Printf("thread,%d,spinnum,%d\n", index, spinNum)
	}
    
}