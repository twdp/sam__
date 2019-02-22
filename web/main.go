package main

import (
	"github.com/astaxie/beego/orm"
	"tianwei.pro/business/controller"
	"tianwei.pro/micro/conf"
	"tianwei.pro/micro/di/single"
	"tianwei.pro/micro/rpc/client"
	"tianwei.pro/micro/rpc/server"
	"tianwei.pro/sam/agent"
	_ "tianwei.pro/sam/core/init"
	_ "tianwei.pro/sam/web/routers"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	// set default database
	orm.RegisterDataBase("default", "mysql", "root:anywhere@tcp(127.0.0.1:3306)/sam?charset=utf8&loc=Asia%2FShanghai", 30)

	// create table
	orm.RunSyncdb("default", false, true)

	if conf.Conf.DefaultString("runmode", "prod") != "prod" {
		orm.Debug = true
	}
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	defer server.Close()
	defer client.Close()

	single.Resolve()

	beego.ErrorController(&controller.RestfulErrorHandler{})

	beego.BConfig.RecoverFunc = controller.RestfulErrorHandle
	beego.InsertFilter("/*", beego.BeforeRouter, agent.SamFilter)
	beego.Run()
}
