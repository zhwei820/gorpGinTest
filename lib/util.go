package lib

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/timehop/jimmy/redis"
)

// SetCache : 设置redis缓存
func SetCache(c *gin.Context, key string, val string, ttl int) {
	pool := c.Keys["Pool"].(redis.Pool)
	err := pool.SetEx(key, val, ttl)
	if err != nil {
		log.Println("SetCache err", err)
	}
}

// GetGache : 获取redis缓存
func GetGache(c *gin.Context, key string) (string, error) {
	pool := c.Keys["Pool"].(redis.Pool)
	res, err := pool.Get(key)
	return res, err
}
