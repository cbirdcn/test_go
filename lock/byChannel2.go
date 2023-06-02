package main

import "fmt"

// 问题：不同的操作需要不同的channel处理。比方说余额操作就需要增、减、查三个信道。对于操作更多的业务很不友好
var deposits = make(chan int) // send amout to deposit
var balances = make(chan int) // receive balance

var balance int // balance 只在 teller 中可以访问

func Deposit(amount int) {
    deposits <- amount
}

func Balance() int {
    return <-balances
}

func teller() {
	// 问题：循环阻塞可能导致死锁
    for {
        select {
        case amount := <-deposits:
            balance += amount
        case balances <- balance:
        }
    }
}

func init() {
    go teller()
}

func main() {
	go Deposit(100)
	fmt.Println(balance)
	go Balance()
	fmt.Println(balance)
}