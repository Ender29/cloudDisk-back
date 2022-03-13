package main

import (
	"cloudDisk/src/controller"
	"cloudDisk/src/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	root := router.Group("/")
	{
		root.POST("/login", controller.LoginHandler)
		root.POST("/register", controller.RegisterHandler)
		share := root.Group("/share")
		{
			share.GET("/download", controller.ShareDownloadHandler)
			share.GET("/filelist", controller.ShareFileListHandler)
			share.GET("/checkcode", controller.ShareCheckCodeHandler)
		}
	}
	file := root.Group("/file")
	file.Use(middleware.HTTPInterceptor(), middleware.Privilege())
	{
		file.GET("/closeshare", controller.ShareCloseHandler)
		file.GET("/createurl", controller.ShareCreateURLHandler)
		file.GET("/createfolder", controller.CreateFolderHandler)
		file.GET("/deletefile", controller.DeleteFileHandler)
		file.GET("/renamefile", controller.RenameFileHandler)
		file.GET("/movefile", controller.MoveFileHandler)
		file.GET("/searchbycategory", controller.SearchByCategoryHandler)
		file.GET("/filelist", controller.FileListHandler)
		file.GET("/sharelist", controller.ShareListHandler)
		file.GET("/check", controller.CheckFileHandler)
		file.POST("/upload", controller.UploadFileHandler)
		file.GET("/merge", controller.MergeFileHandler)
		file.GET("/download", controller.DownloadFileHandler)
	}
	if router.Run(":8080") != nil {
		fmt.Println("地球爆炸")
	}
	fmt.Println("服务监听。。。")
}
