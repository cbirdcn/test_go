package main

import (
	"fmt"
    "time"
	"github.com/gin-gonic/gin"
)

// 使用教程：https://www.topgoer.com/gin%E6%A1%86%E6%9E%B6/gin%E4%B8%AD%E9%97%B4%E4%BB%B6/next%E6%96%B9%E6%B3%95.html
// 其他中间件：https://www.topgoer.com/gin%E6%A1%86%E6%9E%B6/gin%E4%B8%AD%E9%97%B4%E4%BB%B6/%E4%B8%AD%E9%97%B4%E4%BB%B6%E6%8E%A8%E8%8D%90.html

// 定义中间件
func MiddleWare() gin.HandlerFunc {
    return func(c *gin.Context) {
        t := time.Now()
        fmt.Println("前置中间件开始执行")
        // 设置变量到Context的key中，可以通过Get()取
        c.Set("request", "中间件")
        // 前置中间件结束，开始执行函数
		fmt.Println("前置中间件结束")
        c.Next()
        // 后置中间件开始（执行完handleFunc后要做的工作）
		fmt.Println("后置中间件开始")
        status := c.Writer.Status()
        fmt.Println("后置中间件执行完毕", status)
        t2 := time.Since(t)
        fmt.Println("time:", t2)
    }
}

func main() {
    // 1.创建路由
    // 默认使用了2个中间件Logger(), Recovery()
    r := gin.Default()
    // 注册全局中间件
    r.Use(MiddleWare())
    // {}为了代码规范
    {
        r.GET("/ce", func(c *gin.Context) {
            // 取值
            req, _ := c.Get("request")
            fmt.Println("函数执行，获取request:", req)
            // 页面接收
            c.JSON(200, gin.H{"request": req})
        })

    }
	// 打印：
	/*
	前置中间件开始执行
	前置中间件结束
	函数执行，获取request: 中间件
	后置中间件开始
	后置中间件执行完毕 200
	*/

	// 局部中间件
    r.GET("/ce2", MiddleWare(), func(c *gin.Context) {
        // 取值
        req, _ := c.Get("request")
        fmt.Println("request:", req)
        // 页面接收
        c.JSON(200, gin.H{"request": req})
    })
	// 输出：相当于两次经过同样的中间件
	/*
	前置中间件开始执行
	前置中间件结束
	前置中间件开始执行
	前置中间件结束
	request: 中间件
	后置中间件开始
	后置中间件执行完毕 200
	time: 39.598µs
	后置中间件开始
	后置中间件执行完毕 200
	time: 60.384µs
	*/

    r.Run(":8082")
}