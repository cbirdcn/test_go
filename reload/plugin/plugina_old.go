package main

import (
    "fmt"
)

func IamPluginA() {
    fmt.Println("Hello, I am PluginA Old!")
}

// go build --buildmode=plugin -o plugina_old.so plugina_old.go