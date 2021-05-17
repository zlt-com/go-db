package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/zlt-com/go-config"
	"github.com/zlt-com/go-logger"
)

// Mysql 全局数据库实例
var (
	Mysql *gorm.DB
)

func initMysql() {
	//初始化Mysql
	var err error
	if Mysql, err = gorm.Open(config.Config.DBType, config.Config.DBSource); err != nil {
		// fmt.Println((err))
		logger.Error(err)
	}
	Mysql.LogMode(config.Config.LogMode)
	Mysql.SetLogger(&logger.MyLogger{})
	Mysql.DB().SetMaxIdleConns(1)
	Mysql.DB().SetMaxOpenConns(5)

}
