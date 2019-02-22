package main

import (
	"tianwei.pro/business/controller"
	_ "tianwei.pro/sam/routers"

	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.BConfig.RecoverFunc = controller.RestfulErrorHandle
	beego.Run()
}
