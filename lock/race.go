package main

import (
	"fmt"
	"time"
)

var balance int

//存款
func Deposit(amount int) {
	balance = balance + amount
}

//读取余额
func Balance() int {
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

	time.Sleep(time.Second) // 为了简单，不使用waitGroup
	fmt.Println(balance)
}
