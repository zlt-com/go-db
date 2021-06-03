package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/zlt-com/go-config"
	"github.com/zlt-com/go-logger"
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
	mysql.DB().SetMaxIdleConns(5)
	mysql.DB().SetMaxOpenConns(16)
	return
}

func Mysql() *gorm.DB {
	return defaultDB
}
