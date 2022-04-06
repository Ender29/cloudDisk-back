package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	_ "github.com/go-sql-driver/mysql"

	"github.com/casbin/xorm-adapter/v2"
)

func main() {
	a, err := xormadapter.NewAdapter("mysql", "root:123456@tcp(127.0.0.1:3306)/")
	if err != nil {
		fmt.Println(err)
	}
	e, _ := casbin.NewEnforcer("test/casbinTest/model.conf", a)
	e.LoadPolicy()
	//role := e.GetFilteredNamedGroupingPolicy("g", 0, "ender")[0][1]
	//fmt.Println(e.GetPolicy())
	//sub := "ender"
	//obj := "/file/sharelist"
	//act := "GET"
	//added, err := e.AddPolicy(sub, obj, act)
	//added, err := e.UpdateGroupingPolicy([]string{"ender", "normal"}, []string{"ender", "special"})
	//added, err := e.AddGroupingPolicy(sub, "normal")
	//added, err := e.RemovePolicy(sub, obj, act)
	added, err := e.RemoveFilteredNamedGroupingPolicy("g", 0, "asd")

	fmt.Println(added)
	fmt.Println(err)
	//ok, err := e.Enforce(sub, obj, act)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//if ok == true {
	//	fmt.Println("you pass!")
	//	e.SavePolicy()
	//} else {
	//	fmt.Println("forbid!")
	//}
}
