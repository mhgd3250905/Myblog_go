package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"Myblog/models"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Data["IsHome"] = true
	this.TplName = "home.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	fmt.Printf("isLogin：%t",this.Data["IsLogin"])
	cate:=this.Input().Get("cate")
	label:=this.Input().Get("label")

	topics,err:=models.GetAllTopics(cate,label,true)
	if err != nil {
		beego.Error(err)
	}else {
		this.Data["Topics"]=topics
	}
	categories,err:=models.GetAllCategories()
	if err!=nil{
		beego.Error(err)
	}

	this.Data["Categories"]=categories
}
