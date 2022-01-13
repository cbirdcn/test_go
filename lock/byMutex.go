package main

import (
    "fmt"
    "time"
    "sync"
)


var balance int
var mu sync.Mutex

//存款
func Deposit(amount int) { 
    mu.Lock() 
    defer mu.Unlock()
    balance = balance + amount
}
//读取余额
func Balance() int { 
    mu.Lock() 
    defer mu.Unlock()
    return balance
}

func main(){
    //小王：存600，并读取余额
    go func(){
        Deposit(600)
        fmt.Println(Balance())
    }()
    //小张：存500
    go func(){
        Deposit(500)
		fmt.Println(Balance())
    }()
    
    time.Sleep(time.Second)
    fmt.Println(Balance()) // 外部也要加锁，因为主协程没有加waitGroup，所以可能和子协程并发
}
