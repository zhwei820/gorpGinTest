package lib

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

// SetCache : 设置redis缓存
func SetCache(c *gin.Context, key string, val string, ttl int) {
	pool := c.Keys["Pool"].(*redis.Pool)
	con := pool.Get()
	defer con.Close()

	con.Do("SET", key, val, ttl)
}

// GetGache : 获取redis缓存
func GetGache(c *gin.Context, key string) (string, error) {
	pool := c.Keys["Pool"].(*redis.Pool)
	con := pool.Get()
	defer con.Close()

	user, err := redis.String(con.Do("GET", key))
	return user, err
}
