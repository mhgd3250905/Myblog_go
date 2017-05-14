package controllers

import (
	"Myblog/models"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
)

type CategoryController struct {
	beego.Controller
}

func (this *CategoryController) Get() {
	op := this.Input().Get("op")
	beego.Warn("当前URL为：", this.Ctx.Request.RequestURI)
	beego.Warn("当前操作为：", op)

	this.Data["IsLogin"] = checkLoginAccount(this.Controller)

	if strings.EqualFold(op, "del") {
		if !checkLoginAccount(this.Controller) {
			beego.Warn("删除category操作，没有登录，请登录！")
			this.Redirect("/login", 302)
			return
		}

		id := this.Input().Get("id")
		if len(id) == 0 {
			return
		}
		err := models.DelCategory(id)
		if err != nil {
			beego.Error(err)
		}
		this.Redirect("/category", 301)
		return
	}
	this.TplName = "category.html"
	this.Data["IsCategory"] = true
	var err error
	this.Data["Categories"], err = models.GetAllCategories()
	beego.Info(this.Data["Categories"])
	if err != nil {
		beego.Error(err)
	}

}

func (this *CategoryController) Post() {
	op := this.Input().Get("op")
	beego.Warn("当前URL为：", this.Ctx.Request.RequestURI)
	beego.Warn("当前操作为：", op)

	if strings.EqualFold(op, "add") {
		if !checkLoginAccount(this.Controller) {
			beego.Warn("添加category操作，没有登录，请登录！")
			this.Redirect("/login", 302)
			return
		}

		name := this.Input().Get("name")
		fmt.Printf("接收到name为：%s\n", name)
		if len(name) == 0 {
			return
		}
		err := models.AddCategory(name)
		if err != nil {
			beego.Error("添加文章到数据库错误：%s", err)
		}
		beego.Info("添加文章到数据库成功！：%s", err)
		this.Redirect("/category", 301)
		return
	}
}
