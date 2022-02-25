package controller

import (
	"cloudDisk/src/dao"
	"cloudDisk/src/service"
	"cloudDisk/src/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)


// ShareCloseHandler : 取消分享
func ShareCloseHandler(c * gin.Context)  {
	shareAddr := c.Query("shareAddr")
	status := service.ShareClose(shareAddr)
	c.JSON(200, gin.H{
		"status": status,
	})
}

// ShareListHandler : 用户分享文件列表
func ShareListHandler(c *gin.Context) {
	userName := c.Query("userName")
	status, list := service.GetShareList(userName)
	c.JSON(200, gin.H{
		"status": status,
		"list": list,
	})
}

// SearchByCategoryHandler : 查找不同类型文件
func SearchByCategoryHandler(c *gin.Context) {
	userName := c.Query("userName")
	category := c.Query("category")
	files := service.SearchByCategory(userName, category)
	c.JSON(200, gin.H{
		"files": files,
	})
}

// ShareDownloadHandler : 下载分享文件
func ShareDownloadHandler(c * gin.Context) {
	shareName := c.Query("shareName")
	parentPath := c.Query("parentPath")
	fileName := c.Query("fileName")

	if !strings.HasSuffix(parentPath, "/") {
		parentPath += "/"
	}
	//service.IsDir(userName, fileName, parentPath)
	isDir, data := service.IsDir(shareName, fileName, parentPath)
	var rename string
	for _, char := range fileName {
		rename += "\\u" + fmt.Sprintf("%04x", char)
	}
	//fmt.Println(rename)
	if isDir {
		c.Header("content-disposition", "attachment; filename="+ rename + ".zip")
	} else {
		c.Header("content-disposition", "attachment; filename=" + rename)
	}
	c.Header("Content-Type", "application/octect-stream;")
	c.Writer.Write(data)
}

// ShareFileListHandler : 获取文件列表
func ShareFileListHandler(c *gin.Context) {
	shareName := c.Query("shareName")
	path := c.Query("parentPath")
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	parentPath := path
	var message []dao.FileMessage
	if shareName != "" && parentPath != "" {
		message = service.FileList(shareName, parentPath)
	}
	data, _ := json.Marshal(message)
	c.Header("content-type", "text/json")
	c.Writer.Write(data)
}

// ShareCheckCodeHandler : 校验提取码
func ShareCheckCodeHandler(c *gin.Context) {
	shareAddr := c.Query("shareAddr")
	shareCode := c.Query("shareCode")
	status, shareName, file := service.ShareCheckCode(shareAddr, shareCode)
	c.JSON(200, gin.H{
		"status": status,
		"shareName": shareName,
		"file": file,
	})
}

// ShareCreateURLHandler : 获取分享链接和提取码
func ShareCreateURLHandler(c *gin.Context) {
	parentPath := c.Query("parentPath")
	userName := c.Query("userName")
	fileName := c.Query("fileName")
	days := c.Query("days")
	addr, code := service.CreateURL(userName, fileName, parentPath, days)
	c.JSON(200, gin.H{
		"shareAddr": addr,
		"shareCode": code,
	})
}

// CheckFileHandler 判断文件是否存在
func CheckFileHandler(c *gin.Context) {
	fileMd5 := c.Query("fileMD5")

	isFolderExist, err := util.IsExist(service.UploadDir + fileMd5)
	if err != nil {
		fmt.Println("md5路径不存在")
	}
	var fileList []string
	chunks := 0
	if isFolderExist {
		tmp, err := util.ListDir(service.UploadDir + fileMd5, fileList)
		if err != nil {
			fmt.Println("获取文件列表失败")
		}
		chunks = len(tmp)
	} else {
		os.Mkdir(service.UploadDir + fileMd5, 0666)
	}
	c.JSON(200, gin.H{
		"status": 0,
		"chunks": chunks,
	})
}


// DownloadFileHandler :下载文件
func DownloadFileHandler(c *gin.Context) {
	userName := c.Query("userName")
	parentPath := c.Query("parentPath")
	fileName := c.Query("fileName")

	if !strings.HasSuffix(parentPath, "/") {
		parentPath += "/"
	}
	//service.IsDir(userName, fileName, parentPath)
	isDir, data := service.IsDir(userName, fileName, parentPath)
	var rename string
	for _, char := range fileName {
		rename += "\\u" + fmt.Sprintf("%04x", char)
	}
	//fmt.Println(rename)
	if isDir {
		c.Header("content-disposition", "attachment; filename="+ rename + ".zip")
	} else {
		c.Header("content-disposition", "attachment; filename=" + rename)
		//c.Header("content-disposition", "inline")
	}
	c.Header("Content-Type", "application/octect-stream;")
	c.Writer.Write(data)
}


// MergeFileHandler 合并文件
func MergeFileHandler(c *gin.Context) {
	fileMD5 := c.Query("fileMD5")
	fileName := c.Query("fileName")
	parentPath := c.Query("parentPath")
	fileSize := c.Query("fileSize")
	userName := c.Query("userName")
	srcDir := service.UploadDir + fileMD5
	var fileList []string
	fileList, err := util.ListDir(srcDir, fileList)
	fileSuffix := path.Ext(fileName)
	if err != nil {
		fmt.Println("获取文件列表失败", err)
	}
	uploadName := service.UploadDir + fileMD5 +"_file" + fileSuffix
	isExit, _ := util.IsExist(uploadName)
	if isExit {
		service.UploadFile(userName, fileMD5, fileName, parentPath, fileSize)
	} else {
		f, err := os.OpenFile(uploadName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("生成文件失败:", err)
		}
		// 合并区块
		length := len(fileList)
		for i := 0; i < length; i++ {
			data, _ := ioutil.ReadFile(srcDir + "/" + strconv.Itoa(i))
			f.Write(data)
		}
		defer f.Close()
		if !strings.HasSuffix(parentPath, "/") {
			parentPath += "/"
		}
		service.UploadFile(userName, fileMD5, fileName, parentPath, fileSize)
	}
	c.JSON(200, gin.H{
		"status": 0,
	})
}


// UploadFileHandler :处理文件上传
func UploadFileHandler(c *gin.Context) {
	c.Request.ParseMultipartForm(32 << 20)
	form, _ := c.MultipartForm()
	file := form.File["data"][0]
	current := form.Value["current"][0]
	fileMD5 := form.Value["fileMD5"][0]
	uploadFile, _ := file.Open()
	defer uploadFile.Close()
	data, _ :=ioutil.ReadAll(uploadFile)
	err  := ioutil.WriteFile(service.UploadDir + fileMD5 + "/" + current, data, 0666)
	if err != nil {
		fmt.Println("upload: ", err)
	}
	c.JSON(200, gin.H{
			"status": 0,
			"chunks": current,
		})
}


// CreateFolderHandler : 创建目录
func CreateFolderHandler(c *gin.Context) {
	userName := c.Query("userName")
	fileName := c.Query("fileName")
	//  父路径
	parentPath := c.Query("parentPath")
	if !strings.HasSuffix(parentPath, "/") {
		parentPath += "/"
	}
	var message dao.FileMessage
	if userName != "" && parentPath != "" && fileName != "" {
		curTime := time.Now().Format("2006-01-02 15:04:05")
		message.FileName = fileName
		message.FilePath = parentPath
		message.FileTime = curTime
		message.Category = 5
		message.Status = service.CreateCatalog(userName, fileName, curTime, parentPath)
	} else {
		message.Status = 4
	}
	data, _ := json.Marshal(message)
	c.Header("content-type", "text/json")
	c.Writer.Write(data)
}

// DeleteFileHandler : 删除文件
func DeleteFileHandler(c *gin.Context) {
	userName := c.Query("userName")
	var message dao.FileMessage
	message.FilePath = c.Query("filePath")
	message.FileName = c.Query("fileName")
	if userName != "" {
		service.DeleteFile(&message, userName)
	}
	message.FileTime = time.Now().Format("2006-01-02 15:04:05")
	data, _ := json.Marshal(message)
	c.Header("content-type", "text/json")
	c.Writer.Write(data)
}

// RenameFileHandler : 文件重命名
func RenameFileHandler(c *gin.Context) {
	var message dao.FileMessage
	message.FilePath = c.Query("filePath")
	message.FileName = c.Query("oldName")
	newName := c.Query("newName")
	userName := c.Query("userName")
	if !strings.HasSuffix(message.FilePath, "/") {
		message.FilePath  += "/"
	}
	category :=c.Query("category")
	if message.FileName != "" && userName != "" && newName != "" && category != "" && message.FilePath != "" {
		temp, _ := strconv.Atoi(category)
		message.Category = int8(temp)
		service.RenameFile(userName, &message, newName)
		data, _ := json.Marshal(message)
		c.Header("content-type", "text/json")
		c.Writer.Write(data)
	}
}

// MoveFileHandler : 移动文件
func MoveFileHandler(c *gin.Context) {
	var message dao.FileMessage
	message.FilePath = c.Query("filePath")
	message.FileName = c.Query("fileName")
	newPath := c.Query("newPath")
	if !strings.HasSuffix(newPath, "/") {
		newPath += "/"
	}
	userName := c.Query("userName")
	if !strings.HasSuffix(message.FilePath, "/") {
		message.FilePath += "/"
	}
	category := c.Query("category")
	if message.FileName != "" && newPath != "" && userName != "" && category != "" && message.FilePath != "" {
		if !(len(newPath) > len(message.FilePath) && message.FilePath[:len(message.FilePath)] == newPath[:len(message.FilePath)]) {
			fmt.Println("start move")
			temp, _ := strconv.Atoi(category)
			message.Category = int8(temp)
			service.MoveFile(userName, &message, newPath)
			data, _ := json.Marshal(message)
			c.Header("content-type", "text/json")
			c.Writer.Write(data)
		}
	}
}

// FileListHandler : 文件列表
func FileListHandler(c *gin.Context) {
	userName := c.Query("userName")
	path := c.Query("parentPath")
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	parentPath := path
	var message []dao.FileMessage
	if userName != "" && parentPath != "" {
		message = service.FileList(userName, parentPath)
	}
	data, _ := json.Marshal(message)
	c.Header("content-type", "text/json")
	c.Writer.Write(data)
}

