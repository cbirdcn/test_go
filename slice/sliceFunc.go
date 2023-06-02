package main

import "fmt"

func change(s ...string) {
	s[0] = "Go"
	s = append(s, "playground")
	fmt.Println(s)
	// 修改：返回append后新创建的临时数组对应的slice
	// return s
}

func main() {
	welcome := []string{"hello", "world"}
	// 修改：如果要影响到当前welcome，需要welcome = change(...)
	change(welcome...)
	fmt.Println(welcome)
}

// 原输出：
// [Go world playground]
// [Go world]
// 修改后输出：
// [Go world playground]
// [Go world playground]