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
		fmt.Println("初始化Mysql")
		// 初始化Mysql
		if db, err := initMysql(); err != nil {
			fmt.Println(err)
		} else {
			defaultDB = db
		}

	}

	if common.Contains("redis", dbType) {
		fmt.Println("初始化redis")
		//初始化redis
		initRedis()
	}
}

type Database struct {
	order           interface{}
	offset          interface{}
	limit           interface{}
	model           interface{}
	where           interface{}
	whereConditions []Condition
}

type Condition struct {
	Key   string
	Op    string
	Value interface{}
}

func (ct *Condition) ToString(tablename string) string {
	// switch v := ct.Value.(type) {
	// case string:
	// 	return tablename + "." + ct.Key + ct.Op + "'" + v + "'"
	// case int:
	// 	return tablename + "." + ct.Key + ct.Op + strconv.Itoa(v)
	// default:
	// 	return ""
	// }
	return ct.Key + ct.Op + "?"
}

func (m *Database) Model(model interface{}) *Database {
	c := m.clone()
	c.model = model
	return c
}

func (m *Database) Order(order interface{}) *Database {
	c := m.clone()
	c.order = order
	return c
}

func (m *Database) Offset(offset interface{}) *Database {
	c := m.clone()
	c.offset = offset
	return c
}

func (m *Database) Limit(limit interface{}) *Database {
	c := m.clone()
	c.limit = limit
	return c
}

func (m *Database) Where(where interface{}) *Database {
	c := m.clone()
	switch value := where.(type) {
	case map[string]interface{}:
		for key, val := range value {
			c.whereConditions = append(c.whereConditions, Condition{Key: key, Op: "=", Value: val})
		}
	case []Condition:
		c.whereConditions = value
	}
	return c
}

func (m *Database) clone() *Database {
	return &Database{
		order:           m.order,
		offset:          m.offset,
		limit:           m.limit,
		model:           m.model,
		where:           m.where,
		whereConditions: m.whereConditions,
	}
}
