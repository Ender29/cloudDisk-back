package db

import (
	"cloudDisk/src/dao"
	"database/sql"
	"log"

	// mysql srvier
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User dao.User

func init() {
	db, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/netWork?charset=utf8")
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
