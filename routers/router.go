package routers

import (
	"demoIM/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/ws", &controllers.WsController{}, "get:Connect")
}