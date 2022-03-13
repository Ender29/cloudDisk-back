package util

import "github.com/gin-gonic/gin"

func GetName(c *gin.Context) string {
	name, bl := c.Get("userName")
	if bl {
		return name.(string)
	}
	return ""
}

func GetOA(c *gin.Context) (string, string) {
	return c.Request.URL.Path, c.Request.Method
}
