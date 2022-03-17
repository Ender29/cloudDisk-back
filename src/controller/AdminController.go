package controller

import (
	"cloudDisk/src/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

// FilesHandler 获取所有文件信息列表
func FilesHandler(c *gin.Context) {
	list := service.GetFiles()
	c.JSON(200, gin.H{
		"list": list,
	})
}

// SharesHandler 获取分享文件信息列表
func SharesHandler(c *gin.Context) {
	list := service.GetShares()
	fmt.Println(list)
	c.JSON(200, gin.H{
		"list": list,
	})
}

// UsersHandler 获取分享文件信息列表
func UsersHandler(c *gin.Context) {
	list := service.GetUsers()
	c.JSON(200, gin.H{
		"list": list,
	})
}

// PoliciesHandler 获取所有政策
func PoliciesHandler(c *gin.Context) {
	list := service.GetPolicies()
	c.JSON(200, gin.H{
		"list": list,
	})
}
