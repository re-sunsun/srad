package cache

import (
	"github.com/gomodule/redigo/redis"
)

var Pool redis.Pool

func RedisInit() {
	pool := redis.Pool{
		MaxIdle:     16,
		MaxActive:   32,
		IdleTimeout: 120,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
	Pool = pool
}
