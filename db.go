package db

import (
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/zlt-com/go-config"
	"github.com/zlt-com/go-logger"
)

// Mysql 全局数据库实例
var (
	Mysql          *gorm.DB
	OauthCode      = RedisDB{0}
	TokenDB        = RedisDB{1}
	RefreshTokenDB = RedisDB{2}
)

func init() {
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

	//初始化redis
	redisClient = &redis.Pool{
		MaxIdle:     2,
		MaxActive:   8,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				config.Config.RedisType,
				config.Config.RedisHost,
				redis.DialPassword(config.Config.RedisPass),
			)
			if err != nil {
				logger.Error(err.Error())
				return nil, err
			}
			// 选择CacheDB
			// c.Do("SELECT", 1)
			return c, nil
		},
	}
}
