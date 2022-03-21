package middleware

import (
	"cloudDisk/src/util"
	"cloudDisk/src/util/db"
	"github.com/gin-gonic/gin"
)

// 权限控制

func Privilege() gin.HandlerFunc {
	return func(c *gin.Context) {
		sub := util.GetName(c)
		obj, act := util.GetOA(c)
		// 获取 enforcer 并加载策略
		enforcer := db.Enforcer
		enforcer.LoadPolicy()
		bl, _ := enforcer.Enforce(sub, obj, act)
		if bl {
			c.Next()
		} else {
			c.JSON(200, gin.H{
				"privilege": "fail",
			})
			c.Abort()
		}
	}

}
