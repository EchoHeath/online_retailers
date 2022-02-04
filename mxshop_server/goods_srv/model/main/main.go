package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mxshop_server/goods_srv/model"
	"os"
	"time"
)

func main() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // 日志级别
			IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,         // 禁用彩色打印
		},
	)

	dsn := "root:root@tcp(127.0.0.1:3306)/mxshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //去掉结构体名称加 "s"
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	_ = db.AutoMigrate(&model.Category{},
		&model.Brands{},
		&model.GoodsCategoryBrand{},
		&model.Banner{},
		&model.Goods{})

	//var users []model.User
	//res := db.Find(&users)
	//if res.Error != nil {
	//	fmt.Printf("mysql for find err: %v", res.Error)
	//	return
	//}
	//rsp := &proto.UserListResponse{}
	//rsp.Total = int32(res.RowsAffected)
	//db.Scopes(Paginate(1, 10)).Find(&users)
	//for k, user := range users {
	//	fmt.Println(k, user)
	//	userInfoResp := ModelToResponse(user)
	//	rsp.Data = append(rsp.Data, &userInfoResp)
	//
	//}
	//fmt.Println(rsp)

	//option := &password.Options{16, 100, 32, sha512.New}
	//salt, pwd := password.Encode("admin123", option)
	//newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, pwd)
	//fmt.Println(newPassword)
	//
	//for i := 0; i < 10; i++ {
	//	user := model.User{
	//		NickName: fmt.Sprintf("A%d" , i),
	//		Mobile: fmt.Sprintf("1388888889%d" , i),
	//		Password: newPassword,
	//	}
	//	db.Save(&user)
	//}

}

//Paginate 分页逻辑
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
