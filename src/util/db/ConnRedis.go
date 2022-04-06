package db

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var Pool *redis.Pool

func init() {
	// 建立连接池
	Pool = &redis.Pool{
		MaxIdle:     50,
		MaxActive:   1000,
		IdleTimeout: 300 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", "127.0.0.1:6379",
				//redis.DialPassword("cloudGolang29"),
				//redis.DialDatabase(),
				//redis.DialConnectTimeout(),
				redis.DialReadTimeout(300*time.Second),
				redis.DialWriteTimeout(500*time.Second),
			)
			if err != nil {
				log.Fatalln("redis connect error!")
			}
			return con, nil
		},
	}
}
