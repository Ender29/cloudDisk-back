package service

import (
	util "cloudDisk/src/util"
	"cloudDisk/src/vo"
	"log"
	"time"
)

type LoginMessage vo.LoginMessage
type RegisterMessage vo.RegisterMessage

func (lm *LoginMessage)Login(userPwd string)  {
	token := ""
	var times string
	var status int8 = 0
	var fileSize int64
	sql := "select sum(file_size) sum from " + lm.UserName
	rows, err := util.DBConn().Query(sql)
	if err!=nil {
		status = 5
	}else {
		for rows.Next() {
			var size int64
			err := rows.Scan(&size)
			if err != nil {
				size = 0
			}
			fileSize = size
		}
	}

	sql = "select user_name, user_pwd, last_active, user_token from tbl_user where user_name=? and user_pwd=?"
	rows, err = util.DBConn().Query(sql, lm.UserName, userPwd)
	tmp := true
	for rows.Next() {
		tmp = false
		var username string
		var lastime string
		var userpwd string
		var userToken string
		// 获取查询结果
		err := rows.Scan(&username, &userpwd, &lastime, &userToken)
		if err != nil {
			status = 2
		}
		if userToken == "0" {
			userToken, _ = util.GenerateToken(lm.UserName, userPwd)
		}
		times = lastime
		timeNow := time.Now().Format("2006-01-02 15:04:05")

		sql = "update tbl_user set last_active=?,user_token=? where user_name=?"
		stmt, _ := util.DBConn().Prepare(sql)
		result, err := stmt.Exec(timeNow, userToken, lm.UserName)
		if err != nil {
			status = 3
		}
		affect, _ := result.RowsAffected()

		if affect == 0 {
			status = 4
		}
		token = userToken
	}
	if tmp && status!=5 {
		status = 6
		fileSize = 0
	}
	lm.UserToken = token
	lm.Status = status
	lm.FileSize = fileSize
	lm.LatestTime = times
}

func (rm *RegisterMessage)Register(userPwd string)  {
	var status int8 = 0
	realPwd, _ := util.EnPwdCode([]byte(userPwd))
	// 插入用户信息
	stmt, err := util.DBConn().Prepare("INSERT tbl_user set user_name=?,user_pwd=?")
	_, err = stmt.Exec(rm.UserName, realPwd)
	if err != nil {
		status = 1
	}
	var table string = `CREATE TABLE ` + rm.UserName + ` (
		id int NOT NULL AUTO_INCREMENT,
		parent_path varchar(1024) NOT NULL,
		file_name varchar(1024) NOT NULL,
		file_sha1 varchar(64) NOT NULL DEFAULT '',
		file_size bigint DEFAULT '0',
		category int DEFAULT '0',
		upload_at datetime DEFAULT CURRENT_TIMESTAMP,
		change_time datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		ext1 int DEFAULT 0,
		ext2 text,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci `
	if status == 0 {
		// 创建用户文件表
		stmt, err := util.DBConn().Prepare(table)
		if err != nil {
			util.DBConn().Prepare("delete from tbl_user where user_name=?")
			_, _ = stmt.Exec(rm.UserName)
			status = 2
		} else {
			_, err = stmt.Exec()
			if err != nil {
				status = 3
			}
		}
	}
	// 返回状态码
	rm.Status = status
}

// Logout : 注销账户
func Logout(userName, userToken, status *string) {
	sql := "delete from tbl_user where user_name=? and user_token=?"
	stmt, _ := util.DBConn().Prepare(sql)
	_, err := stmt.Exec(userName, userToken)
	if err != nil {
		*status = "1"
	}
	sql = "drop table " + *userName
	stmt, _ = util.DBConn().Prepare(sql)
	_, err = stmt.Exec()
	if err != nil {
		*status = "2"
	}
}

// ChangePwd : 修改密码
func ChangePwd(userName, userPwd, newPwd string) string {
	sql := "update tbl_user set user_pwd=?,user_token='0' where user_name=? and user_pwd=?;"
	stmt, _ := util.DBConn().Prepare(sql)
	_, err := stmt.Exec(newPwd, userName, userPwd)
	if err != nil {
		log.Fatalln(err)
		return "1"
	}
	return "0"
}
