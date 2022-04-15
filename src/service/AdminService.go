package service

import (
	"cloudDisk/src/dao"
	"cloudDisk/src/util"
	"cloudDisk/src/util/db"
	"fmt"
	"os"
)

// AdminUploadFile 管理员上传
func AdminUploadFile(fileMD5, fileName, fileSize string) {
	sql := "insert into tbl_file (file_md5,file_name,file_size) values('" + fileMD5 + "',\"" + fileName + "\",'" + fileSize + "')"
	stmt, err := db.DBConn().Prepare(sql)
	if err != nil {
		fmt.Println("UploadFile insert tbl_file:", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
	}
}

// DelFile 管理员删除文件
func DelFile(fileMD5 string) {
	os.Remove(UploadDir + fileMD5 + "_file")
}

// Logout 注销账户
func Logout(userName string) int {
	enforcer := db.Enforcer
	enforcer.LoadPolicy()
	bl, _ := enforcer.RemoveFilteredNamedGroupingPolicy("g", 0, userName)
	if bl {
		sql := "delete from " + userName
		db.DBConn().Exec(sql)
		sql = "drop table " + userName + "_share"
		db.DBConn().Exec(sql)
		sql = "drop table " + userName
		db.DBConn().Exec(sql)
		sql = "delete from tbl_user where user_name='" + userName + "'"
		db.DBConn().Exec(sql)
		sql = "delete from tbl_share where share_name='" + userName + "'"
		db.DBConn().Exec(sql)
		conn := db.Pool.Get()
		defer conn.Close()
		conn.Do("del", userName+"'s photo")
		return 0
	}
	return -1
}

// SetRoleService 设置角色
func SetRoleService(oldPolicy, newPolicy []string) bool {
	enforcer := db.Enforcer
	enforcer.LoadPolicy()
	bl, _ := enforcer.UpdateGroupingPolicy(oldPolicy, newPolicy)
	return bl
}

// DelPolicyService 添加policy
func DelPolicyService(policy dao.Policy) bool {
	enforcer := db.Enforcer
	enforcer.LoadPolicy()
	bl, _ := enforcer.RemovePolicy(policy.Sub, policy.Obj, policy.Act)
	return bl
}

// SetPolicyService 添加policy
func SetPolicyService(oldPolicy, newPolicy []string) bool {
	enforcer := db.Enforcer
	enforcer.LoadPolicy()
	bl, _ := enforcer.UpdatePolicy(oldPolicy, newPolicy)
	return bl
}

// AddPolicyService 添加policy
func AddPolicyService(policy dao.Policy) bool {
	enforcer := db.Enforcer
	enforcer.LoadPolicy()
	bl, _ := enforcer.AddPolicy(policy.Sub, policy.Obj, policy.Act)
	return bl
}

// GetFiles 获取文件列表
func GetFiles() []dao.FilesList {
	sql := "select file_md5, file_size, create_at, update_at from tbl_file order by file_size desc"
	rows, _ := db.DBConn().Query(sql)
	var list []dao.FilesList
	for rows.Next() {
		var temp dao.FilesList
		err := rows.Scan(&temp.FileMD5, &temp.FileSize, &temp.CreateTime, &temp.UpdateTime)
		if err == nil {
			bl, _ := util.IsExist(UploadDir + temp.FileMD5 + "_file")
			if bl {
				list = append(list, temp)
			}
		}
	}
	return list
}

// GetShares 获取分享文件列表
func GetShares() dao.ShareList {
	sql := "select share_addr, share_name, signup_at, share_code, (days-DATEDIFF(now(),signup_at)) from tbl_share"
	rows, _ := db.DBConn().Query(sql)
	var list dao.ShareList
	for rows.Next() {
		var temp dao.ShareFile
		err := rows.Scan(&temp.ShareAddr, &temp.FileName, &temp.SignupAt, &temp.ShareCode, &temp.Days)
		if err == nil {
			list = append(list, temp)
		}
	}
	return list
}

// GetUsers 用户表
func GetUsers() []dao.Users {
	sql := "select user_name, signup_at, last_active from tbl_user"
	rows, _ := db.DBConn().Query(sql)
	var list []dao.Users
	for rows.Next() {
		var temp dao.Users
		err := rows.Scan(&temp.UserName, &temp.SignupAt, &temp.LastActive)
		if err == nil {
			list = append(list, temp)
		}
	}
	enforcer := db.Enforcer
	enforcer.LoadPolicy()
	for k, v := range list {
		list[k].Role = enforcer.GetFilteredNamedGroupingPolicy("g", 0, v.UserName)[0][1]
	}
	return list
}

// GetPolicies 获取所有政策
func GetPolicies() dao.Policies {
	enforcer := db.Enforcer
	enforcer.LoadPolicy()
	list := enforcer.GetPolicy()
	var policies dao.Policies
	for _, i := range list {
		policies = append(policies, dao.Policy{
			Sub: i[0],
			Obj: i[1],
			Act: i[2],
		})
	}
	return policies
}
