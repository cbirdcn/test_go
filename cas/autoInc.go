package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

/* 使用自旋锁+atomic.CompareAndSwapInt32()为并发请求分配自增ID*/

var (
    counter int32          //计数器
    wg      sync.WaitGroup
	sm sync.Map
)

func main() {
    threadNum := 1000
    wg.Add(threadNum)
    for i := 0; i < threadNum; i++ {
        go incCounter(i)
    }
    wg.Wait()


	kset := make(map[int]struct{}) // map中的key集合
	vset := make(map[int32]struct{}) // map中的val集合
	
	fmt.Printf("final counter:%d\n", counter)
	fmt.Println("final map is:")
	sm.Range(func(key, value interface{}) bool {
		fmt.Printf("key=%d,val=%d\n", key, value)
		kset[key.(int)] = struct{}{}
		vset[value.(int32)] = struct{}{}
		return true
	})
	fmt.Println("key set len is:", len(kset))
	fmt.Println("val set len is:", len(vset))

}

func incCounter(index int) {
    defer wg.Done()

	// 逻辑：自旋锁中持续CAS，修改成功就把index=>目标值的关系存入sync.Map，然后break。不成功就进入下一循环
    for {
        // 原子操作
        old := counter
        ok := atomic.CompareAndSwapInt32(&counter, old, old+1)
        if ok {
			sm.Store(index, old+1)
            break
        } 
    }
    
}

/*
输出：
final counter:1000
final map is:
key=26,val=13
...
key=864,val=877
key=948,val=963
key set len is: 1000
val set len is: 1000
*/