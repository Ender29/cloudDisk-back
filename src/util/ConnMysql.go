package util

import (
	"cloudDisk/src/dao"
	"database/sql"
	"fmt"
	"log"

	// mysql srvier
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User dao.User

func init() {
	db, _ = sql.Open("mysql", "root:123456@tcp(124.223.78.104:3306)/netWork2?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		log.Fatalln("建立连接失败:", err)
	}
}

// DBConn : 返回数据库连接
func DBConn() *sql.DB {
	return db
}

func GetUser(userName string) *User {
	//fmt.Println("GetUser()-->userName:", userName)
	stmt, err := DBConn().Prepare("select * from tbl_user where user_name = ?")
	defer stmt.Close()
	if err != nil {
		fmt.Println("sql预编译失败")
		return nil
	}
	row := stmt.QueryRow(userName)
	user := User{}
	err = row.Scan(&user.ID, &user.UserName, &user.UserPwd, &user.Email, &user.Phone, &user.EmailValidated,
		&user.PhoneValidated, &user.SignupAt, &user.LastActive, &user.Profile, &user.Status, &user.UserToken)
	if err == sql.ErrNoRows {
		// 没有查询到任何一条数据
		return nil
	} else if err != nil {
		fmt.Println("获取唯一表数据失败--", err.Error())
		return nil
	}
	return &user
}
