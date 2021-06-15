package database

import (
	"time"

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
	mysql.LogMode(false)
	mysql.SetLogger(&logger.MyLogger{})
	mysql.DB().SetMaxIdleConns(4)
	mysql.DB().SetMaxOpenConns(16)
	mysql.DB().SetConnMaxLifetime(3 * time.Second)
	return
}

func Mysql() (mysql *gorm.DB, err error) {
	if mysql, err = gorm.Open(config.Config.DBType, config.Config.DBSource); err != nil {
		// fmt.Println((err))
		logger.Error(err)
		return
	}
	return
}

func GetIdleConn() (int, int, int, int, int) {
	return defaultDB.DB().Stats().Idle, defaultDB.DB().Stats().InUse, redisClient.Stats().IdleCount, redisClient.Stats().ActiveCount, redisClient.MaxActive
}
