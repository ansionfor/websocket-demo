package controllers

import (
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

type IndexController struct {
	BaseController
}

func (this * IndexController) Get() {
	this.Ctx.WriteString("index")
}