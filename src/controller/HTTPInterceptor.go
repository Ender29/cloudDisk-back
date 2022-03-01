package controller

import (
	"cloudDisk/src/util"
	"cloudDisk/src/vo"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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

		// err := r.ParseMultipartForm(32 << 20)
		err := c.Request.ParseMultipartForm(32 << 20)
		if err != nil {
			//fmt.Println(err.Error())
		}
		userName := c.Query("userName")
		token := c.Query("userToken")

		test := c.Request.Header.Get("token")
		fmt.Println("token: ", test)
		//claims, status := util.ParseToken(test)

		// 验证token
		user := util.GetUser(userName)
		if c.Request.Method == "OPTIONS" {
			//c.JSON(http.StatusMethodNotAllowed, vo.NewRespMsg(-1, "不行哦~", nil))
			//c.Abort()
			c.JSON(http.StatusOK, vo.NewRespMsg(0, "ok", nil))
			c.Next()
		} else if user != nil && user.UserToken == token {
			c.Next()
		} else if user == nil {
			fmt.Println("没有用户")
			c.JSON(http.StatusMethodNotAllowed, vo.NewRespMsg(-1, "没有用户", nil))
			c.Abort()
		} else if user.UserToken != token {
			fmt.Println("状态无效")
			c.JSON(http.StatusMethodNotAllowed, vo.NewRespMsg(-1, "状态无效", nil))
			c.Abort()
		} else {
			c.JSON(http.StatusMethodNotAllowed, vo.NewRespMsg(-1, "不可以哦~", nil))
			c.Abort()
		}
	}
}
