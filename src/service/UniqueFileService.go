package service

import (
	"archive/zip"
	"cloudDisk/src/dao"
	mysql "cloudDisk/src/util"
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type UniqueFile dao.UniqueFile

// GetUniqueFileMeta :根据hash从数据库获取一条信息
func GetUniqueFileMeta(fileSha1 string) *UniqueFile {
	stmt, err := mysql.DBConn().Prepare("select * from tbl_file where file_sha1=?")
	defer stmt.Close()
	if err != nil {
		fmt.Println("sql预编译失败")
		return nil
	}
	row := stmt.QueryRow(fileSha1)
	uFile := UniqueFile{}
	err = row.Scan(&uFile.Id, &uFile.FileSha1, &uFile.FileName, &uFile.FileSize, &uFile.FileAddr, &uFile.CreateAt, &uFile.UpdateAt, &uFile.Status, &uFile.Ext1)
	if err == sql.ErrNoRows {
		// 没有查询到任何一条数据
		return nil
	} else if err != nil {
		fmt.Println("获取唯一表数据失败--", err.Error())
		return nil
	}
	return &uFile

}

func InsertUniqueFileMeta(fileSha1, fileName, fileSize, fileAddr string) bool {
	stmt, err := mysql.DBConn().Prepare("insert into tbl_file(file_sha1, file_name, file_size, file_addr) values(?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println("sql预编译失败")
		return false
	}
	result, err := stmt.Exec(fileSha1, fileName, fileSize, fileAddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := result.RowsAffected(); rf < 1 && err == nil {
		return true
	}
	return false
}

// UpdateUserFileTable :将上传文件信息写入用户表
func InsertUserFileInfo(userName, parentPath, fileName, fileSha1, fileSize, category string) bool {
	// 更新插入文件夹信息

	stmt, err := mysql.DBConn().Prepare("insert into " + userName + "(parent_path, file_name, file_sha1, file_size, category) values(?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	result, err := stmt.Exec(parentPath, fileName, fileSha1, fileSize, category)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := result.RowsAffected(); rf != 0 && err == nil {
		return true
	}
	return false
}

// CheckUserFileInfo :查询用户是否有这条信息
func CheckUserFileInfo(userName, parent, fileName string) bool {
	stmt, err := mysql.DBConn().Prepare("select * from " + userName + " where parent_path=? and file_name=?")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	row := stmt.QueryRow(parent, fileName)
	userFile := dao.UserFile{}
	err = row.Scan(&userFile.Id, &userFile.ParentPath, &userFile.FileName, &userFile.FileSha1, &userFile.FileSize, &userFile.Category, &userFile.UpdateAt, &userFile.ChangeTime, &userFile.Ext1, &userFile.Ext2)
	if err == sql.ErrNoRows {
		return true
	}
	return false

}

// GetUserFileInfo : 获取一条用户文件信息
func GetUserFileInfo(userName, parent, fileName string) *dao.UserFile {
	if !strings.HasSuffix(parent, "/") {
		parent += "/"
	}
	stmt, err := mysql.DBConn().Prepare("select * from " + userName + " where parent_path=? and file_name=?")
	fmt.Println("select * from " + userName + " where parent_path='" + parent + "' and file_name=" + fileName)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	row := stmt.QueryRow(parent, fileName)
	userFile := dao.UserFile{}
	err = row.Scan(&userFile.Id, &userFile.ParentPath, &userFile.FileName, &userFile.FileSha1, &userFile.FileSize, &userFile.Category, &userFile.UpdateAt, &userFile.ChangeTime, &userFile.Ext1, &userFile.Ext2)
	if err == sql.ErrNoRows {
		fmt.Println("no rows")
		return nil
	}
	return &userFile
}

func GetUserFileInfoList(userName, parent string) []*dao.UserFile {
	if !strings.HasSuffix(parent, "/") {
		parent += "/"
	}
	var fileInfoList []*dao.UserFile
	stmt, err := mysql.DBConn().Prepare("select * from " + userName + " where parent_path=?")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	rows, err := stmt.Query(parent)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	for rows.Next() {
		fileInfo := dao.UserFile{}
		rows.Scan(&fileInfo.Id, &fileInfo.ParentPath, &fileInfo.FileName, &fileInfo.FileSha1, &fileInfo.FileSize, &fileInfo.Category, &fileInfo.UpdateAt, &fileInfo.ChangeTime, &fileInfo.Ext1, &fileInfo.Ext2)
		fileInfoList = append(fileInfoList, &fileInfo)
	}
	return fileInfoList
}

var lock sync.Mutex

// CheckAndInsertDirByLock :插入文件信息（分布式锁）
func CheckAndInsertDirByLock(userName, parentPath, fileName string) {
	lock.Lock()
	if CheckUserFileInfo(userName, parentPath, fileName) {
		InsertUserFileInfo(userName, parentPath, fileName, "0", "0", "5")
	}
	lock.Unlock()
}

type Sha1Stream struct {
	_sha1 hash.Hash
}

func (obj *Sha1Stream) Update(data []byte) {
	if obj._sha1 == nil {
		obj._sha1 = sha1.New()
	}
	obj._sha1.Write(data)
}

func (obj *Sha1Stream) Sum() string {
	return hex.EncodeToString(obj._sha1.Sum([]byte("")))
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func MultipartFileSha1(file multipart.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

// func PathExists(path string) (bool, error) {
// 	_, err := os.Stat(path)
// 	if err == nil {
// 		return true, nil
// 	}
// 	if os.IsNotExist(err) {
// 		return false, nil
// 	}
// 	return false, err
// }

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

func BuildLocation(location, fileName string) string {

	if !PathExists(location) {
		os.Mkdir(location, 0777)
	}
	if !PathExists(location + fileName) {
		return location
	}
	location += fmt.Sprintf("%v/", time.Now().Unix())
	return BuildLocation(location, fileName)
}

func GetFileSizeByMulti(file multipart.File) int {

	sum := 0
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		sum += n
		if err == io.EOF {
			break
		}
	}
	return sum
}

// GetFilePathAndName :解析文件相对路径和文件名
func GetFilePathAndName(values map[string][]string, file *multipart.FileHeader) ([]string, string, int){
	var pathList []string
	if v, ok := values["fullPath"]; ok {
		// 拖拽上传 fullPath=相对路径+文件名 name=文件名
		pathList = strings.Split(v[0], "/")
	}else {
		// 点击上传 name=相对路径+文件名
		pathList = strings.Split(file.Filename, "/")
	}
	pathLen := len(pathList)
	return pathList[:pathLen - 1], pathList[pathLen-1], pathLen - 1
}

func GetFileType(fileName string) int {
	fileNameSlice := strings.Split(fileName, ".")
	length := len(fileNameSlice)
	suffix := fileNameSlice[length - 1]
	switch suffix {
	case "txt","doc","docx","exl","ppt","hlp","rtf","html":
		return 0
	case "bmp","gif","jpg","pic","png","tif":
		return 1
	case "avi","mpg","mov","swf":
		return 2
	case "wav","aif","au","mp3","ram","wma","mmf":
		return 3
	}
	return 4
}

func Compress(file *dao.UserFile, zw *zip.Writer, prefix string, userName string, parentPath string) {

	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	if !strings.HasSuffix(parentPath, "/") {
		parentPath += "/"
	}

	// 文件夹
	if file.Category == "5" {
		fileList := GetUserFileInfoList(userName, parentPath+file.FileName)
		for _, f := range fileList {
			Compress(f, zw, prefix + file.FileName, userName, parentPath + file.FileName)
		}
		return
	}
	// 文件
	// 获取文件位置
	uFile := GetUniqueFileMeta(file.FileSha1)
	openFile, err := os.Open(uFile.FileAddr + uFile.FileName)
	if err != nil {
		fmt.Println("打开真实文件失败")
		return
	}
	info, err := openFile.Stat()
	header, err := zip.FileInfoHeader(info)
	header.Name = prefix + file.FileName
	writer, err := zw.CreateHeader(header)
	io.Copy(writer, openFile)
}


