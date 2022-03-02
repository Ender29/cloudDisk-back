package middleware

import (
	"cloudDisk/src/util"
	"fmt"
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
		if c.Request.Method == "OPTIONS" {
			fmt.Println("sadas")
			c.JSON(200, gin.H{
				"msg": "ok",
			})
			c.Next()
		} else if status == 1 {
			refreshToken := c.Request.Header.Get("RefreshToken")
			claims2, status2 := util.ParseToken(refreshToken)
			if status2 == 0 {
				newToken, _ := util.GenerateToken(claims2.Username, claims2.Password, time.Second*30)
				fmt.Println("time out")
				c.Header("Authorization", newToken)
				//c.Set("newToken", newToken)
				c.Set("userName", claims2.Username)
				c.Next()
			} else {
				fmt.Println("登录已过期")
				c.JSON(200, gin.H{
					"time": 1,
					"msg":  "登录已过期",
				})
				c.Abort()
			}
		} else if status == 0 {
			fmt.Println("there1")
			c.Set("userName", claims.Username)
			c.Next()
		} else {
			fmt.Println("claims: ", claims)
			fmt.Println("status: ", status)
			fmt.Println("there2")
			c.Abort()
		}
	}
}
