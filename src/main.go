package main

import (
	"cloudDisk/src/controller"
	"cloudDisk/src/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
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
	root.Use(middleware.HTTPInterceptor(), middleware.Privilege())
	file := root.Group("/file")
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
		file.POST("/preview", controller.PreviewFileHandler)
		file.GET("/searchbyname", controller.SearchByNameHandler)
	}
	admin := root.Group("/admin")
	{
		admin.GET("/filelist", controller.FilesHandler)
		admin.GET("/shares", controller.SharesHandler)
		admin.GET("/users", controller.UsersHandler)
		admin.GET("/policies", controller.PoliciesHandler)
		admin.POST("/addpolicy", controller.AddPolicyHandler)
		admin.POST("/setpolicy", controller.SetPolicyHandler)
		admin.GET("/delpolicy", controller.DelPolicyHandler)
		admin.POST("/setrole", controller.SetRoleHandler)
		admin.POST("/logout", controller.LogoutHandler)
		admin.GET("/closeshare", controller.CloseShareHandler)
		admin.POST("/delfile", controller.DelFileHandler)
		admin.POST("/upload", controller.AdminUpload)
		admin.GET("/check", controller.AdminCheck)
		admin.GET("/merge", controller.AdminMerge)
	}
	user := root.Group("/user")
	{
		user.POST("/uploadphoto", controller.UploadPhotoHandler)
		user.POST("/uprole", controller.UploadRoleHandler)
	}
	if router.Run(":8080") != nil {
		fmt.Println("地球爆炸")
	}
	fmt.Println("服务监听。。。")
}
