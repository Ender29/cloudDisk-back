package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	_ "github.com/go-sql-driver/mysql"

	"github.com/casbin/xorm-adapter/v2"
)

func main() {
	a, err := xormadapter.NewAdapter("mysql", "root:123456@tcp(124.223.78.104:3306)/")
	if err != nil {
		fmt.Println(err)
	}
	e, _ := casbin.NewEnforcer("test/casbinTest/model.conf", a)
	e.LoadPolicy()
	sub := "zhangsan"
	obj := "data2"
	act := "read"
	//added, err := e.AddPolicy(sub, obj, act)
	//added, err := e.AddGroupingPolicy(sub, "admin")
	//fmt.Println(added)
	//fmt.Println(err)
	ok, err := e.Enforce(sub, obj, act)
	if err != nil {
		fmt.Println(err)
	}
	if ok == true {
		fmt.Println("you pass!")
		e.SavePolicy()
	} else {
		fmt.Println("forbid!")
	}

}
