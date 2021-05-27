package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/zlt-com/go-config"
	"github.com/zlt-com/go-logger"
)

// Mysql 全局数据库实例
var (
	mysql *gorm.DB
)

func initMysql() (mysql *gorm.DB, err error) {
	//初始化Mysql
	if mysql, err = gorm.Open(config.Config.DBType, config.Config.DBSource); err != nil {
		// fmt.Println((err))
		logger.Error(err)
		return
	}
	mysql.LogMode(config.Config.LogMode)
	mysql.SetLogger(&logger.MyLogger{})
	mysql.DB().SetMaxIdleConns(1)
	mysql.DB().SetMaxOpenConns(5)
	return
}

func Mysql() *gorm.DB {
	return mysql
}
