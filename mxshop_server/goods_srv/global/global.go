package global

import (
	"gorm.io/gorm"
	"mxshop_server/goods_srv/config"
)

var (
	DB *gorm.DB
	ServerConfig config.ServiceConfig
	NacosConfig config.NacosConfig
)

//func init() {
//	newLogger := logger.New(
//		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
//		logger.Config{
//			SlowThreshold:             time.Second,   // 慢 SQL 阈值
//			LogLevel:                  logger.Silent, // 日志级别
//			IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
//			Colorful:                  false,         // 禁用彩色打印
//		},
//	)
//
//	dsn := "root:root@tcp(127.0.0.1:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
//	var err error
//	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//		NamingStrategy: schema.NamingStrategy{
//			SingularTable: true, //去掉结构体名称加 "s"
//		},
//		Logger: newLogger,
//	})
//	if err != nil {
//		panic(err)
//	}
//}
