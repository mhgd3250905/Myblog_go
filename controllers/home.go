package controllers

import (
	"Myblog/models"
	"fmt"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	//获取Session中保存的数据
	uname := this.GetSession("uname")
	pwd := this.GetSession("pwd")

	this.Data["IsHome"] = true
	this.TplName = "home.html"
	this.Data["IsLogin"] = checkAccount(uname, pwd)
	fmt.Printf("isLogin：%t", this.Data["IsLogin"])
	cate := this.Input().Get("cate")
	label := this.Input().Get("label")

	topics, err := models.GetAllTopics(cate, label, true)
	if err != nil {
		beego.Error(err)
	} else {
		this.Data["Topics"] = topics
	}
	categories, err := models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}

	this.Data["Categories"] = categories
}
