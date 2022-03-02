package main

import (
	"cloudDisk/src/controller"
	"cloudDisk/src/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/user/login", controller.LoginHandler)
	router.POST("/user/register", controller.RegisterHandler)
	router.GET("/share/download", controller.ShareDownloadHandler)
	router.GET("/share/filelist", controller.ShareFileListHandler)
	router.GET("/share/checkcode", controller.ShareCheckCodeHandler)
	router.Use(middleware.HTTPInterceptor())
	//router.POST("/user/logout", controller.LogoutHandler)
	//router.POST("/user/changepwd", controller.ChangePwdHandler)
	router.GET("/share/searchbycategory", controller.SearchByCategoryHandler)
	router.GET("/share/createurl", controller.ShareCreateURLHandler)
	router.GET("/file/closeshare", controller.ShareCloseHandler)
	router.GET("/file/createfolder", controller.CreateFolderHandler)
	router.GET("/file/deletefile", controller.DeleteFileHandler)
	router.GET("/file/renamefile", controller.RenameFileHandler)
	router.GET("/file/movefile", controller.MoveFileHandler)
	router.GET("/file/filelist", controller.FileListHandler)
	router.GET("/file/sharelist", controller.ShareListHandler)
	router.GET("/file/check", controller.CheckFileHandler)
	router.POST("/file/upload", controller.UploadFileHandler)
	router.GET("/file/merge", controller.MergeFileHandler)
	router.GET("/file/download", controller.DownloadFileHandler)
	if router.Run(":8080") != nil {
		fmt.Println("地球爆炸")
	}
	fmt.Println("服务监听。。。")
}
