package dao

type ShareList []ShareFile

type ShareFile struct {
	FileName  string
	ShareAddr string
	ShareCode string
	SignupAt  string
	Days      int
}

type DownloadList struct {
	FilePath string
	FileName string
	FileMD5  string
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

type FilesList struct {
	FileMD5    string
	FileSize   int
	CreateTime string
	UpdateTime string
}
