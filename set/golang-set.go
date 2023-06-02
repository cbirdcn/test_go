package main

// golang-set使用map和空value实现set，封装了初始化、添加、包含、长度、清除(重新初始化)、set相等（遍历两个set的每个元素）、子集（遍历判断包含）等很多功能
// 区分线程不安全（默认）和安全两个版本

import (
    "fmt"
    mapset "github.com/deckarep/golang-set"
)

func main() {
    // 默认创建的线程安全的，如果无需线程安全
    // 可以使用 NewThreadUnsafeSet 创建，使用方法都是一样的。
    s1 := mapset.NewSet(1, 2, 3, 4)  
    fmt.Println("s1 contains 3: ", s1.Contains(3))
    fmt.Println("s1 contains 5: ", s1.Contains(5))
    fmt.Println("s1 length: ", s1.Cardinality()) // 基数（集合长度）

    // interface 参数，可以传递任意类型
    s1.Add("poloxue")
    fmt.Println("s1 contains poloxue: ", s1.Contains("poloxue"))
    s1.Remove(3)
    fmt.Println("s1 contains 3: ", s1.Contains(3))

    s2 := mapset.NewSet(1, 3, 4, 5)

    // 并集
    fmt.Println(s1.Union(s2))
}