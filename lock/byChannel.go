package main

import (
	"fmt"
	"time"
)

var balance int

//存款
func Deposit(ch chan bool, amount int) {
	ch <- true // 抢channel，没抢到的被阻塞
	balance = balance + amount
	<-ch // 释放channel，其他被阻塞的可以继续抢
}

//读取余额
func Balance(ch chan bool) int {
	ch <- true
	ret := balance
	<-ch
	return ret
}

func main() {
	ch := make(chan bool, 1)
	//小王：存600，并读取余额
	go func() {
		Deposit(ch, 600)
		fmt.Println(Balance(ch))
	}()
	//小张：存500
	go func() {
		Deposit(ch, 500)
		fmt.Println(Balance(ch))
	}()

	// 为了简单，不使用 sync.WaitGroup ， 所以需要defer和Sleep
	defer close(ch) // 记得关闭channel。单一程序是可以自动被GC清理掉channel的。但是如果goroutine很多，机器性能差，长时间积累或短时间爆发会压垮机器的CPU和内存，手动关闭channel是个好习惯。

	time.Sleep(time.Second)
	fmt.Println(balance)
}