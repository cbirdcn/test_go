package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

var sqlDB *sql.DB
var gormDB *gorm.DB
var db *sql.DB

func init() {
	//host := "host.docker.internal" // 假设在容器中运行，访问其他ip段的容器映射到宿主机的服务
	host := "127.0.0.1"
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", "root", "root", host, 3306, "global")

	// database/sql创建数据库连接
	sqlDB, err = sql.Open("mysql", dsn)
	if !ExistError(err) {
		// GORM 允许通过一个现有的数据库连接来初始化 *gorm.DB
		gormDB, err = gorm.Open(mysql.New(mysql.Config{
			Conn: sqlDB,
		}), &gorm.Config{})

		if !ExistError(err) {
			db, err = gormDB.DB()
			if !ExistError(err) {
				db.SetMaxIdleConns(10) // 设置空闲连接池中连接的最大数量
				db.SetMaxOpenConns(100) // 设置打开数据库连接的最大数量
				db.SetConnMaxLifetime(time.Hour) // 设置了连接可复用的最大时间
				fmt.Println("mysql connection ok.")
			}

			err = gormDB.AutoMigrate(&originalModel{}) // gormDB有迁移工具，init时根据struct自动迁移。启动服务后可以检查original_models表的status字段
			ExistError(err)
		}
	}

}


func main() {
	router := gin.Default()
	v1 := router.Group("/api/v1/")

	{
		v1.GET("/", getAll)
		v1.GET("/:id", getOne)
		v1.POST("/", create)
		v1.PUT("/:id", updateOne)
		v1.DELETE("/:id", deleteOne)
	}

	router.Run(":8088")
}

// 两种Model结构体：第一个结构体代表原始的数据库字段，第二个结构体用来定义向 api 返回的字段。我们之所以在第二个结构体中重新定义返回的字段主要考虑到数据库中数据的安全性，我们不希望将数据库中的原始字段名
type (
	originalModel struct {
		gorm.Model // 把 ID，CreatedAt，UpdatedAt 和 DeletedAt 这四个字段嵌入到我们定义好的 todoModel 结构体中，一般数据表中都会用到这四个字段。
		Title string `json:"title"`
		Result int `json:"result"` // 结果状态（int）
	}

	transformedData struct {
		ID uint `json:"id"`
		Title string `json:"title"`
		Completed bool `json:"completed"` // 结果状态（bool）
	}
)

/*
curl --location --request GET 'http://127.0.0.1:8088/api/v1/'
*/
func getAll(c *gin.Context) {
	var originalDataList []originalModel
	gormDB.Find(&originalDataList) // 不加条件和分页
	if len(originalDataList) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "item not found",
		})
		return
	}

	// 对数据库的数据优化后返回
	var transformedDataList []transformedData
	for _,v := range originalDataList {
		completed := false
		if v.Result == 1 {
			completed = true
		}
		transformedDataList = append(transformedDataList, transformedData{
			ID: v.ID,
			Title:     v.Title,
			Completed: completed,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status" : http.StatusOK,
		"data" : transformedDataList,
	})
}

/*
curl --location --request GET 'http://127.0.0.1:8088/api/v1/1'
*/
func getOne(c *gin.Context) {
	var originalData originalModel
	id := c.Param("id")
	gormDB.First(&originalData, id)

	if originalData.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "item not found",
		})
		return
	}

	completed := false
	if originalData.Result == 1 {
		completed = true
	}

	transformedData :=  transformedData{
		ID:        originalData.ID,
		Title:     originalData.Title,
		Completed: completed,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": transformedData,
	})

}

/*
请求：
curl --location --request POST 'http://127.0.0.1:8088/api/v1' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'title=1' \
--data-urlencode 'completed=true'
*/
// 增
func create(c *gin.Context) {
	// gin.Context可以传递上下文，比如post请求的参数
	completed := c.PostForm("completed") // 从PostForm拿到字符串"true"
	completedBool,_ := strconv.ParseBool(completed) // 将字符串转为bool
	result := 0 // 将bool转成int
	if completedBool {
		result = 1
	}
	// 组装入库数据
	data := originalModel{
		Title:  c.PostForm("title"),
		Result: result, // 存储int型
	}
	gormDB.Save(&data) // orm save
	// 返回，注意rest api的http code
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"message": "item created successfully",
		"resourceId": data.ID, // 返回ID，由于orm save时使用了data地址，所以保存后data生成了ID
	})
}

// 对http put请求的响应，解释如下：
// 使用惯例是，在 PUT 请求中进行资源更新，但是不需要改变当前展示给用户的页面，那么返回 204 No Content。如果创建了资源，则返回 201 Created 。如果应将页面更改为新更新的页面，则应改用 200 。
// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Status/204
/*
curl --location --request PUT 'http://127.0.0.1:8088/api/v1/1' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'completed=false' \
--data-urlencode 'title=1'
*/
func updateOne(c *gin.Context) {
	var originalModel originalModel
	id := c.Param("id")
	gormDB.First(&originalModel, id)
	if originalModel.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "item not found",
		})
		return
	}

	title := c.PostForm("title")
	completed := c.PostForm("completed") // 从PostForm拿到字符串"true"
	completedBool,_ := strconv.ParseBool(completed) // 将字符串转为bool
	result := 0 // 将bool转成int
	if completedBool {
		result = 1
	}

	res := gormDB.Model(&originalModel).Updates(map[string]interface{}{
		"title": title,
		"result": result,
	})
	if res.RowsAffected == 0 {
		// put失败返回http 400，可以在响应正文中解释失败原因
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"message": "item update error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusNoContent, // 更新成功，不返回数据:204
		"message": "item update successfully",
	})
}

/*
curl --location --request DELETE 'http://127.0.0.1:8088/api/v1/1'
*/
func deleteOne(c *gin.Context) {
	var originalModel originalModel
	id := c.PostForm("id")
	gormDB.First(&originalModel, id)

	if originalModel.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "item not found",
		})
		return
	}

	gormDB.Delete(&originalModel) // 会给deleted_at字段添加datetime

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": "item delete successfully",
	})
}

// 判断是否存在error
// TODO：添加exitFlag参数，如果提供exitFlag=true表示中断程序执行，抛出fatal error，否则只是抛出error。可以用log.Fatal()和log.Error()区分
func ExistError(err error) bool {
	if err != nil {
		msg := err.Error()
		_, file, line, ok := runtime.Caller(1)
		if ok {
			msg = fmt.Sprintf("file:%s, line:%d, error:%s", file, line, err.Error())
		}
		fmt.Println(msg)
		return true
	}
	return false
}