package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"log"
)

type LoginController struct {
	beego.Controller
}

func (this*LoginController) Get() {
	isExit := this.Input().Get("exit") == "true"
	if isExit {
		beego.Info("----------")
		beego.Info("删除Session")
		beego.Info("----------\n")

		log.Println(isExit)
		this.DelSession("uname")
		this.DelSession("pwd")
		this.Redirect("/", 302)
		return
	}
	this.TplName = "login.html"
}

func (this*LoginController) Post() {
	uname := this.Input().Get("uname")
	pwd := this.Input().Get("pwd")
	autoLogin := this.Input().Get("autoLogin") == "on"
	beego.Info("Post:uname" + uname + " pwd" + pwd + " ")
	beego.Info("是否勾选:", autoLogin)
	beego.Info("验证账户密码:", beego.AppConfig.String("uname") == uname &&
		beego.AppConfig.String("pwd") == pwd)

	beego.Info("测试URL:", this.Ctx.Request.URL)

	if beego.AppConfig.String("uname") == uname &&
		beego.AppConfig.String("pwd") == pwd {
		if autoLogin { //保存N小时
			beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 1<<31 - 1
		}

		this.SetSession("uname", uname)
		this.SetSession("pwd", pwd)

		this.Redirect("/", 302)
	}
}

////检测账号密码是否正确
//func checkAccount(ctx *context.Context) bool {
//	uname:=ctx.GetCookie("uname")
//	pwd :=ctx.GetCookie("pwd")
//	beego.Info("获取到Cookie：uname: ",uname," pwd: ",pwd)
//	fmt.Printf("uname: %s,pwd: %s",uname,pwd)
//	return (beego.AppConfig.String("uname") == uname && beego.AppConfig.String("pwd") == pwd)
//}

//检测账号密码是否正确

/**
使用传入的数据来判断用户是否登录
*/
func checkLoginAccount(ctx beego.Controller) bool {
	uname := ctx.GetSession("uname")
	pwd := ctx.GetSession("pwd")
	beego.Info("获取到Sessionn内容：uname: ", uname, " pwd: ", pwd)
	fmt.Printf("uname: %s,pwd: %s", uname, pwd)
	return (beego.AppConfig.String("uname") == uname && beego.AppConfig.String("pwd") == pwd)
}
