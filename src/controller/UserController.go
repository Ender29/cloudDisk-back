package controller

import (
	"cloudDisk/src/service"
	"cloudDisk/src/util"
	"cloudDisk/src/util/db"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"path"
	"time"
)

// UploadRoleHandler 升级权限
func UploadRoleHandler(c *gin.Context) {
	username := util.GetName(c)
	role := c.PostForm("role")
	if role == "admin" {
		c.JSON(200, gin.H{
			"msg":  "升级失败！",
			"type": "error",
		})
		return
	}
	service.UpRoleService(username, role)
	c.JSON(200, gin.H{
		"msg":  "升级成功！",
		"type": "success",
	})
}

// UploadPhotoHandler 修改头像
func UploadPhotoHandler(c *gin.Context) {
	form, _ := c.MultipartForm()
	img := form.File["img"][0]
	act := "fail"
	img64 := ""
	fileSuffix := path.Ext(img.Filename)
	if img.Size < 10<<10 && fileSuffix == ".png" {
		img64 = service.UploadPhotoService(util.GetName(c), img)
		act = "success"
	}
	c.JSON(200, gin.H{
		"act": act,
		"img": img64,
	})
}

// LoginHandler 登录
func LoginHandler(c *gin.Context) {
	userName := c.PostForm("userName")
	userPwd := c.PostForm("userPwd")
	var message service.LoginMessage
	message.UserName = userName
	//fmt.Println(userName)
	if userName != "" && userPwd != "" {
		userPwd, _ = util.EnPwdCode([]byte(userPwd))
		//message.UserToken, message.LatestTime, message.Status, message.FileSize = service.Login(userName, userPwd)
		message.Login(userPwd)
		enforcer := db.Enforcer
		enforcer.LoadPolicy()
		if enforcer.HasNamedGroupingPolicy("g", message.UserName, "admin") {
			message.Role = "admin"
			message.Status = 0
		} else {
			message.Role = "user"
		}
	} else {
		message.Status = -1
	}
	// 返回字段
	//data, _ := json.Marshal(message)
	c.Header("content-type", "text/json")
	c.JSON(200, message)
}

// RegisterHandler : 用户注册
func RegisterHandler(c *gin.Context) {
	userName := c.PostForm("userName")
	userPwd := c.PostForm("userPwd")
	var message service.RegisterMessage
	message.UserName = userName
	if userName != "" && userPwd != "" {
		message.UserName = userName
		message.Register(userPwd)
		if message.Status == 0 {
			enforcer := db.Enforcer
			enforcer.LoadPolicy()
			enforcer.AddGroupingPolicy(userName, "normal")
			message.RegisterTime = time.Now().Format("2006-01-02 15:04:05")
		}
	} else {
		message.Status = 4
	}
	c.Header("content-type", "text/json")
	c.JSON(200, message)
}

// ChangePwdHandler : 修改密码
func ChangePwdHandler(c *gin.Context) {
	userName := util.GetName(c)
	userPwd := c.Query("userPwd")
	newPwd := c.Query("newPwd")
	status := "0"
	params := make(map[string]string, 2)
	if userName != "" && userPwd != "" && newPwd != "" {
		userPwd, _ = util.EnPwdCode([]byte(userPwd))
		newPwd, _ = util.EnPwdCode([]byte(newPwd))
		status = service.ChangePwd(userName, userPwd, newPwd)
	}
	params["userName"] = userName
	params["status"] = status
	data, _ := json.Marshal(params)
	c.Header("content-type", "text/json")
	c.Writer.Write(data)
}
