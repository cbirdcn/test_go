package main

import (
    "fmt"
)

func IamPluginA() {
    fmt.Println("Hello, I am PluginA!")
}

// 用新源码编译好的so替换旧so
// go build --buildmode=plugin -o plugina.so plugina.go