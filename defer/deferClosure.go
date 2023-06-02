package main

// defer在匿名函数与闭包中

/*
闭包与匿名函数

匿名函数:没有函数名的函数。

闭包:可以使用另外一个函数作用域中的变量的函数。闭包把函数和运行时的引用环境打包成为一个新的整体，当每次调用包含闭包的函数时都将返回一个新的闭包实例，这些实例之间是隔离的，分别包含调用时不同的引用环境现场。
*/

import (
    "fmt"
)

func main() {
    simpleClosure()
    fmt.Println(".....")
    closureWithParam()
}

func simpleClosure() {
        for i := 0; i < 5; i++ {
            defer func() {
                fmt.Println(i)
                // 输出:5 5 5 5 5
                // 因为,defer 表达式中的 i 是对 for 循环中 i 的引用。循环结束，i 加到 5。退出循环，main结束前，执行defer，最后全部打印 5。
            }()
        }
    }

func closureWithParam() {
        for i := 0; i < 5; i++ {
            defer func(i int) {
                fmt.Println(i)
                // 输出:4 3 2 1 0
                // 因为,将 i 作为参数传入 defer 表达式中，在传入最初就会进行求值保存，只是没有执行延迟函数而已。
            }(i)
        }
    }
