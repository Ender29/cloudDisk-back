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
		//fmt.Println("进入拦截器")
		/*
		c.Header("access-control-allow-credentials", "true")
		c.Header("Access-Control-Allow-Origin", "*") //允许访问所有域
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, cache-control,postman-token,Cookie, Accept,x-requested-with")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Cache-Control,postman-token,Cookie, Accept,x-requested-with, No-Cache, If-Modified-Since, Pragma, Last-Modified, Expires, Access-Control-Allow-Credentials")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, OPTIONS")
		*/
		c.Header("Access-Control-Allow-Origin", "*")                                    // 这是允许访问所有域
            	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
          	//  header的类型
            	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
            //  允许跨域设置   可以返回其他子段
                c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
	        c.Header("Access-Control-Max-Age", "172800")  // 缓存请求信息 单位为秒
		c.Header("Access-Control-Allow-Credentials", "false")  //  跨域请求是否需要带cookie信息 默认设置为true
            	c.Set("content-type", "application/json")

		// err := r.ParseMultipartForm(32 << 20)
		err := c.Request.ParseMultipartForm(32 << 20)
		if err != nil {
			//fmt.Println(err.Error())
		}
		userName := c.Query("userName")
		//fmt.Println(userName)
		token := c.Query("userToken")
		//fmt.Println(token)
		// 验证token
		user := util.GetUser(userName)
		if c.Request.Method == "OPTIONS" {
			//c.JSON(http.StatusMethodNotAllowed, vo.NewRespMsg(-1, "不行哦~", nil))
			//c.Abort()
			c.JSON(http.StatusOK, vo.NewRespMsg(0, "ok", nil))
			c.Next()
		}else if user != nil && user.UserToken == token {
			c.Next()
		}else if user == nil {
			fmt.Println("没有用户")
			c.JSON(http.StatusMethodNotAllowed, vo.NewRespMsg(-1, "没有用户", nil))
			c.Abort()
		}else if user.UserToken != token {
			fmt.Println("状态无效")
			c.JSON(http.StatusMethodNotAllowed, vo.NewRespMsg(-1, "状态无效", nil))
			c.Abort()
		} else {
			c.JSON(http.StatusMethodNotAllowed, vo.NewRespMsg(-1, "不可以哦~", nil))
			c.Abort()
		}
	}
}
