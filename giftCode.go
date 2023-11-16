package main

import "fmt"

func main() {
    code_map := make(map[string]struct{})
    code_map["a"]=struct{}{}
    if _, ok := code_map["a"]; ok{
        fmt.Println("a exist.")
    }
}
