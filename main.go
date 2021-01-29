package main

import (
	"github.com/astaxie/beego"
	_ "demoIM/routers"
)

func main() {
	beego.Info(beego.BConfig.AppName, beego.AppConfig.String("appversion"))
	beego.Run()
}

func init() {
	
}


