package service

import (
	"cloudDisk/src/dao"
	mysql "cloudDisk/src/util"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"
)


// CreateCatalog : 创建文件夹
func CreateCatalog(userName, fileName, theTime, path string) (id int, status int8) {
	status = 0
	sql := "INSERT INTO " + userName + "(parent_path,file_name,category,upload_at,change_time) SELECT '" + path + "', '" + fileName + "', 5,'" + theTime + "','" + theTime + "' from DUAL where NOT exists(select id from " + userName + " where parent_path='" + path + "' and file_name = '" + fileName + "');"
	fmt.Println(sql)
	stmt, _ := mysql.DBConn().Prepare(sql)
	_, err := stmt.Exec()
	if err != nil {
		status = 2
	}
	// 获取ID值
	sql = "select id from " + userName + " where parent_path=? and file_name=?"
	rows, _ := mysql.DBConn().Query(sql, path, fileName)
	for rows.Next() {
		var fid int
		err := rows.Scan(&fid)
		if err != nil {
			status = 3
		} else {
			id = fid
		}
	}
	return id, status
}

// DeleteFile : 删除文件
func DeleteFile(message *dao.FileMessage, userName string) {
	sql := "select parent_path, file_name, file_size, category from " + userName + " where id=?"
	rows, err := mysql.DBConn().Query(sql, message.FileID)
	if err != nil {
		message.Status = 1
	}
	for rows.Next() {
		// 查询文件信息
		err := rows.Scan(&message.FilePath, &message.FileName, &message.FileSize, &message.Category)
		if err != nil {
			message.Status = 2
		}
	}
	if message.Status == 0 {
		sql = "delete from " + userName + " where id=?"
		stmt, _ := mysql.DBConn().Prepare(sql)
		_, err = stmt.Exec(message.FileID)
		if err != nil {
			message.Status = 3
		}
		if message.Category == 5 {
			str := message.FilePath + message.FileName
			if !strings.HasSuffix(str, "/") {
				str += "/"
			}
			length := strconv.Itoa(utf8.RuneCountInString(str))
			sql = "delete from " + userName + " where mid(parent_path, 1, " + length + ")='" + str + "'"
			fmt.Println(sql)
			stmt, _ = mysql.DBConn().Prepare(sql)
			_, err = stmt.Exec()
			if err != nil {
				message.Status = 4
			}
		}
	}
}

// RenameFile : 文件重命名
func RenameFile(userName string, message *dao.FileMessage, newName string) {
	sql := "update " + userName + " set file_name=? where id=?"
	stmt, _ := mysql.DBConn().Prepare(sql)
	_, err := stmt.Exec(newName, message.FileID)
	if err != nil {
		message.Status = 1
	}
	if message.Category == 5 {
		message.Isdir = 1
		length1 := strconv.Itoa(utf8.RuneCountInString(message.FilePath) + 1)
		length2 := strconv.Itoa(utf8.RuneCountInString(message.FilePath))
		length3 := strconv.Itoa(utf8.RuneCountInString(message.FileName) + utf8.RuneCountInString(message.FilePath) + 1)
		sql := "update " + userName + " set parent_path=insert(parent_path, " + length1 + ", " + length2 + ",'/" + newName + "') where mid(parent_path, 1, " + length3 + ")='" + message.FilePath + "/" + message.FileName + "'"

		_, err := mysql.DBConn().Exec(sql)
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
	sql := "update " + userName + " set parent_path=? where id=?"
	stmt, _ := mysql.DBConn().Prepare(sql)
	_, err := stmt.Exec(newPath, message.FileID)
	if err != nil {
		message.Status = 1
	}
	if message.Category == 5 {
		message.Isdir = 1
		fPath := message.FilePath + message.FileName + "/"
		flen := strconv.Itoa(len(fPath))
		sPath := newPath + message.FileName + "/"
		// slen := strconv.Itoa(len(sPath))
		sql := "update " + userName + " set parent_path=insert(parent_path, 1, " + flen + ", '" + sPath + "') where mid(parent_path, 1, " + flen + ")='" + fPath + "'"
		fmt.Println(sql)
		stmt, _ = mysql.DBConn().Prepare(sql)
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
	sql := "select id, file_name, category, change_time, file_size from " + userName + " where parent_path=?"
	rows, _ := mysql.DBConn().Query(sql, parentPath)
	var list []dao.FileMessage
	for rows.Next() {
		var status int8
		var id int
		var fileSize int
		var name string
		var category int
		var ctime string
		var isDir int8
		err := rows.Scan(&id, &name, &category, &ctime, &fileSize)
		if err != nil {
			status = 1
		} else {
			if category == 5 {
				isDir = 1
			}
			list = append(list, dao.FileMessage{
				FileID:   id,
				FileName: name,
				Isdir:    isDir,
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

