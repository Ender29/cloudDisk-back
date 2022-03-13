package service

import (
	"archive/zip"
	"bytes"
	"cloudDisk/src/dao"
	util "cloudDisk/src/util"
	"cloudDisk/src/util/db"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const UploadDir = "D:/upload/"

// ShareClose : 取消分享
func ShareClose(shareAddr string) int {
	sql := "delete from tbl_share where share_addr='" + shareAddr + "'"
	//fmt.Println(sql)
	_, err := db.DBConn().Exec(sql)
	if err != nil {
		return -1
	}
	return 0
}

// GetShareList : 查询用户已分享的文件
func GetShareList(userName string) (int8, dao.ShareList) {
	sql := "select a.share_addr,a.share_code,a.signup_at,(a.days-DATEDIFF(now(),a.signup_at)),b.file_name from tbl_share as a inner join " + userName + "_share as b;"
	var list dao.ShareList
	rows, err := db.DBConn().Query(sql)
	if err != nil {
		return 1, list
	}
	for rows.Next() {
		var file dao.ShareFile
		err = rows.Scan(&file.ShareAddr, &file.ShareCode, &file.SignupAt, &file.Days, &file.FileName)
		if err != nil {
			return 2, list
		}
		if file.Days < 0 {
			file.Days = 0
		}
		list = append(list, file)
	}
	return 0, list
}

// SearchByCategory : 通过类型查找文件
func SearchByCategory(userName, category string) []dao.FileMessage {
	sql := "select parent_path,file_name,category,change_time,file_size from " + userName + " where category='" + category + "'"
	rows, _ := db.DBConn().Query(sql)
	var list []dao.FileMessage
	for rows.Next() {
		var file dao.FileMessage
		err := rows.Scan(&file.FilePath, &file.FileName, &file.Category, &file.FileTime, &file.FileSize)
		if err != nil {
			return list
		} else {
			list = append(list, file)
		}
	}
	return list
}

// ShareCheckCode : 检查分享地址和分享码
func ShareCheckCode(shareAddr, shareCode string) (int8, string, dao.FileMessage) {
	sql := "delete from tbl_share where share_addr='" + shareAddr + "' and days<=DATEDIFF(now(),signup_at)"
	res, _ := db.DBConn().Exec(sql)
	eff, _ := res.RowsAffected()
	var file dao.FileMessage
	if eff > 0 {
		return 1, "", file
	}
	sql = "select share_code,share_name from tbl_share where share_addr='" + shareAddr + "'"
	code := ""
	shareName := ""
	row := db.DBConn().QueryRow(sql)
	row.Scan(&code, &shareName)
	if len(code) < 1 {
		return 2, "", file
	}
	if code != shareCode {
		return 3, "", file
	}
	sql = "select parent_path,file_name,category,change_time,file_size from " + shareName + " where (parent_path,file_name) in (select parent_path,file_name from " + shareName + "_share where share_addr='" + shareAddr + "')"
	row = db.DBConn().QueryRow(sql)

	row.Scan(&file.FilePath, &file.FileName, &file.Category, &file.FileTime, &file.FileSize)

	return 0, shareName, file
}

// CreateURL 获得分享地址和分享码
func CreateURL(userName, fileName, parentPath, days string) (string, string) {

	sql := "select share_addr from " + userName + "_share where parent_path='" + parentPath + "' and file_name='" + fileName + "'"
	var findAddr = ""
	row := db.DBConn().QueryRow(sql)
	row.Scan(&findAddr)
	addr := ""
	code := ""
	if findAddr == "" {
		for i := 0; i < 10; i++ {
			addr = util.GetRandStr(15)
			code = util.GetRandStr(4)
			sql = "insert into tbl_share (share_addr,share_name,share_code,days) values('" + addr + "','" + userName + "','" + code + "','" + days + "')"
			stmt, _ := db.DBConn().Prepare(sql)
			_, err := stmt.Exec()
			if err == nil {
				sql = "insert into " + userName + "_share (parent_path,file_name,share_addr) values('" + parentPath + "','" + fileName + "','" + addr + "')"
				fmt.Println(sql)
				db.DBConn().Exec(sql)
				break // 插入成功就中断循环
			}
			fmt.Println("err:", err)
		}
	}

	return addr, code

}

// IsDir : 判断是否是文件夹
func DownloadService(userName, fileName, parentPath string) (bool, []byte) {
	sql := "select category,file_md5 from " + userName + " where file_name='" + fileName + "' and parent_path='" + parentPath + "' limit 1"
	//fmt.Println(sql)
	// warning: 没查到也可能是0
	row := db.DBConn().QueryRow(sql)
	var category int
	var fileMD5 string
	row.Scan(&category, &fileMD5)
	//fmt.Println(category)
	// 判断是否是文件夹
	if category == 5 {
		time := time.Now().UnixNano()
		parPath := parentPath + fileName + "/"
		pathLen := len(parPath)
		sql = "SELECT parent_path,file_name,file_md5 FROM " + userName + " where MID(parent_path,1," + strconv.Itoa(pathLen) + ")='" + parPath + "' AND category!='5'"
		//fmt.Println(sql)
		rows, _ := db.DBConn().Query(sql)
		var list []dao.DownloadList
		for rows.Next() {
			var file_name string
			var file_md5 string
			var parent_path string
			err := rows.Scan(&parent_path, &file_name, &file_md5)
			if err == nil {
				list = append(list, dao.DownloadList{
					FileName: file_name,
					FilePath: parent_path,
					FileMD5:  file_md5,
				})
			}
		}
		downPath := UploadDir + userName + strconv.FormatInt(time, 10)
		os.MkdirAll(downPath+"/", 0666)
		for i := range list {
			dirPath := downPath + list[i].FilePath
			fileSuffix := path.Ext(list[i].FileName)
			os.MkdirAll(dirPath, 0666)
			util.Copy(UploadDir+list[i].FileMD5+"_file"+fileSuffix, dirPath+list[i].FileName)
		}
		// 预防：旧文件无法覆盖
		//zipName := downPath + "/" + fileName + ".zip"
		//os.RemoveAll(zipName)
		buf := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buf)
		util.Zip(downPath+"/", zipWriter)
		ioutil.WriteFile(downPath+".zip", buf.Bytes(), 0666)
		os.RemoveAll(downPath)
		os.Remove(downPath + ".zip")
		return true, buf.Bytes()
	}
	// 文件直接读取数据
	file, _ := os.Open(UploadDir + fileMD5 + "_file" + path.Ext(fileName))
	defer file.Close()
	data, _ := io.ReadAll(file)
	return false, data
}

// UploadFile : 上传文件
func UploadFile(userName, fileMD5, fileName, parentPath, fileSize string) int8 {
	fileSuffix := path.Ext(UploadDir + fileMD5 + "_" + fileName)
	fileSuffix = strings.ToLower(fileSuffix)
	var category int8 = 0
	if fileSuffix == ".wav" || fileSuffix == ".mp3" || fileSuffix == ".au" || fileSuffix == ".aif" || fileSuffix == ".aiff" || fileSuffix == ".ra" || fileSuffix == ".mid" {
		category = 1
	} else if fileSuffix == ".avi" || fileSuffix == ".mp4" || fileSuffix == ".mkv" || fileSuffix == ".wmv" || fileSuffix == ".3gp" || fileSuffix == ".mod" || fileSuffix == ".mov" || fileSuffix == ".ogg" || fileSuffix == ".rm" || fileSuffix == ".rmvb" || fileSuffix == ".dat" || fileSuffix == ".webm" {
		category = 2
	} else if fileSuffix == ".png" || fileSuffix == ".gif" || fileSuffix == ".jpg" || fileSuffix == ".raw" || fileSuffix == ".bmp" || fileSuffix == ".tiff" || fileSuffix == ".psd" || fileSuffix == ".svg" {
		category = 3
	} else if fileSuffix == ".xls" || fileSuffix == ".txt" || fileSuffix == ".xlsx" || fileSuffix == ".csv" || fileSuffix == ".ppt" || fileSuffix == ".doc" || fileSuffix == ".docx" || fileSuffix == ".pptx" {
		category = 4
	}
	// warning:文件名带单引号
	sql := "insert into tbl_file (file_md5,file_name,file_size) values('" + fileMD5 + "',\"" + fileName + "\",'" + fileSize + "')"
	//fmt.Println(sql)
	stmt, err := db.DBConn().Prepare(sql)
	if err != nil {
		fmt.Println("UploadFile insert tbl_file:", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println("UploadFile insert tbl_file:", err)
	}
	if parentPath != "/" {
		paths := strings.Split(parentPath, "/")
		root := "/"
		for i := 0; i < len(paths); i++ {
			if paths[i] == "" {
				continue
			}
			status := CreateCatalog(userName, paths[i], time.Now().Format("2006-01-02 15:04:05"), root)
			fmt.Println("CreateCatalog: ", status)
			root += paths[i] + "/"
		}
	}
	if !strings.HasSuffix(parentPath, "/") {
		parentPath += "/"
	}
	sql = "INSERT INTO " + userName + " (parent_path,file_name,category,file_size,file_md5) values('" + parentPath + "','" + fileName + "','" + strconv.Itoa(int(category)) + "','" + fileSize + "','" + fileMD5 + "')"
	fmt.Println(sql)
	stmt, _ = db.DBConn().Prepare(sql)
	_, err = stmt.Exec()
	if err != nil {
		return 2
	}
	return 0
}

// CreateCatalog : 创建文件夹
func CreateCatalog(userName, fileName, theTime, path string) int8 {
	var status int8 = 0
	sql := "INSERT INTO " + userName + "(parent_path,file_name,category,upload_at,change_time) VALUES ('" + path + "','" + fileName + "','5','" + theTime + "','" + theTime + "')"
	fmt.Println(sql)
	stmt, _ := db.DBConn().Prepare(sql)
	_, err := stmt.Exec()
	if err != nil {
		fmt.Println(err)
		status = 2
	}
	return status
}

// DeleteFile : 删除文件
func DeleteFile(message *dao.FileMessage, userName string) {
	sql := "select parent_path, file_name, file_size, category from " + userName + " where parent_path=? and file_name=?"
	rows, err := db.DBConn().Query(sql, message.FilePath, message.FileName)
	if err != nil {
		message.Status = 1
	}
	for rows.Next() {
		// 查询文件信息
		err = rows.Scan(&message.FilePath, &message.FileName, &message.FileSize, &message.Category)
		if err != nil {
			message.Status = 2
		}
	}
	if message.Status == 0 {
		sql = "delete from tbl_share where share_addr=(select share_addr from " + userName + "_share where parent_path='" + message.FilePath + "' and file_name='" + message.FileName + "')"
		fmt.Println(sql)
		db.DBConn().Exec(sql)
		sql = "delete from " + userName + " where parent_path=? and file_name=?"
		stmt, _ := db.DBConn().Prepare(sql)
		_, err = stmt.Exec(message.FilePath, message.FileName)
		if err != nil {
			message.Status = 3
		}
		if message.Category == 5 {
			str := message.FilePath + message.FileName
			if !strings.HasSuffix(str, "/") {
				str += "/"
			}
			length := strconv.Itoa(utf8.RuneCountInString(str))
			// 先删了分享表中的
			sql = "delete from tbl_share where share_addr=(select share_addr from " + userName + "_share where mid(parent_path, 1, " + length + ")='" + str + "')"
			fmt.Println(sql)
			db.DBConn().Exec(sql)
			sql = "delete from " + userName + " where mid(parent_path, 1, " + length + ")='" + str + "'"
			stmt, _ = db.DBConn().Prepare(sql)
			_, err = stmt.Exec()
			if err != nil {
				message.Status = 4
			}
		}
	}
}

// RenameFile : 文件重命名
func RenameFile(userName string, message *dao.FileMessage, newName string) {
	sql := "update " + userName + " set file_name=? where parent_path=? and file_name=?"
	stmt, _ := db.DBConn().Prepare(sql)
	_, err := stmt.Exec(newName, message.FilePath, message.FileName)
	if err != nil {
		message.Status = 1
	}
	if message.Category == 5 {
		insertStart := strconv.Itoa(utf8.RuneCountInString(message.FilePath) + 1)
		insertLen := strconv.Itoa(utf8.RuneCountInString(message.FileName))
		midLen := strconv.Itoa(utf8.RuneCountInString(message.FileName) + utf8.RuneCountInString(message.FilePath) + 1)
		sql = "update " + userName + " set parent_path=insert(parent_path, " + insertStart + ", " + insertLen + ",'" + newName + "') where mid(parent_path, 1, " + midLen + ")='" + message.FilePath + message.FileName + "/'"
		fmt.Println(sql)
		_, err = db.DBConn().Exec(sql)
		if err != nil {
			message.Status = 2
		}
	}
	if message.Status == 0 {
		message.FileName = newName
	}
}

// MoveFile ：移动文件
func MoveFile(userName string, message *dao.FileMessage, newPath string) {
	sql := "update " + userName + " set parent_path=? where parent_path=? and file_name=?"
	fmt.Println(sql)
	stmt, _ := db.DBConn().Prepare(sql)
	_, err := stmt.Exec(newPath, message.FilePath, message.FileName)
	if err != nil {
		message.Status = 1
	}
	if message.Category == 5 {
		fPath := message.FilePath + message.FileName + "/"
		flen := strconv.Itoa(len(fPath))
		sPath := newPath + message.FileName + "/"
		// slen := strconv.Itoa(len(sPath))
		sql = "update " + userName + " set parent_path=insert(parent_path, 1, " + flen + ", '" + sPath + "') where mid(parent_path, 1, " + flen + ")='" + fPath + "'"
		fmt.Println(sql)
		stmt, _ = db.DBConn().Prepare(sql)
		_, err = stmt.Exec()
		if err != nil {
			log.Fatalln(err)
		}
	}
	if message.Status == 0 {
		message.FilePath = newPath
	}
}

// FileList : 文件信息
func FileList(userName, parentPath string) []dao.FileMessage {
	sql := "select file_name, category, change_time, file_size from " + userName + " where parent_path=?"
	rows, _ := db.DBConn().Query(sql, parentPath)
	var list []dao.FileMessage
	for rows.Next() {
		var status int8
		var fileSize int
		var name string
		var category int
		var ctime string
		err := rows.Scan(&name, &category, &ctime, &fileSize)
		if err != nil {
			status = 1
		} else {
			list = append(list, dao.FileMessage{
				FileName: name,
				Category: int8(category),
				FilePath: parentPath,
				FileTime: ctime,
				FileSize: fileSize,
				Status:   status,
			})
		}
	}
	return list
}
