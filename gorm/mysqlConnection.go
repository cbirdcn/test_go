package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"runtime"
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
			// TODO:迁移
			db, err = gormDB.DB()

			if !ExistError(err) {
				db.SetMaxIdleConns(10) // 设置空闲连接池中连接的最大数量
				db.SetMaxOpenConns(100) // 设置打开数据库连接的最大数量
				db.SetConnMaxLifetime(time.Hour) // 设置了连接可复用的最大时间

				fmt.Println("mysql connection ok.")
			}
		}
	}

}


func main() {
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