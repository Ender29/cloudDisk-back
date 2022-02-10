package controller

import (
	"archive/zip"
	"bytes"
	"cloudDisk/src/dao"
	"cloudDisk/src/service"
	"cloudDisk/src/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var uploadDir string = "D:\\upload\\"

// CheckFileHandler 判断文件是否存在
func CheckFileHandler(c *gin.Context) {
	fileMd5 := c.Query("fileMD5")

	isFolderExist, err := util.IsExist(uploadDir + fileMd5)
	if err != nil {
		fmt.Println("md5路径不存在")
	}
	var fileList []string
	chunks := 0
	if isFolderExist {
		tmp, err := util.ListDir(uploadDir + fileMd5, fileList)
		if err != nil {
			fmt.Println("获取文件列表失败")
		}
		chunks = len(tmp)
	} else {
		os.Mkdir(uploadDir + fileMd5, 0666)
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

	file := service.GetUserFileInfo(userName, parentPath, fileName)
	if file == nil {
		//w.Write([]byte("下载失败，没有此文件"))
		c.Writer.Write([]byte("下载失败，没有此文件"))
		return
	} else if file.Category == "5" {
		fmt.Println("this is a Dir")
		// 以文件夹名创建压缩文件
		//zipFile, err := os.Create(file.FileName + ".zip")
		//defer zipFile.Close()
		//if err != nil {
		//	fmt.Println("创建压缩文件失败")
		//}
		//---- 创建一个缓冲流
		buf := new(bytes.Buffer)

		zipWriter := zip.NewWriter(buf)

		// 获取文件夹下的文件写入压缩文件
		fileList := service.GetUserFileInfoList(userName, parentPath+file.FileName)
		//fileList := db.GetUserFileInfoList(userName, parentPath+file.FileName)
		for _, f := range fileList{
			fmt.Println(f.FileName)
			service.Compress(f, zipWriter, fileName, userName, parentPath+fileName)
		}

		zipWriter.Close()
		c.Header("Content-Type", "application/octect-stream")
		c.Header("content-disposition", "attachment; filename=\""+file.FileName+".zip"+"\"")
		c.Writer.Write(buf.Bytes())

	} else {
		uniqueFile := service.GetUniqueFileMeta(file.FileSha1)
		fmt.Println(uniqueFile.FileAddr, ":", uniqueFile.FileName)
		f, err := os.Open(uniqueFile.FileAddr + uniqueFile.FileName)
		defer f.Close()
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Write([]byte("下载失败"))
			return
		}
		data, _ := ioutil.ReadAll(f)
		c.Header("Content-Type", "application/octect-stream")
		c.Header("content-disposition", "attachment; filename=\""+fileName+"\"")
		c.Writer.Write(data)
	}

}

// MergeFileHandler 合并文件
func MergeFileHandler(c *gin.Context) {
	md5 := c.Query("fileMD5")
	fileName := c.Query("fileName")
	srcDir := uploadDir + md5
	var fileList []string
	fileList, err := util.ListDir(srcDir, fileList)
	if err != nil {
		fmt.Println("获取文件列表失败", err)
	}
	f, err := os.OpenFile(uploadDir + fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
	err  := ioutil.WriteFile(uploadDir + fileMD5 + "/" + current, data, 0666)
	if err != nil {
		fmt.Println("upload: ", err)
	}
	c.JSON(200, gin.H{
			"status": 0,
			"chunks": current,
		})
}
/*
func UploadFileHandler(c *gin.Context) {

	//if r.Method != "POST" {
	//	return
	//}

	c.Request.ParseMultipartForm(32 << 20)

	if c.MultipartForm == nil {
		fmt.Println("form is nil")
		return
	}

	values, _ := c.MultipartForm()
	fmt.Println("VALUE:", values.Value)
	userName := c.Query("userName")
	parent := values.Value["parentPath"][0]
	fmt.Println(userName)
	if !strings.HasSuffix(parent, "/") {
		parent += "/"
	}
	fmt.Println("FILE:", values.File)
	files := values.File["file"]

	for _, file := range files {

		uploadFile, err := file.Open()
		defer uploadFile.Close()
		if err != nil {
			fmt.Println("上传文件读取失败")
			return
		}
		fileContent, err := ioutil.ReadAll(uploadFile)
		if err != nil {
			fmt.Println("读取上传文件失败")
		}
		fileSha1 := service.Sha1(fileContent) // 计算文件sha1
		fmt.Println("FILEHASH:", fileSha1)
		uFile := service.GetUniqueFileMeta(fileSha1)
		path, fileName, pathNum := service.GetFilePathAndName(values.Value, file)
		fmt.Println("path:", path)
		fmt.Println("fileName:", fileName)
		fmt.Println("pathNum:", pathNum)

		// 没有存过这份文件
		if uFile == nil {
			// 存下文件
			// ---- 建立文件存放位置
			location := "/tmp/" // 全部文件存放地方
			fileAddr := service.BuildLocation(location, fileName)
			newFile, err := os.Create(fileAddr + fileName)
			defer newFile.Close()
			if err != nil {
				fmt.Println("文件存储失败")
				return
			}
			fileSize, err := newFile.Write(fileContent)
			//fileSize, err := io.Copy(newFile,uploadFile)
			if err != nil {
				fmt.Println("文件存储失败")
				return
			} else if fileSize == 0 {
				fmt.Println("写入为0")
			}
			// 更新唯一表
			service.InsertUniqueFileMeta(
				fileSha1,
				fileName,
				strconv.Itoa(fileSize),
				fileAddr,
			)
		}
		// 计算文件类型
		fileType := service.GetFileType(fileName)
		fmt.Println(path)
		parentPath := parent
		// 存过文件
		for index := 0; index < pathNum; index++ {
			//if db.CheckUserFileInfo(userName, parentPath, path[index]) {
			//	db.InsertUserFileInfo(userName, parentPath, path[index], "0", "0", "5")
			//}
			service.CheckAndInsertDirByLock(userName, parentPath, path[index])
			parentPath += path[index] + "/"
		}
		// 更新用户表
		if service.CheckUserFileInfo(userName, parentPath, fileName) {
			service.InsertUserFileInfo(userName, parentPath, fileName, fileSha1, strconv.FormatInt(file.Size, 10), strconv.Itoa(fileType))
		}

	}
	c.Writer.Write([]byte("上传成功"))

}*/

// CcatalogHandler : 创建目录
func CcatalogHandler(c *gin.Context) {
	userName := c.Query("userName")
	fileName := c.Query("fileName")
	//  父路径
	path := c.Query("path")
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	var message dao.FileMessage
	if userName != "" && path != "" && fileName != "" {
		thetime := time.Now().Format("2006-01-02 15:04:05")
		message.FileName = fileName
		message.FilePath = path
		message.FileTime = thetime
		message.Isdir = 1
		message.Category = 5
		message.FileID, message.Status = service.CreateCatalog(userName, fileName, thetime, path)
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
	fileID := c.Query("fileID")
	var message dao.FileMessage
	if userName != "" && fileID != "0" {
		message.FileID, _ = strconv.Atoi(fileID)
		service.DeleteFile(&message, userName)
	}
	if message.Category == 5 {
		message.Isdir = 1
	}
	message.FileTime = time.Now().Format("2006-01-02 15:04:05")
	data, _ := json.Marshal(message)
	c.Header("content-type", "text/json")
	c.Writer.Write(data)
}

// RenameFileHandler : 文件重命名
func RenameFileHandler(c *gin.Context) {
	var message dao.FileMessage
	id := c.Query("fileID")
	message.FileName = c.Query("oldName")
	newName := c.Query("newName")
	userName := c.Query("userName")
	path := c.Query("path")
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	message.FilePath = path
	category :=c.Query("category")
	if id != "0" && message.FileName != "" && userName != "" && newName != "" && category != "" && message.FilePath != "" {
		message.FileID, _ = strconv.Atoi(id)
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
	id := c.Query("fileID")
	message.FileName = c.Query("fileName")
	newPath := c.Query("newPath")
	if !strings.HasSuffix(newPath, "/") {
		newPath += "/"
	}
	userName := c.Query("userName")
	path := c.Query("path")
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	message.FilePath = path
	category := c.Query("category")
	if id != "0" && message.FileName != "" && newPath != "" && userName != "" && category != "" && message.FilePath != "" {
		message.FileID, _ = strconv.Atoi(id)
		temp, _ := strconv.Atoi(category)
		message.Category = int8(temp)
		service.MoveFile(userName, &message, newPath)
		data, _ := json.Marshal(message)
		c.Header("content-type", "text/json")
		c.Writer.Write(data)
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

