package controller

import (
	"cloudDisk/src/service"
	"cloudDisk/src/util"
	"cloudDisk/src/util/db"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

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

//	LogoutHandler : 注销账户
func LogoutHandler(c *gin.Context) {
	userName := util.GetName(c)
	userToken := c.Query("token")
	status := "0"
	params := make(map[string]string, 3)
	if userName != "" && userToken != "" {
		service.Logout(&userName, &userToken, &status)
	}
	params["userName"] = userName
	params["userToken"] = userToken
	params["status"] = status
	data, _ := json.Marshal(params)
	c.Header("content-type", "text/json")
	c.Writer.Write(data)
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
