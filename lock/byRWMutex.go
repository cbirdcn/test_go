package main

import (
	"fmt"
	"sync"
	"time"
)

var balance int
var mu sync.RWMutex // 读锁

//存款
func Deposit(amount int) {
	mu.Lock()
	defer mu.Unlock()
	balance = balance + amount
}

//读取余额
func Balance() int {
	mu.RLock() // 读取余额要获取读锁
	defer mu.RUnlock()
	return balance
}

func main() {
	//小王：存600，并读取余额
	go func() {
		Deposit(600)
		fmt.Println(Balance())
	}()
	//小张：存500
	go func() {
		Deposit(500)
		fmt.Println(Balance())
	}()

	time.Sleep(time.Second)
	fmt.Println(balance)
}
