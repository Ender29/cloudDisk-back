package middleware

import (
	"cloudDisk/src/util"
	"cloudDisk/src/util/db"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"time"
)

// HTTPInterceptor :拦截请求，验证token
func HTTPInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
		c.Header("Access-Control-Max-Age", "172800")
		c.Header("Access-Control-Allow-Credentials", "false")
		c.Set("content-type", "application/json")

		accessToken := c.Request.Header.Get("AccessToken")
		claims, status := util.ParseToken(accessToken)
		//fmt.Println("status:", status)
		if c.Request.Method == "OPTIONS" {
			c.JSON(200, gin.H{
				"msg": "ok",
			})
			return
		} else if status == 1 {
			//refreshToken := c.Request.Header.Get("RefreshToken")
			conn := db.Pool.Get()
			defer conn.Close()
			refreshToken, _ := redis.String(conn.Do("get", accessToken))
			//fmt.Println("refresh token:", refreshToken)
			claims2, status2 := util.ParseToken(refreshToken)
			if status2 == 0 && refreshToken != "" {
				newToken, _ := util.GenerateToken(claims2.Username, claims2.Password, time.Minute*5)
				conn.Do("set", newToken, refreshToken, "EX", 3600*24)
				c.Header("Authorization", newToken)
				c.Set("userName", claims2.Username)
				c.Next()
			} else {
				c.JSON(200, gin.H{
					"timeout": 1,
					"msg":     "登录已过期",
				})
				c.Abort()
			}
			conn.Do("del", accessToken)
		} else if status == 0 {
			c.Set("userName", claims.Username)
			c.Next()
		} else {
			c.Abort()
		}
	}
}
