package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"log"
	"fmt"
)

type LoginController struct {
	beego.Controller
}

func (this*LoginController) Get() {
	isExit := this.Input().Get("exit")=="true"
	if isExit {
		log.Println(isExit)
		this.Ctx.SetCookie("uname","",-1,"/")
		this.Ctx.SetCookie("pwd", "",-1, "/")
		this.Redirect("/",302)
		return
	}
	this.TplName = "login.html"
}

func (this*LoginController) Post() {
	uname := this.Input().Get("uname")
	pwd := this.Input().Get("pwd")
	autoLogin := this.Input().Get("autoLogin") == "on"
	beego.Info("Post:uname"+uname+" pwd"+pwd+" ")
	beego.Info("是否勾选:",autoLogin)
	beego.Info("验证账户密码:",beego.AppConfig.String("uname") == uname &&
		beego.AppConfig.String("pwd") == pwd)

	beego.Info("测试URL:",this.Ctx.Request.URL)


	if beego.AppConfig.String("uname") == uname &&
		beego.AppConfig.String("pwd") == pwd {
		maxAge := 3600//保存一小时
		if autoLogin {
			maxAge = 1<<31 - 1//保存N小时
		}
		this.Ctx.SetCookie("uname", uname, maxAge,"/")
		this.Ctx.SetCookie("pwd", pwd, maxAge, "/")
		this.Redirect("/", 302)
	}
}

//检测账号密码是否正确
func checkAccount(ctx *context.Context) bool {
	uname:=ctx.GetCookie("uname")
	pwd :=ctx.GetCookie("pwd")
	beego.Info("获取到Cookie：uname: ",uname," pwd: ",pwd)
	fmt.Printf("uname: %s,pwd: %s",uname,pwd)
	return (beego.AppConfig.String("uname") == uname && beego.AppConfig.String("pwd") == pwd)
}

