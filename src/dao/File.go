package dao

import "database/sql"

type ShareList []ShareFile

type ShareFile struct {
	FileName string
	ShareAddr string
	ShareCode string
	SignupAt string
	Days int
}

type DownloadList struct {
	FilePath string
	FileName string
	FileMD5	 string
}

// FileMessage : 文件信息
type FileMessage struct {
	FileName string
	Category int8
	FilePath string
	FileTime string
	FileSize int
	Status   int8
}

type UniqueFile struct {
	Id       string
	FileSha1 string
	FileName string
	FileSize string
	FileAddr string
	CreateAt string
	UpdateAt string
	Status   string
	Ext1     string
}

type UserFile struct {
	Id         string
	ParentPath string
	FileName   string
	FileSha1   string
	FileSize   string
	Category   string
	UpdateAt   string
	ChangeTime string
	Ext1       string
	Ext2       sql.NullString
}