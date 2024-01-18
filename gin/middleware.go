package main

import (
	"fmt"
    "time"
	"github.com/gin-gonic/gin"
)

// 使用教程：https://www.topgoer.com/gin%E6%A1%86%E6%9E%B6/gin%E4%B8%AD%E9%97%B4%E4%BB%B6/next%E6%96%B9%E6%B3%95.html
// 其他中间件：https://www.topgoer.com/gin%E6%A1%86%E6%9E%B6/gin%E4%B8%AD%E9%97%B4%E4%BB%B6/%E4%B8%AD%E9%97%B4%E4%BB%B6%E6%8E%A8%E8%8D%90.html

// 全局中间件
func MiddleWareGlobal() gin.HandlerFunc {
    return func(c *gin.Context) {
        t := time.Now()
        fmt.Println("Global前置中间件开始执行")
        // 设置变量到Context的key中，可以通过Get()取
        c.Set("request", "中间件")
        // 前置中间件结束，开始执行函数
		fmt.Println("Global前置中间件结束")
        c.Next()
        // 后置中间件开始（执行完handleFunc后要做的工作）
		fmt.Println("Global后置中间件开始")
        status := c.Writer.Status()
        fmt.Println("Global后置中间件执行完毕", status)
        t2 := time.Since(t)
        fmt.Println("Global time:", t2)
    }
}

// 局部中间件
func MiddleWareLocal() gin.HandlerFunc {
    return func(c *gin.Context) {
        t := time.Now()
        fmt.Println("Local前置中间件开始执行")
        // 设置变量到Context的key中，可以通过Get()取
        c.Set("request", "中间件")
        // 前置中间件结束，开始执行函数
		fmt.Println("Local前置中间件结束")
        c.Next()
        // 后置中间件开始（执行完handleFunc后要做的工作）
		fmt.Println("Local后置中间件开始")
        status := c.Writer.Status()
        fmt.Println("后置中间件执行完毕", status)
        t2 := time.Since(t)
        fmt.Println("Local time:", t2)
    }
}

// 错误处理中间件
func MiddleWareError() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("Error前置中间件开始执行")
        // 延迟处理
        defer func() {
            // 通过recover获取错误信息
            if err := recover(); err != nil {
                // 统一以json格式返回错误信息
                c.JSON(500, gin.H{
                    "status": 500,
                    "error":  fmt.Sprint(err),
                })
            }
        }()
        status := c.Writer.Status()
        fmt.Println("Error前置中间件执行完毕", status)
        c.Next()
    }
}

func main() {
    // 1.创建路由
    // 默认使用了2个中间件Logger(), Recovery()
    r := gin.Default()
    // 注册全局中间件，对所有路由生效
    r.Use(MiddleWareGlobal())

    // {}为了代码规范
    {
        // curl http://localhost:8888/ce?request=1
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
    Global前置中间件开始执行
    Global前置中间件结束
    函数执行，获取request: 中间件
    Global后置中间件开始
    Global后置中间件执行完毕 200
    Global time: 96.193µs
	*/

	// 局部中间件
    // curl http://localhost:8888/ce2?request=2
    r.GET("/ce2", MiddleWareLocal(), func(c *gin.Context) {
        // 取值
        req, _ := c.Get("request")
        fmt.Println("request:", req)
        // 页面接收
        c.JSON(200, gin.H{"request": req})
    })
	// 输出：相当于两次经过同样的中间件，把global和local理解为栈即可。前置相当于入栈，global在前local在后。后置相当于出栈，local在前，global在后。
	/*
    Global前置中间件开始执行
    Global前置中间件结束
    Local前置中间件开始执行
    Local前置中间件结束
    request: 中间件
    Local后置中间件开始
    后置中间件执行完毕 200
    Local time: 87.926µs
    Global后置中间件开始
    Global后置中间件执行完毕 200
    Global time: 316.255µs
	*/

    // 错误处理中间件
    // curl http://localhost:8888/error
    r.GET("/error", MiddleWareError(), func(c *gin.Context) {
        panic("error occurred")
    })
    /*
    Global前置中间件开始执行
    Global前置中间件结束
    Error前置中间件开始执行
    Error前置中间件执行完毕 200
    Global后置中间件开始
    Global后置中间件执行完毕 500
    Global time: 358.711µs
    */

    r.Run(":8888")
}