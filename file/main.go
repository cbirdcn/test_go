package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
    file, err := os.OpenFile("./a.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY , 0666)
    defer file.Close()
    if err != nil {
        fmt.Println("open file failed", err)
    }
    writer := bufio.NewWriter(file)
    _,err = writer.WriteString("123")
    if err != nil {
        fmt.Println("write failed", err)
    } else {
	err = writer.Flush()
	if err != nil {
		panic(err)
	} else {
        	fmt.Println("write success")
	}
    }
}
