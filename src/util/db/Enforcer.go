package db

import (
	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	"log"
)

var Enforcer *casbin.Enforcer

//初始化 enforcer
func init() {
	adapter, err := xormadapter.NewAdapter("mysql", "root:123456@tcp(127.0.0.1:3306)/")
	if err != nil {
		log.Fatalf("init adapter fail: %s\n", err)
	}
	Enforcer, err = casbin.NewEnforcer("src/config/model.conf", adapter)
	if err != nil {
		log.Fatalf("init enforcer fail: %s\n", err)
	}
	//Enforcer.EnableLog(true)
}
