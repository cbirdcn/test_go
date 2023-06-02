package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// https://geektutu.com/post/quick-go-gin.html

func main() {
	r := gin.Default() // 生成一个应用程序实例

	// 声明一个路由，以及触发的函数
	// curl http://localhost:8081/
	r.GET("/", func(c *gin.Context) {
		// 函数内返回客户端想要的响应
		c.String(http.StatusOK, "hello world.")
	})

	// 路由方法有 GET, POST, PUT, PATCH, DELETE 和 OPTIONS，还有Any，可匹配以上任意类型的请求。

	// 带参数的URL。:name表示传入不同的 name。/user/:name/*role，* 代表可选。
	// curl http://localhost:8081/user/geektutu
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	// 获取Query参数
	// 匹配users?name=xxx&role=xxx，role可选
	// 注意query和param格式的不同
	// curl "http://localhost:8081/users?name=Tom&role=student"
	r.GET("/users", func(c *gin.Context) {
		name := c.Query("name")
		role := c.DefaultQuery("role", "teacher") // 为query设置默认值
		c.String(http.StatusOK, "%s is a %s", name, role)
	})

	// POST
	// curl http://localhost:8081/form  -X POST -d 'username=geektutu&password=1234'
	r.POST("/form", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.DefaultPostForm("password", "000000") // 可设置默认值

		// c.JSON()返回JSON消息
		// gin.H{} is a shortcut for map[string]interface{}
		// 注意，会被自动排序
		c.JSON(http.StatusOK, gin.H{
			"username": username,
			"password": password,
		})
	})

	// GET 和 POST 混合
	// curl "http://localhost:8081/posts?id=9876&page=7"  -X POST -d 'username=geektutu&password=1234'
	r.POST("/posts", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		username := c.PostForm("username")
		password := c.DefaultPostForm("username", "000000") // 可设置默认值

		c.JSON(http.StatusOK, gin.H{
			"id":       id,
			"page":     page,
			"username": username,
			"password": password,
		})
	})

	// Map参数(字典参数)
	// curl -g "http://localhost:8081/post?ids[Jack]=001&ids[Tom]=002" -X POST -d 'names[a]=Sam&names[b]=David'
	// {"ids":{"Jack":"001","Tom":"002"},"names":{"a":"Sam","b":"David"}}
	r.POST("/post", func(c *gin.Context) {
		ids := c.QueryMap("ids")
		names := c.PostFormMap("names")

		c.JSON(http.StatusOK, gin.H{
			"ids":   ids,
			"names": names,
		})
	})

	// 重定向(Redirect)
	// https://pkg.go.dev/gopkg.in/gin-gonic/gin.v1#section-readme
	// curl -i http://localhost:8081/redirect
	// curl "http://localhost:8081/reindex"
	// Issuing a HTTP redirect is easy. Both internal and external locations are supported.
	r.GET("/redirect", func(c *gin.Context) {
		// http重定向
		c.Redirect(http.StatusMovedPermanently, "/reindex")
	})
	// Issuing a Router redirect, use HandleContext like below.
	r.GET("/reindex", func(c *gin.Context) {
		c.Request.URL.Path = "/" // 路由重定向，目标是"/"
		r.HandleContext(c)       // 执行调用
	})

	// 分组路由(Grouping Routes)
	// curl http://localhost:8081/v1/posts
	// curl http://localhost:8081/v2/posts
	defaultHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"path": c.FullPath(),
		})
	}
	// group: v1
	v1 := r.Group("/v1")
	{
		v1.GET("/posts", defaultHandler)
		v1.GET("/series", defaultHandler)
	}
	// group: v2
	v2 := r.Group("/v2")
	{
		v2.GET("/posts", defaultHandler)
		v2.GET("/series", defaultHandler)
	}

	// 上传文件
	// 单个文件
	r.POST("/upload1", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		// c.SaveUploadedFile(file, dst)
		c.String(http.StatusOK, "%s uploaded!", file.Filename)
	})
	// 多文件
	r.POST("/upload2", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)
			// c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, "%d files uploaded!", len(files))
	})

	// HTML模板
	// curl http://localhost:8081/arr
	// Gin默认使用模板Go语言标准库的模板text/template和html/template
	// 标准库：https://pkg.go.dev/text/template
	// https://pkg.go.dev/html/template
	type student struct {
		Name string
		Age  int8
	}
	r.LoadHTMLGlob("templates/*") // LoadHTMLGlob loads HTML files identified by glob pattern and associates the result with HTML renderer.
	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/arr", func(c *gin.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gin.H{
			"title":     "Gin",
			"stuArr":    [2]*student{stu1, stu2},  // 写法1：[2]*student{stu1, stu2}表示长度=2的slice，slice内每个元素都是*student指针类型
			"stuValArr": [2]student{*stu1, *stu2}, // 写法2：[2]student{*stu1, *stu2}。两者模板写法一致，原因应该是模板编译过程中，自动把写法1的变量var编译成了*var，不管指针还是值最终都是拿到值的拷贝了
		})
	})

	r.Run(":8081") //让应用运行在一个本地服务器上，监听端口默认为8080
}
