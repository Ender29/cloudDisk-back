package service

import (
	util "cloudDisk/src/util"
	"cloudDisk/src/util/db"
	"cloudDisk/src/vo"
	"github.com/garyburd/redigo/redis"
	"io"
	"log"
	"mime/multipart"
	"time"
)

type LoginMessage vo.LoginMessage
type RegisterMessage vo.RegisterMessage

// UploadPhotoService 存储头像
func UploadPhotoService(userName string, img *multipart.FileHeader) string {
	file, _ := img.Open()
	defer file.Close()
	b, _ := io.ReadAll(file)
	imgBase64 := util.ToBase64(b)
	conn := db.Pool.Get()
	defer conn.Close()
	imgBase64 = "data:image/png;base64," + imgBase64
	conn.Do("set", userName+"'s photo", imgBase64)
	return imgBase64
}

// Login 登录
func (lm *LoginMessage) Login(userPwd string) {
	var times string
	var status int8 = 0
	var fileSize int64
	sql := "select sum(file_size) sum from " + lm.UserName
	rows, err := db.DBConn().Query(sql)
	if err != nil {
		status = 5
	} else {
		for rows.Next() {
			var size int64
			err = rows.Scan(&size)
			if err != nil {
				size = 0
			}
			fileSize = size
		}
	}

	sql = "select user_name, user_pwd, last_active from tbl_user where user_name=? and user_pwd=?"
	rows, err = db.DBConn().Query(sql, lm.UserName, userPwd)
	tmp := true
	for rows.Next() {
		tmp = false
		var username string
		var lastime string
		var userpwd string
		// 获取查询结果
		err = rows.Scan(&username, &userpwd, &lastime)
		if err != nil {
			status = 2
		}
		times = lastime
		timeNow := time.Now().Format("2006-01-02 15:04:05")

		sql = "update tbl_user set last_active=? where user_name=?"
		stmt, _ := db.DBConn().Prepare(sql)
		result, err := stmt.Exec(timeNow, lm.UserName)
		if err != nil {
			status = 3
		}
		affect, _ := result.RowsAffected()

		if affect == 0 {
			status = 4
		}
	}
	if tmp && status != 5 {
		status = 6
		fileSize = 0
	}
	lm.AccessToken, _ = util.GenerateToken(lm.UserName, userPwd, time.Minute*5)
	conn := db.Pool.Get()
	defer conn.Close()
	// redis 绑定token,设置过期时间
	refreshToken, _ := util.GenerateToken(lm.UserName, userPwd, time.Hour*24)
	_, err = conn.Do("set", lm.AccessToken, refreshToken, "EX", 3600*24)
	lm.Status = status
	lm.FileSize = fileSize
	lm.LatestTime = times
	lm.HeadPhoto, _ = redis.String(conn.Do("get", lm.UserName+"'s photo"))
}

func (rm *RegisterMessage) Register(userPwd string) {
	var status int8 = 0
	realPwd, _ := util.EnPwdCode([]byte(userPwd))
	// 插入用户信息
	stmt, err := db.DBConn().Prepare("INSERT tbl_user set user_name=?,user_pwd=?")
	_, err = stmt.Exec(rm.UserName, realPwd)
	if err != nil {
		status = 1
	}
	var table = `CREATE TABLE ` + rm.UserName + ` (
		parent_path varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  		file_name varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
		file_md5 varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
		file_size bigint DEFAULT '0',
		category int DEFAULT '0',
		upload_at datetime DEFAULT CURRENT_TIMESTAMP,
		change_time datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (parent_path, file_name) USING BTREE
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci `
	if status == 0 {
		// 创建用户文件表
		stmt, err := db.DBConn().Prepare(table)
		if err != nil {
			db.DBConn().Prepare("delete from tbl_user where user_name=?")
			_, _ = stmt.Exec(rm.UserName)
			status = 2
		} else {
			_, err = stmt.Exec()
			if err != nil {
				status = 3
			}
			table = `CREATE TABLE ` + rm.UserName + `_share (
				parent_path varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
				file_name varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
				share_addr varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
				FOREIGN KEY(parent_path,file_name) REFERENCES ` + rm.UserName + `(parent_path,file_name) ON UPDATE CASCADE ON DELETE CASCADE,
				FOREIGN key (share_addr) REFERENCES tbl_share(share_addr) ON UPDATE CASCADE ON DELETE CASCADE, 
				UNIQUE INDEX id_share_addr (share_addr) USING BTREE
			) ENGINE = InnoDB AUTO_INCREMENT = 128 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;`
			stmt, _ = db.DBConn().Prepare(table)
			_, err = stmt.Exec()
			if err != nil {
				status = 4
			}
		}
	}
	// 返回状态码
	rm.Status = status
}

// Logout : 注销账户
func Logout(userName, userToken, status *string) {
	sql := "delete from tbl_user where user_name=? and user_token=?"
	stmt, _ := db.DBConn().Prepare(sql)
	_, err := stmt.Exec(userName, userToken)
	if err != nil {
		*status = "1"
	}
	sql = "drop table " + *userName
	stmt, _ = db.DBConn().Prepare(sql)
	_, err = stmt.Exec()
	if err != nil {
		*status = "2"
	}
	sql = "drop table " + *userName + "_share"
	stmt, _ = db.DBConn().Prepare(sql)
	_, err = stmt.Exec()
	if err != nil {
		*status = "3"
	}
}

// ChangePwd : 修改密码
func ChangePwd(userName, userPwd, newPwd string) string {
	sql := "update tbl_user set user_pwd=?,user_token='0' where user_name=? and user_pwd=?;"
	stmt, _ := db.DBConn().Prepare(sql)
	_, err := stmt.Exec(newPwd, userName, userPwd)
	if err != nil {
		log.Fatalln(err)
		return "1"
	}
	return "0"
}
