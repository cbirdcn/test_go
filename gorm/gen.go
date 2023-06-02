package main

import 
(
	"gorm.io/gen"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

// Dynamic SQL
type Querier interface {
  // SELECT * FROM @@table WHERE user_id = @user_id
  FilterWithUser(user_id int) ([]gen.T, error)
}

func main() {
  g := gen.NewGenerator(gen.Config{
    OutPath: "./gen_query",
    Mode: gen.WithoutContext|gen.WithDefaultQuery|gen.WithQueryInterface, // generate mode
  })

  gormdb, _ := gorm.Open(mysql.Open("root:root@(127.0.0.1:3307)/c_login_device_history?charset=utf8mb4&parseTime=True&loc=Local"))
  g.UseDB(gormdb) // reuse your gorm db

  // Generate basic type-safe DAO API for struct `model.User` following conventions
  g.ApplyBasic(model.LoginDeviceHistory{})

  // Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
  g.ApplyInterface(func(Querier){}, model.LoginDeviceHistory{})

  // Generate the code
  g.Execute()
}
