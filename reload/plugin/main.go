package main

import (
    "fmt"
    "os"
    "plugin"
)

func main() {
    p, err := plugin.Open("./plugina.so")
    if err != nil {
        fmt.Println("error open plugin: ", err)
        os.Exit(-1)
    }
    s, err := p.Lookup("IamPluginA")
    if err != nil {
        fmt.Println("error lookup IamPluginA: ", err)
        os.Exit(-1)
    }
    if x, ok := s.(func()); ok {
        x()
    }
}