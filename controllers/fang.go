package controllers

import "github.com/astaxie/beego"

type FangController struct {
	beego.Controller
}

func (this *FangController)Get()  {
	this.TplName="2048.html"
}
