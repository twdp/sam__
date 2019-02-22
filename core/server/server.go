package main

import (
	"flag"
	"github.com/astaxie/beego/orm"
	"tianwei.pro/micro/conf"
	"tianwei.pro/micro/rpc/server"
	_ "tianwei.pro/micro/rpc/server"
	_ "tianwei.pro/sam/core/init"
)

func init() {
	// set default database
	orm.RegisterDataBase("default", "mysql", conf.Conf.String("db.url"), 30)

	// create table
	orm.RunSyncdb("default", false, true)

	if conf.Conf.DefaultString("runmode", "prod") != "prod" {
		orm.Debug = true
	}
}

func main() {
	flag.Parse()

	server.Server.Start()

	defer server.Close()
}
