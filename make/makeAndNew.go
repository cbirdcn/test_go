package main

import "fmt"

type myChan chan int

func main() {
	var1 := make(chan int)
	var2 := make(myChan)
	var3 := new(myChan)

	fmt.Printf("%T, %T, %T\n", var1, var2, var3)
}