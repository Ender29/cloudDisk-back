package main

import (
	"cloudDisk/src/controller"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/user/login", controller.LoginHandler)
	router.POST("/user/register", controller.RegisterHandler)
	router.Use(controller.HTTPInterceptor())
	router.POST("/user/logout", controller.LogoutHandler)
	router.POST("/user/changepwd", controller.ChangePwdHandler)
	router.POST("/file/createcatalog", controller.CcatalogHandler)
	router.POST("/file/deletefile", controller.DeleteFileHandler)
	router.POST("/file/renamefile", controller.RenameFileHandler)
	router.POST("/file/movefile", controller.MoveFileHandler)
	router.GET("/file/filelist", controller.FileListHandler)

	router.GET("/file/check", controller.CheckFileHandler)
	router.POST("/file/upload", controller.UploadFileHandler)
	router.GET("/file/merge", controller.MergeFileHandler)
	router.GET("/download", controller.DownloadFileHandler)
	if router.Run(":8080") != nil {
		fmt.Println("地球爆炸")
	}
	fmt.Println("服务监听。。。")
}
