package database

import (
	"bytes"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/zlt-com/go-config"
	"github.com/zlt-com/go-logger"
)

// RedisDB redis 基础类
type RedisDB struct {
	DBNum int
}

var Redisdb *gorm.DB

// 初始化redis
func initRedis() {

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

var (
	redisClient *redis.Pool
)

// connect 连接到Redis数据库
func (redisDB *RedisDB) connect() redis.Conn {
	c := redisClient.Get()
	c.Do("SELECT", redisDB.DBNum)
	return c
}

// 查看redis信息
func (redisDB *RedisDB) Info() (reply interface{}, err error) {
	c := redisDB.connect()
	defer c.Close()
	reply, err = c.Do("info")
	return
}

// Set 存储String型数据
func (redisDB *RedisDB) Set(value ...interface{}) (reply interface{}, err error) {
	c := redisDB.connect()
	defer c.Close()
	reply, err = c.Do("set", value...)
	if err != nil {
		logger.Error("redis.CacheDB set error", err, value)
	}
	return
}

// Psetex 存储String型数据
func (redisDB *RedisDB) Psetex(value ...interface{}) (reply interface{}, err error) {
	c := redisDB.connect()
	defer c.Close()
	reply, err = c.Do("psetex", value...)
	if err != nil {
		logger.Error("redis.CacheDB Psetex error", err, value)
	}
	return
}

// Get 获取String型数据
func (redisDB *RedisDB) Get(key string) (reply interface{}, err error) {
	c := redisDB.connect()
	defer c.Close()
	reply, err = redis.String(c.Do("get", key))
	if err != nil {
		logger.Error("redis.CacheDB Get error", err, key)
	}
	return
}

// Del 删除key
func (redisDB *RedisDB) Del(key string) (reply interface{}, err error) {
	c := redisDB.connect()
	defer c.Close()
	reply, err = c.Do("del", key)
	if err != nil {
		logger.Error("redis.CacheDB Del error", err, key)
		return
	}
	return
}

// String 转换redis数据为字符串型，默认是uint8数组
func (redisDB *RedisDB) String(arg interface{}) string {
	value, err := redis.String(arg, nil)
	if err != nil {
		logger.Error("redis.CacheDB Strings failed:", err, arg)
		return "error string"
	}
	return value
}

// DelKey 删除key
func (redisDB *RedisDB) DelKey(key string) {
	c := redisDB.connect()
	defer c.Close()
	_, err := c.Do("del", key)
	if err != nil {
		logger.Error("redis.CacheDB DelKey error", key, err)
	}
}

// Hget 获取Hash型数据
func (redisDB *RedisDB) Hget(key string, value interface{}) (reply interface{}, err error) {
	c := redisDB.connect()
	defer c.Close()
	reply, err = c.Do("hget", key, value)
	if err != nil {
		logger.Error("redis.CacheDB hget error ", err, key, value)
		return
	}
	return
}

// Hset 存储Hash型数据
func (redisDB *RedisDB) Hset(key string, field, value interface{}) error {
	c := redisDB.connect()
	defer c.Close()
	_, err := c.Do("hset", key, field, value)
	if err != nil {
		logger.Error("redis.CacheDB hset error ", err, key, field, value)
	}
	return err
}

// Hdel 删除Hash型数据
func (redisDB *RedisDB) Hdel(key string, field interface{}) error {
	c := redisDB.connect()
	defer c.Close()
	_, err := c.Do("hdel", key, field)
	if err != nil {
		logger.Error("redis.CacheDB hdel error", err, key, field)
	}
	return err
}

// Hexists 判断field是否存在
func (redisDB *RedisDB) Hexists(key, field interface{}) (has bool, err error) {
	c := redisDB.connect()
	defer c.Close()
	has, err = redis.Bool(c.Do("hexists", key, field))
	if err != nil {
		logger.Error("redis.CacheDB hexists error", err, key, field)
	}

	return
}

// Hsetappend 存储Hash型数据，把字符串用逗号分割
func (redisDB *RedisDB) Hsetappend(key string, field string, value string) {
	c := redisDB.connect()
	defer c.Close()
	exists, _ := redis.Bool(c.Do("hexists", key, field))
	if exists {
		var buf bytes.Buffer
		result, _ := c.Do("hget", key, field)
		vv, ok := result.([]byte)
		if ok {
			buf.Write(vv)
		}

		_, err := c.Do("hset", key, field, value+","+buf.String()) //写
		if err != nil {
			logger.Error("redis.CacheDB Hsetappend failed:", err)
		}
	} else {
		_, err := c.Do("hset", key, field, value) //写
		if err != nil {
			logger.Error("redis.CacheDB set failed:", err)
		}
	}

}

// Hmget 批量获取Hash型数据
func (redisDB *RedisDB) Hmget(field ...interface{}) ([]interface{}, error) {
	c := redisDB.connect()
	defer c.Close()
	result, err := redis.Values(c.Do("hmget", field...))
	if err != nil {
		logger.Error("redis.CacheDB Hmget failed:", err, field)
		return nil, err
	}
	return result, err
}

// Hgetall 批量获取Hash型数据
func (redisDB *RedisDB) Hgetall(key string) []interface{} {
	c := redisDB.connect()
	defer c.Close()
	result, err := redis.Values(c.Do("HVALS", key))
	if err != nil {
		logger.Error("redis.CacheDB Hgetall failed:", err, key)
		return nil
	}
	return result
}

// Hmset 批量存储Hash型数据
func (redisDB *RedisDB) Hmset(field ...interface{}) error {
	c := redisDB.connect()
	defer c.Close()
	_, err := c.Do("hmset", field...)
	if err != nil {
		logger.Error("redis.CacheDB Hmset failed:", err, field)
	}
	return err
}

// Hkeys 获取Hash型数据所有field
func (redisDB *RedisDB) Hkeys(key string) []interface{} {
	c := redisDB.connect()
	defer c.Close()
	result, err := redis.Values(c.Do("hkeys", key))
	if err != nil {
		logger.Error("redis.CacheDB Hkeys error", err, "key:", key)
		return nil
	}
	return result
}

// Hvals 获取Hash型数据所有值
func (redisDB *RedisDB) Hvals(key string) []interface{} {
	c := redisDB.connect()
	defer c.Close()
	result, err := redis.Values(c.Do("hvals", key))
	if err != nil {
		logger.Error("redis.CacheDB Hvals error", err, "key:", key)
		return nil
	}
	return result
}

// Sadd 存储Set型数据
func (redisDB *RedisDB) Sadd(key string, value interface{}) (err error) {
	c := redisDB.connect()
	defer c.Close()
	_, err = c.Do("sadd", key, value)
	if err != nil {
		logger.Error("redis.CacheDB Sadd error", err, key, value)
	}
	return
}

// Sadd 存储Set型数据
func (redisDB *RedisDB) Smembers(key, value interface{}) (reply interface{}, err error) {
	c := redisDB.connect()
	defer c.Close()
	reply, err = c.Do("smembers", key, value)
	if err != nil {
		logger.Error("redis.CacheDB Sadd error", err, key, value)
	}
	return
}

// Lpush 存储List型数据
func (redisDB *RedisDB) Lpush(args ...interface{}) {
	c := redisDB.connect()
	defer c.Close()
	_, err := c.Do("lpush", args...)
	if err != nil {
		logger.Error("redis.CacheDB Lpush error", err, args)
	}
}

// Lrange 获取List型数据
func (redisDB *RedisDB) Lrange(args ...interface{}) []interface{} {
	c := redisDB.connect()
	defer c.Close()
	v, err := redis.Values(c.Do("lrange", args...))
	if err != nil {
		logger.Error("redis.CacheDB Lpush error", err, args)
	}
	return v
}

// Llen 获取List型数据数量
func (redisDB *RedisDB) Llen(key interface{}) int {
	c := redisDB.connect()
	defer c.Close()
	v, err := redis.Int(c.Do("llen", key))
	if err != nil {
		logger.Error("redis.CacheDB Llen error", err, key)
	}
	return v
}

// Zadd 存储sortset型数据
func (redisDB *RedisDB) Zadd(key string, sort int, value interface{}) (err error) {
	c := redisDB.connect()
	defer c.Close()
	_, err = c.Do("zadd", key, sort, value) //写
	if err != nil {
		logger.Error("redis.CacheDB set failed:", err, key, sort, value)
	}
	return
}

// Zcard 获取sortset型数据数量
func (redisDB *RedisDB) Zcard(key string) (count int, err error) {
	c := redisDB.connect()
	defer c.Close()
	count, err = redis.Int(c.Do("zcard", key))
	if err != nil {
		logger.Error("redis.CacheDB Zcard failed:", err, key)
	}
	return
}

// Zrange 获取sortset型数据
func (redisDB *RedisDB) Zrange(key string, start, end interface{}) ([]interface{}, error) {
	c := redisDB.connect()
	defer c.Close()
	values, err := redis.Values(c.Do("zrange", key, start, end, " WITHSCORES"))
	if err != nil {
		logger.Error("redis.CacheDB zrange failed:", err, key, start, end)
	}
	return values, err
}

// Zrevrange 获取sortset型数据，从大到小
func (redisDB *RedisDB) Zrevrange(key string, start, end interface{}) ([]interface{}, error) {
	c := redisDB.connect()
	defer c.Close()
	values, err := redis.Values(c.Do("zrevrange", key, start, end))
	if err != nil {
		logger.Error("redis.CacheDB Zrevrange failed:", err, key, start, end)
	}
	return values, err
}

// ZrevrangeStrings 获取sortset型数据，从大到小，转化成字符串
func (redisDB *RedisDB) ZrevrangeStrings(args ...interface{}) []string {
	c := redisDB.connect()
	defer c.Close()
	values, err := redis.Strings(c.Do("zrevrange", args...))
	if err != nil {
		logger.Error("redis.CacheDB ZrevrangeStrings failed:", err, args)
	}
	return values
}

// Zrank 返回指定key的索引未知，空为没有找到
func (redisDB *RedisDB) Zrank(args ...interface{}) interface{} {
	c := redisDB.connect()
	defer c.Close()
	values, err := c.Do("zrank", args...)
	if err != nil {
		logger.Error("redis.CacheDB Zrank failed:", err, args)
	}
	return values
}

// ZrangeByScore 获取sortset型数据
func (redisDB *RedisDB) ZrangeByScore(args ...interface{}) []interface{} {
	c := redisDB.connect()
	defer c.Close()
	values, err := redis.Values(c.Do("ZRANGEBYSCORE", args...))
	if err != nil {
		logger.Error("redis.CacheDB ZrangeByScore failed:", err, args)
	}
	return values
}

// ZrangeByLex 获取sortset型数据，参数中包含查询字符，类似于like
func (redisDB *RedisDB) ZrangeByLex(args ...interface{}) []interface{} {
	c := redisDB.connect()
	defer c.Close()
	values, err := redis.Values(c.Do("zrangeByLex", args...))
	if err != nil {
		logger.Error("redis.CacheDB ZrangeByLex failed:", err, args)
	}
	return values
}

// Zunionstore 合并多个sortset
func (redisDB *RedisDB) Zunionstore(args ...interface{}) {
	c := redisDB.connect()
	defer c.Close()
	_, err := c.Do("zunionstore", args...)
	if err != nil {
		logger.Error("redis.CacheDB Zunionstore failed:", err, args)
	}
}

// Zrem 删除sortset成员
func (redisDB *RedisDB) Zrem(key string, args ...interface{}) (err error) {
	c := redisDB.connect()
	defer c.Close()
	_, err = c.Do("zrem", args...)
	if err != nil {
		logger.Error("redis.CacheDB Zrem failed:", err, key, args)
	}
	return
}

// Zdelete 删除sortset型数据
func (redisDB *RedisDB) Zdelete(key, member string) error {
	c := redisClient.Get()
	defer c.Close()
	_, err := c.Do("zrem", key, member)
	if err != nil {
		logger.Error("redis.CacheDB Zdelete failed:", err, key, member)
	}
	return err
}

// Keys 获取符合条件的key
func (redisDB *RedisDB) Keys(rex string) []interface{} {
	c := redisDB.connect()
	defer c.Close()
	values, err := redis.Values(c.Do("keys", rex))
	if err != nil {
		logger.Error("redis.CacheDB Keys failed:", err, rex)
	}
	return values
}

// Type 返回 key 所储存的值的类型
func (redisDB *RedisDB) Type(key string) string {
	c := redisDB.connect()
	defer c.Close()
	values, err := redis.String(c.Do("type", key))
	if err != nil {
		logger.Error("redis.CacheDB type failed:", err, key)
	}
	return values
}

// Incr 自增数值
func (redisDB *RedisDB) Incr() int {
	c := redisDB.connect()
	defer c.Close()
	maxid, err := redis.Int(c.Do("incr", "maxid"))
	if err != nil {
		logger.Error("redis.CacheDB Incr failed:", err)
	}
	return maxid
}

// Hlen hash的数量
func (redisDB *RedisDB) Hlen(key string) int {
	c := redisDB.connect()
	defer c.Close()
	lens, err := redis.Int(c.Do("hlen", key))
	if err != nil {
		logger.Error("redis.CacheDB Hlen failed:", err, key)
	}
	return lens
}

// Strings 转换redis数据为字符串型，默认是uint8数组
func (redisDB *RedisDB) Strings(arg interface{}) []string {
	value, err := redis.Strings(arg, nil)
	if err != nil {
		logger.Error("redis.CacheDB Strings failed:", err, arg)
		return nil
	}
	return value
}

// Bool 转换redis数据为Bool型，默认是uint8数组
func (redisDB *RedisDB) Bool(arg interface{}) bool {
	value, err := redis.Bool(arg, nil)
	if err != nil {
		logger.Error("redis.CacheDB Bool failed:", err, arg)
		return false
	}
	return value
}

// Expire 设置 key 的过期时间
func (redisDB *RedisDB) Expire(key string, time int64) {
	c := redisDB.connect()
	defer c.Close()
	_, err := redis.Int(c.Do("Expire", key, time))
	if err != nil {
		logger.Error("redis.CacheDB Expire failed:", err, key)
	}
}

// EXISTS key 是否存在
func (redisDB *RedisDB) Exists(key string) (has bool, err error) {
	c := redisDB.connect()
	defer c.Close()
	has, err = redis.Bool(c.Do("Exists", key))
	if err != nil {
		logger.Error("redis.CacheDB Expire failed:", err, key)
	}
	return
}
