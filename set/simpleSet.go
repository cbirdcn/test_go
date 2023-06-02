package main

import "fmt"

type void struct{}

var val void

func main() {
	set := make(map[interface{}]void)

	// 集合赋值
	set["a"] = val

	// 输出
	for k := range set {
		fmt.Printf("%v\n", k)
	}

	// 集合删除元素及长度查询
	delete(set, "a")
	size := len(set)
	fmt.Println(size)
	
	// 判断集合中是否存在某元素
	_, ok := set["a"]
	fmt.Println(ok)
}
