package main

import (
    "fmt"
    "sort"
)

func main(){
    main2()
}

func main1() {
    var m = map[string]int{
        "hello":         0,
        "morning":       1,
        "keke":          2,
        "jame":          3,
    }
	// fmt.Println(m) // 打印得到固定顺序：map[hello:0 jame:3 keke:2 morning:1]
    var keys []string
    for k := range m { // 忽略v，把k入slice
		// fmt.Println(k) // 从map读数据是无序的
        keys = append(keys, k)
    }
	// fmt.Println(keys) // 打印会得到不同的顺序，比如：[jame hello morning keke] 
    sort.Strings(keys) // 对key slice 原地排序
	// fmt.Println(keys) // 打印得到固定顺序：[hello jame keke morning]
    for _, k := range keys {
        fmt.Println("Key:", k, "Value:", m[k])
    }
}


func main2() {
    /* 声明索引类型为字符串的map */
    var testMap = make(map[string]string)
    testMap["Bda"] = "B"
    testMap["Ada"] = "A"
    testMap["Dda"] = "D"
    testMap["Cda"] = "C"
    testMap["Eda"] = "E"

    for key, value := range testMap {
        fmt.Println(key, ":", value)
    }
    var testSlice []string
    testSlice = append(testSlice, "Bda", "Ada", "Dda", "Cda", "Eda")

    /* 对slice数组进行排序，然后就可以根据key值顺序读取map */
    sort.Strings(testSlice)
    fmt.Println("排序输出:")
    for _, Key := range testSlice {
        /* 按顺序从MAP中取值输出 */
        if Value, ok := testMap[Key]; ok {
            fmt.Println(Key, ":", Value)
        }
    }
}