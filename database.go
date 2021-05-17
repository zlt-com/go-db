package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/zlt-com/go-common"
)

func Open(dbType []string) {

	if common.Contains("mysql", dbType) {
		// 初始化Mysql
		initMysql()
	}
	if common.Contains("redis", dbType) {
		//初始化redis
		initRedis()
	}

}
