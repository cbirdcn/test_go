package main

import "fmt"

func worker(ch chan struct{}) {
	ch <- struct{}{}
	fmt.Println("do something")
	close(ch)
}

func main() {
	ch := make(chan struct{})
	go worker(ch)
	<-ch
}