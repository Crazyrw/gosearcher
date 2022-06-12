package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	// "gorm.io/gorm/logger"
	// "log"
	// "os"
	"time"
)

var MysqlDB *gorm.DB
var err error

//root:122513gzhGZH!!@tcp(decs.pcl.ac.cn:1762)/search_engine?
//root:root@tcp(101.200.128.148:3306)/documents?
func ConnectMySql() {
	dsn := "root:122513gzhGZH!!@tcp(decs.pcl.ac.cn:1762)/search_engine?charset=utf8mb4&parseTime=True&loc=Local"
	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
	// 	logger.Config{
	// 		SlowThreshold:             time.Second, // 慢 SQL 阈值
	// 		LogLevel:                  logger.Info, // 日志级别
	// 		IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
	// 		Colorful:                  true,        // 彩色打印
	// 	},
	// )
	//全局模式
	// MysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
	// 	Logger: newLogger,
	// })
	MysqlDB, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("数据库连接失败:" + err.Error())
	}
	fmt.Println("数据库连接成功")
	//获取通用数据库对象 sql.DB
	MysqlDBPool, err := MysqlDB.DB()
	if err != nil {
		panic("获取数据库对象sql.DB失败")
	}
	//设置连接池中空闲连接的最大数量
	MysqlDBPool.SetMaxIdleConns(5)
	//设置打开数据库连接的最大数量
	MysqlDBPool.SetMaxOpenConns(5)
	//设置连接可复用的最大时间
	MysqlDBPool.SetConnMaxLifetime(time.Hour)
	//设置超时时间 与数据库保持一致 interactive_timeout：交互式连接的超时时间（mysql）
	MysqlDBPool.SetConnMaxLifetime(time.Duration(7*3600) * time.Second)

	//自动迁移 建表
	//err = DB.AutoMigrate(&Docs{})
	//if err != nil {
	//	fmt.Println("自动迁移失败....", err.Error())
	//}
}
