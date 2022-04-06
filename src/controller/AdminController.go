package controller

import (
	"cloudDisk/src/dao"
	"cloudDisk/src/service"
	"cloudDisk/src/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"strings"
)

// AdminUpload 管理员上传
func AdminUpload(c *gin.Context) {
	c.Request.ParseMultipartForm(32 << 20)
	form, _ := c.MultipartForm()
	file := form.File["data"][0]
	current := form.Value["current"][0]
	fileMD5 := form.Value["fileMD5"][0]
	uploadFile, _ := file.Open()
	defer uploadFile.Close()
	data, _ := ioutil.ReadAll(uploadFile)
	err := ioutil.WriteFile(service.UploadDir+fileMD5+"/"+current, data, 0666)
	if err != nil {
		fmt.Println("upload: ", err)
	}
	c.JSON(200, gin.H{
		"status": 0,
		"chunks": current,
	})
}

// AdminMerge 管理员合并
func AdminMerge(c *gin.Context) {
	fileMD5 := c.Query("fileMD5")
	uploadName := service.UploadDir + fileMD5 + "_file"
	isExit, _ := util.IsExist(uploadName)
	if !isExit {
		service.MergeFile(uploadName, fileMD5)
	}
	c.JSON(200, gin.H{
		"status": 0,
	})
}

// AdminCheck 统计文件块
func AdminCheck(c *gin.Context) {
	fileMd5 := c.Query("fileMD5")
	isFolderExist, err := util.IsExist(service.UploadDir + fileMd5)
	if err != nil {
		fmt.Println("md5路径不存在")
	}
	var fileList []string
	chunks := 0
	if isFolderExist {
		tmp, err := util.ListDir(service.UploadDir+fileMd5, fileList)
		if err != nil {
			fmt.Println("获取文件列表失败")
		}
		chunks = len(tmp)
	} else {
		os.Mkdir(service.UploadDir+fileMd5, 0666)
	}
	c.JSON(200, gin.H{
		"status": 0,
		"chunks": chunks,
	})
}

// DelFileHandler 删除文件
func DelFileHandler(c *gin.Context) {
	fileMD5 := c.PostForm("fileMD5")
	service.DelFile(fileMD5)
	c.JSON(200, gin.H{
		"msg":  "取消成功",
		"type": "success",
	})
}

// CloseShareHandler 取消分享
func CloseShareHandler(c *gin.Context) {
	shareAddr := c.Query("shareAddr")
	service.ShareClose(shareAddr)
	c.JSON(200, gin.H{
		"msg":  "取消成功",
		"type": "success",
	})
}

//LogoutHandler 注销账户
func LogoutHandler(c *gin.Context) {
	userName := c.PostForm("userName")
	service.Logout(userName)
	c.JSON(200, gin.H{
		"msg":  "注销成功",
		"type": "success",
	})
}

// SetRoleHandler 设置角色
func SetRoleHandler(c *gin.Context) {
	oldp := c.PostForm("oldPolicy")
	newp := c.PostForm("newPolicy")
	oldPolicy := strings.Split(oldp, ",")
	newPolicy := strings.Split(newp, ",")
	bl := service.SetRoleService(oldPolicy, newPolicy)
	c.JSON(200, gin.H{
		"msg": bl,
	})
}

// DelPolicyHandler 删除policy
func DelPolicyHandler(c *gin.Context) {
	p := dao.Policy{}
	p.Sub = c.Query("sub")
	p.Obj = c.Query("obj")
	p.Act = c.Query("act")
	bl := service.DelPolicyService(p)
	c.JSON(200, gin.H{
		"msg": bl,
	})
}

// SetPolicyHandler 修改policy
func SetPolicyHandler(c *gin.Context) {
	oldp := c.PostForm("oldPolicy")
	newp := c.PostForm("newPolicy")
	oldPolicy := strings.Split(oldp, ",")
	newPolicy := strings.Split(newp, ",")
	bl := service.SetPolicyService(oldPolicy, newPolicy)
	c.JSON(200, gin.H{
		"msg": bl,
	})
}

// AddPolicyHandler 新增policy
func AddPolicyHandler(c *gin.Context) {
	p := dao.Policy{}
	p.Sub = c.PostForm("sub")
	p.Obj = c.PostForm("obj")
	p.Act = c.PostForm("act")
	bl := service.AddPolicyService(p)
	c.JSON(200, gin.H{
		"msg": bl,
	})
}

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
