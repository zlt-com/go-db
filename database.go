package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/zlt-com/go-common"
)

var defaultDB *gorm.DB

func Open(dbType []string) {
	if common.Contains("mysql", dbType) {
		// 初始化Mysql
		if db, err := initMysql(); err != nil {
			fmt.Println(err)
		} else {
			defaultDB = db
		}

	}
	if common.Contains("redis", dbType) {
		//初始化redis
		initRedis()
	}
}

type Database struct {
	order  interface{}
	offset interface{}
	limit  interface{}
	model  interface{}
	where  map[string]interface{}
}

func (m *Database) Model(model interface{}) *Database {
	m.model = model
	return m
}

func (m *Database) Order(order interface{}) *Database {
	m.order = order
	return m
}

func (m *Database) Offset(offset interface{}) *Database {
	m.offset = offset
	return m
}

func (m *Database) Limit(limit interface{}) *Database {
	m.limit = limit
	return m
}

func (m *Database) Where(where map[string]interface{}) *Database {
	m.where = where
	return m
}
