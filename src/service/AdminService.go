package service

import (
	"cloudDisk/src/dao"
	"cloudDisk/src/util/db"
)

// GetFiles 获取文件列表
func GetFiles() []dao.FilesList {
	sql := "select file_md5, file_size, create_at, update_at from tbl_file"
	rows, _ := db.DBConn().Query(sql)
	var list []dao.FilesList
	for rows.Next() {
		var temp dao.FilesList
		err := rows.Scan(&temp.FileMD5, &temp.FileSize, &temp.CreateTime, &temp.UpdateTime)
		if err == nil {
			list = append(list, temp)
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
