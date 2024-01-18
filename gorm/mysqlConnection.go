package main

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"runtime"
	"time"
	"os"
	"log"
	"gorm.io/gorm/logger"
	"database/sql/driver"
)

// 声明全局变量
var sqlDBCentral *sql.DB
var gormDBCentral *gorm.DB
var poolCentral *sql.DB
var err error

// 自动初始化
func init() {
	//host := "host.docker.internal" // 假设在容器中运行，访问其他ip段的容器映射到宿主机的服务
	user := "root"
	pass := "123456"
	host := "host.docker.internal"
	port := 3307
	dbname := "central"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local", user, pass, host, port, dbname)
	/*
	参数说明：
	想要正确的处理 time.Time ，您需要带上 parseTime 参数。mysql中的date和datetime等时间类型字段将自动转为golang中的time.Time类型，类似的0000-00-00 00:00:00 ，会被转为time.Time的零值，否则转为[]byte/string类型
	使用charset指定编码，要支持完整的 UTF-8 编码，您需要将 charset=utf8 更改为 charset=utf8mb4
	loc时区默认是utc，一般用上海时区(Asia/Shanghai)，或者Local
	其他参数见 https://github.com/go-sql-driver/mysql#parameters
	*/

	// database/sql创建数据库连接
	sqlDBCentral, err = sql.Open("mysql", dsn)
	if !ExistError(err) {
		// 初始化数据库日志
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold 慢查询阈值
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger 忽略未找到的错误
				Colorful:                  true,        // Disable color
			},
		)

		// GORM 允许通过一个现有的数据库连接来初始化 *gorm.DB
		gormDBCentral, err = gorm.Open(mysql.New(mysql.Config{
			Conn: sqlDBCentral,
		}), &gorm.Config{
			Logger: newLogger,
		})

		if !ExistError(err) {
			// TODO:迁移
			poolCentral, err = gormDBCentral.DB()
			if !ExistError(err) {
				poolCentral.SetMaxIdleConns(10) // 设置空闲连接池中连接的最大数量
				poolCentral.SetMaxOpenConns(100) // 设置打开数据库连接的最大数量
				poolCentral.SetConnMaxLifetime(time.Hour) // 设置了连接可复用的最大时间（空闲超时时间）

				err = poolCentral.Ping()
				if !ExistError(err) {
					fmt.Println("mysql connection ok.")
				}
			}
		}
	}

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

/********************************模型定义********************************/

/*
CREATE TABLE `c_users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  `created_at` int(10) unsigned NOT NULL,
  `updated_at` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;
*/

// type tag 不会影响实际提交的数据，那有何用？
// 对于有 CreatedAt 字段的模型，创建记录时，如果该字段值为零值，则将该字段的值设为当前时间
// 你可以通过将 autoCreateTime 标签置为 false 来禁用时间戳追踪
// 对于有 UpdatedAt 字段的模型，更新记录时，将该字段的值设为当前时间。创建记录时，如果该字段值为零值，则将该字段的值设为当前时间。可以通过将 autoUpdateTime 标签置为 false 来禁用时间戳追踪
// 索引标签，包括index、uniqueIndex。注意复合索引也就是index:index_name，字段顺序会影响其性能，可以使用 priority 指定顺序，默认优先级值是 10，如果优先级值相同，则顺序取决于模型结构体字段的顺序
// 支持复合主键primaryKey
// 等等
type User struct {
	Id		 	int64		`gorm:"primaryKey;column:id;"`
	Username 	string 		`gorm:"column:username;type:varchar(255);default:(-)" `
	Password 	string 		`gorm:"column:password;type:varchar(255);default:(-)"`     
	CreatedAt int64		`gorm:"column:created_at;type:int(10);default:(-)"`
	UpdatedAt int64		`gorm:"column:updated_at;type:int(10);default:(-)"`
	CreateTime *DateTime    `gorm:"column:create_time;type:datetime;default:(-)"`	// 自定义类型：比如db需要datetime，要自己实现类型（要你何用！），也就是实现几个方法。
	UpdateTime *DateTime      `gorm:"column:update_time;type:datetime;default:(-)"`
	// Deleted    gorm.DeletedAt `gorm:"column:deleted;type:timestamp;default:(-)"` // 软删除，鸡肋
}

// TableName 自定义表名
func (*User) TableName() string {
	return "c_users"
}

/********************************操作方法********************************/

/************新增************/

// 新增单个
func Create(user *User) (err error){
	r := gormDBCentral.Model(&user).Create(&user)
	// 错误处理，后面的方法就不单独处理了
	if r.Error != nil {
		return r.Error
	}
	// log:INSERT INTO `c_users` (`username`,`password`,`created_at`,`updated_at`,`create_time`,`update_time`,`id`) VALUES ('a','b',1704872847,1704872847,'2024-01-10 15:47:27.748','2024-01-10 15:47:27.748',1)
	// db:INSERT INTO `central`.`c_users`(`id`, `username`, `password`, `create_time`, `update_time`, `created_at`, `updated_at`) VALUES (1, 'a', 'b', '2024-01-10 15:47:28', '2024-01-10 15:47:28', 1704872847, 1704872847);
	// 再次插入错误写入日志，但不会影响程序继续执行： Error 1062 (23000): Duplicate entry '1' for key 'PRIMARY'
	return nil
}

// 保存单个
// Create和Save的区别：Save判断存在则不插入，Create无论什么情况都执行插入
func Save(user *User)  {
	gormDBCentral.Model(&user).Save(&user)
	/*
	已存在时：
	UPDATE `c_users` SET `username`='a',`password`='b',`created_at`=1704877053,`updated_at`=1704877053,`create_time`='2024-01-10 16:57:33.161',`update_time`='2024-01-10 16:57:33.161' WHERE `id` = 1
	不存在时：
	UPDATE `c_users` SET `username`='b',`password`='b',`created_at`=1704877053,`updated_at`=1704877053,`create_time`='2024-01-10 16:57:33.161',`update_time`='2024-01-10 16:57:33.161' WHERE `id` = 2
	INSERT INTO `c_users` (`username`,`password`,`created_at`,`updated_at`,`create_time`,`update_time`,`id`) VALUES ('b','b',1704877053,1704877053,'2024-01-10 16:57:33.161','2024-01-10 16:57:33.161',2) ON DUPLICATE KEY UPDATE `id`=`id`
	注意 INSERT 后面的 ON DUPLICATE KEY UPDATE `id`=`id`
	*/
	return
}

// 创建多个
func CreateBatch(user []*User)  {
	gormDBCentral.Model(&user).Create(&user)
	/*
	一个语句执行多行数据
	INSERT INTO `c_users` (`create_time`,`update_time`,`id`,`username`,`password`,`created_at`,`updated_at`) VALUES ('2024-01-10 17:08:31.234','2024-01-10 17:08:31.234',3,'c','b',1704877711,1704877711),('2024-01-10 17:08:31.234','2024-01-10 17:08:31.234',4,'d','b',1704877711,1704877711)
	*/
	return
}

/************查询************/

func GetFirst() (user *User) {
	gormDBCentral.Model(&user).First(&user) // 同理，还有Last()
	/*
	没数据时：
	SELECT * FROM `c_users` ORDER BY `c_users`.`id` LIMIT 1
	&{0   {0001-01-01 00:00:00 +0000 UTC false} 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC}
	有数据时：
	SELECT * FROM `c_users` ORDER BY `c_users`.`id` LIMIT 1
	&{1 a b 1704872847 1704872847 2024-01-10 15:47:28 2024-01-10 15:47:28}
	*/
	return 
}

func GetById(id int64) (user *User){
	// 方法1
	gormDBCentral.Model(&user).Where("id = ?", id).Find(&user) // 硬编码略愚蠢：SELECT * FROM `c_users` WHERE id = 1
	// 方法2
	gormDBCentral.Model(&user).Find(&user, id) // 兼容ids一个条件也要AND略愚蠢：SELECT * FROM `c_users` WHERE `c_users`.`id` = 1 AND `c_users`.`id` = 1
	return 
}

func GetByIds(ids []int64) (user *User){
	// 方法1
	gormDBCentral.Model(&user).Where("id in ?", ids).Find(&user) // SELECT * FROM `c_users` WHERE id in (1)
	// 方法2
	gormDBCentral.Model(&user).Find(&user, ids) // SELECT * FROM `c_users` WHERE `c_users`.`id` = 1 AND `c_users`.`id` = 1
	return 
}

func GetSomeParam(id int64) (user *User) {
	gormDBCentral.Model(&user).Select("username", "password").Find(&user, id) // 获取部分列：SELECT `username`,`password` FROM `c_users` WHERE `c_users`.`id` = 1
	// log：&{0 a b 0 0 <nil> <nil>}
	return
}

func GetPage(limit int, offset int) (user []*User) {
	gormDBCentral.Model(&user).Limit(limit).Offset(offset).Find(&user)
	/*
	offset或limit=0时就不会编入sql，比如limit=1,offset=0：
	SELECT * FROM `c_users` LIMIT 1
	正常情况下：
	SELECT * FROM `c_users` LIMIT 1 OFFSET 1
	注意：返回切片
	*/
	return
}

func GetByOrder() (user []*User) {
	gormDBCentral.Model(&user).Order("id desc, username").Find(&user) // SELECT * FROM `c_users` ORDER BY id desc, username
	gormDBCentral.Model(&user).Order("id desc").Order("username").Find(&user) // 拼接竟然不带空格：SELECT * FROM `c_users` ORDER BY id desc,username
	return
}

/************扫描************/

func ScanRow(id int64) (row *sql.Row) {
	row = gormDBCentral.Table("c_users").Where("id = ?", id).Select("username", "password").Row()
	return
}

func ScanRows(username string) (rows *sql.Rows) {
	rows, _ = gormDBCentral.Model(&User{}).Where("username = ?", username).Select("username, password, id").Rows()
	return
}

/************更新************/

// 更新单个字段
func UpdateUsername(id int64, username string)  {
	gormDBCentral.Model(&User{}).Where("id = ?",id).Update("username",username) // UPDATE `c_users` SET `username`='x',`updated_at`=1704878474 WHERE id = 1
	return
}

// 全量/多列更新（根据结构体）
func UpdateByUser(user *User)  {
	gormDBCentral.Model(&User{}).Where("id = ?",user.Id).Updates(&user)
	/*
	UPDATE `c_users` SET `id`=1,`username`='y',`password`='y',`created_at`=1704878996,`updated_at`=1704878996,`create_time`='2024-01-10 17:29:56.063',`update_time`='2024-01-10 17:29:56.063' WHERE id = 1
	*/
	return
}

/************删除************/

// 简单删除（根据user里的id进行删除）
func DeleteByUser(user *User)  {
	gormDBCentral.Model(&User{}).Delete(&user)
	/*
	鸡肋：
	结构体未加gorm.DeletedAt标记的字段，直接删除，加了将更新deleted字段，即实现软删除
	当前没加，所以是
	DELETE FROM `c_users` WHERE `c_users`.`id` = 4
	*/
	return
}

// 根据id进行删除
func DeleteById(id int64)  {
	gormDBCentral.Model(&User{}).Delete(&User{}, id) // DELETE FROM `c_users` WHERE `c_users`.`id` = 3
	return
}

/********************************事务********************************/

// 匿名事务
func AnonymousTransaction() error {
	err := gormDBCentral.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，使用 'tx' 而不是 'db'）
		var now_datetime = new(DateTime)
		_ = now_datetime.Scan(time.Now())
		now_int := time.Now().Unix()
		// 一定成功的语句
		// log: INSERT INTO `c_users` (`username`,`password`,`created_at`,`updated_at`,`create_time`,`update_time`) VALUES ('m','m',1704880196,1704880196,'2024-01-10 17:49:56.419','2024-01-10 17:49:56.419')
		if err := tx.Create(&User{
			Username: "m",
			Password: "m",
			CreateTime: now_datetime,
			UpdateTime: now_datetime,
			CreatedAt: now_int,
			UpdatedAt: now_int,
		}).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		
		// 一定失败的语句（Error 1062 (23000): Duplicate entry '1' for key 'PRIMARY'）
		// log: INSERT INTO `c_users` (`username`,`password`,`created_at`,`updated_at`,`create_time`,`update_time`,`id`) VALUES ('m','m',1704880455,1704880455,'2024-01-10 17:54:15.318','2024-01-10 17:54:15.318',1)
		if err := tx.Create(&User{
			Id: int64(1),
			Username: "m",
			Password: "m",
			CreateTime: now_datetime,
			UpdateTime: now_datetime,
			CreatedAt: now_int,
			UpdatedAt: now_int,
		}).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}

		// 删除失败不代表发生了错误，只是没匹配而已，不会导致回滚
		// log：DELETE FROM `c_users` WHERE `c_users`.`id` = 100
		if err := tx.Model(&User{}).Delete(&User{Username: "not_exist", Id: int64(100)}).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}

		// 返回 nil 提交事务
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// 手动处理事务(自定义)
func CustomTransaction() error {
	tx := gormDBCentral.Begin()

	// 需要在defer中recover判断整个事务最终提交时是否发生了错误，有错误也要回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
	if err := tx.Create(&User{Username: "lomtom"}).Error; err != nil {
		// Error 1364 (HY000): Field 'password' doesn't have a default value
		// 回滚事务
		tx.Rollback()
		return err 
	}
	if err := tx.Delete(&User{}, 28).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果还没return，说明没有碰到任何错误，提交整个事务。注意：提交整个事务时如果返回错误，需要在defer回滚
	// 提交事务
	return tx.Commit().Error
}

/********************************原生SQL********************************/

func GetOneByRaw(id int64) (user *User) {
	gormDBCentral.Raw("SELECT id, username, password FROM c_users WHERE id = ?", id).Scan(&user)
	// log: {1 y y 0 0 <nil> <nil>}
	// 注意：后面四个是零值，0不是db返回的值
	return
}

func GetMultiByRaw(ids []int64) (users []*User) {
	gormDBCentral.Raw("SELECT * FROM c_users WHERE id IN ?", ids).Scan(&users)
	return
}

func ExecExpressionByRaw(id int64) {
	gormDBCentral.Exec("UPDATE c_users SET created_at = ? WHERE id = ?", gorm.Expr("created_at * ? + ?", 1, 3600), id)
	// UPDATE c_users SET created_at = created_at * 1 + 3600 WHERE id = 1
	return
}

type NamedArgument struct {
    Id int64
    Username string
}

func GetByRawNamedParam(id int64, username string) (user *User) {
	gormDBCentral.Raw("SELECT * FROM c_users WHERE id = @id and username = @username", sql.Named("id", id), sql.Named("username", username)).Scan(&user)
	gormDBCentral.Raw("SELECT * FROM c_users WHERE id = @id and username = @username", map[string]interface{}{"id": id, "username": username}).Scan(&user)
	gormDBCentral.Raw("SELECT * FROM c_users WHERE id = @Id and username = @Username", NamedArgument{Id: id, Username: username}).Scan(&user)
	/*
	拼接：SELECT * FROM c_users WHERE id = 1 and username = 'y'
	方法1：只需要少量具名参数时sql.Named
	方法2：基本类型可以使用map，值用interface{}传入
	方法3：特殊或结构化类型可以使用自定义结构体传入
	*/
	return
}

func DryRun(id int64, user *User) *gorm.Statement{
	stmt := gormDBCentral.Session(&gorm.Session{DryRun: true}).First(&user, id).Statement
	// stmt.SQL.String() //=> SELECT * FROM `c_users` WHERE `id` = $1 ORDER BY `id`
	// stmt.Vars         //=> []interface{}{1}
	return stmt
}

func ToSQL(id int64, limit int) (string){
	return gormDBCentral.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&User{}).Where("id = ?", id).Limit(limit).Order("id desc").Find(&[]User{})
		// SELECT * FROM `c_users` WHERE id = 1 ORDER BY id desc LIMIT 10
	})
}

/********************************主函数********************************/

func main() {
	var now_datetime = new(DateTime)
	_ = now_datetime.Scan(time.Now())

	now_int := time.Now().Unix()
	
	newUser := User{
		Id: 1,
		Username: "a",
		Password: "b",
		CreateTime: now_datetime,
		UpdateTime: now_datetime,
		CreatedAt: now_int,
		UpdatedAt: now_int,
	}
	err = Create(&newUser)
	if err != nil {
		fmt.Println("Create", err)
	}

	u := GetFirst()
	fmt.Println("GetFirst", u)

	id := int64(1)
	u = GetById(id)
	fmt.Println("GetById", u)

	ids := []int64{int64(1)}
	u = GetByIds(ids)
	fmt.Println("GetByIds", u)

	u = GetSomeParam(id)
	fmt.Println("GetSomeParam", u)

	limit := 1
	offset := 0
	var us []*User
	us = GetPage(limit, offset)
	fmt.Println("GetPage", us)
	for _, v := range us {
		fmt.Println("GetPage range", v)
		fmt.Println("GetPage range Username", v.Username)
	}

	us = GetByOrder()
	for _, v := range us {
		fmt.Println("GetByOrder range", v)
	}

	newUser.UpdatedAt = time.Now().Unix()
	Save(&newUser)

	newUser2 := newUser
	newUser2.Id = 2
	newUser2.Username = "b"
	Save(&newUser2)

	newUser3 := newUser2
	newUser3.Id = 3
	newUser3.Username = "c"
	newUser4 := newUser2
	newUser4.Id = 4
	newUser4.Username = "d"
	CreateBatch([]*User{&newUser3, &newUser4})

	var username string
	var password string
	row := ScanRow(id)
	row.Scan(&username, &password)
	fmt.Println("ScanRow to variables")
	fmt.Println(username)
	fmt.Println(password)

	username = "m"
	rows := ScanRows(username)
	defer rows.Close()
	fmt.Println("ScanRows to variables")
	for rows.Next() {
		rows.Scan(&username, &password, &id)
		fmt.Println(username)
		fmt.Println(password)
		fmt.Println(id)
	}

	var u *User
	rows = ScanRows(username)
	defer rows.Close()
	for rows.Next() {
		// ScanRows 将一行扫描至 user
		gormDBCentral.ScanRows(rows, &u)
		// 业务逻辑...
	}
	fmt.Println("ScanRows to model", u) // &{10 m m 0 0 <nil> <nil>} 注意零值

	UpdateUsername(id, "x")

	newUser.Username = "y"
	newUser.Password = "y"
	UpdateByUser(&newUser)

	DeleteByUser(&newUser4)

	id3 := int64(3)
	DeleteById(id3)

	AnonymousTransaction()

	CustomTransaction()

	var res = User{}
	u = GetOneByRaw(id)
	fmt.Println("GetOneByRaw", u)
	fmt.Println("GetOneByRaw CreatedAt", u.CreatedAt) // 注意：这是零值，不是db返回的值

	ids = []int64{1, 2, 3}
	us = GetMultiByRaw(ids)
	fmt.Println(us)
	for _, v := range us {
		fmt.Println("GetMultiByRaw range", v)
	}

	ExecExpressionByRaw(id)

	GetByRawNamedParam(id, "y")

	stmt := DryRun(id, u)
	fmt.Println("DryRun statement", stmt)
	fmt.Println("DryRun statement sql string", stmt.SQL.String()) //=> SELECT * FROM `c_users` WHERE `id` = $1 ORDER BY `id`
	fmt.Println("DryRun statement vars", stmt.Vars) //=> []interface{}{1}

	limit = 10
	sql := ToSQL(id , limit)
	fmt.Println(sql)


}

/********************************实现自定义类型DateTime********************************/

type DateTime time.Time
 
const (
    datetimeFormart = "2006-01-02 15:04:05"
    zone        = "Asia/Shanghai"
)
 
// UnmarshalJSON implements json unmarshal interface.
func (t *DateTime) UnmarshalJSON(data []byte) (err error) {
    now, err := time.ParseInLocation(`"`+datetimeFormart+`"`, string(data), time.Local)
    *t = DateTime(now)
    return
}
 
// MarshalJSON implements json marshal interface.
func (t DateTime) MarshalJSON() ([]byte, error) {
    b := make([]byte, 0, len(datetimeFormart)+2)
    b = append(b, '"')
    b = time.Time(t).AppendFormat(b, datetimeFormart)
    b = append(b, '"')
    return b, nil
}
 
func (t DateTime) String() string {
    return time.Time(t).Format(datetimeFormart)
}
 
func (t DateTime) local() time.Time {
    loc, _ := time.LoadLocation(zone)
    return time.Time(t).In(loc)
}
 
// Value ...
func (t DateTime) Value() (driver.Value, error) {
    var zeroTime time.Time
    var ti = time.Time(t)
    if ti.UnixNano() == zeroTime.UnixNano() {
        return nil, nil
    }
    return ti, nil
}
 
// Scan valueof time.Time 注意是指针类型 method
func (t *DateTime) Scan(v interface{}) error {
    value, ok := v.(time.Time)
    if ok {
        *t = DateTime(value)
        return nil
    }
    return fmt.Errorf("can not convert %v to DateTime", v)
}

/****************************************************************/

/*
参考：
https://juejin.cn/post/7034706176050724878
https://gorm.io/zh_CN/docs/create.html
*/