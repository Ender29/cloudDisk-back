package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var rc *redis.Pool

func init() {
	rc = &redis.Pool{
		MaxIdle:     50,
		MaxActive:   1000,
		IdleTimeout: 300 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", "124.223.78.104:9039") //redis.DialPassword(),
			//redis.DialDatabase(),
			//redis.DialConnectTimeout(),
			//redis.DialReadTimeout(timeout*time.Second),
			//redis.DialWriteTimeout(timeout*time.Second),

			if err != nil {
				log.Fatalln("redis connect error!")
			}
			return con, nil
		},
	}
}

func main() {
	conn := rc.Get()
	defer conn.Close()
	//conn.Do("set", "test", "hello")
	//reply, err := conn.Do("get", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImVuZGVyIiwicGFzc3dvcmQiOiJ6dkwwNVh3aXZJSGp3UFVkMXNEZ2ZBbVJhaDg4eWtHTkx6TXNNajRZSzZhaFZWWlhkRk1paGNwRHh1MnBGTmd0IiwiZXhwIjoxNjQ2NzE3MTI1LCJpc3MiOiJsZXQncyBnbyJ9.h6f2FwXIzHbwmFBY0_5lyZYkuMArKe-slK5IePoPc9I")
	//if err != nil {
	//	log.Fatalln("get error")
	//}
	//s, _ := redis.String(reply, nil)
	//fmt.Println(s)
	_, err1 := conn.Do("set", "keyWithEx", "valueWithEx", 3600*24)
	if err1 != nil {
		fmt.Println("err:", err1)
	}
	valWithEx, _ := redis.String(conn.Do("get", "keyWithEx"))
	if valWithEx != "" {
		fmt.Println("获取带过期时间的key，", valWithEx)
	}
	time.Sleep(3 * time.Second) // 模拟过期时间
	valWithEx1, _ := redis.String(conn.Do("get", "keyWithEx"))
	if valWithEx1 == "" {
		fmt.Println("空")
	} else {
		fmt.Println("获取带过期时间的key，", valWithEx)
	}
}
