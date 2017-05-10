package main

import (
	"Myblog/models"
	_ "Myblog/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"Myblog/controllers"
	"os"
)


func init() {
	models.RegisterDB()
}

func main() {
	orm.Debug = true
	//参数1：同步db 参数2：每次是否强制重建 参数3：是否打印log
	orm.RunSyncdb("default",false,true)
	//注册路由
	beego.Router("/",&controllers.MainController{})
	beego.Router("/login",&controllers.LoginController{})
	beego.Router("/category",&controllers.CategoryController{})
	beego.Router("/topic",&controllers.TopicController{})
	//使用自动路由
	beego.AutoRouter(&controllers.TopicController{})
	beego.Router("/reply",&controllers.ReplyController{})
	beego.Router("/reply/add",&controllers.ReplyController{},"post:Add")
	beego.Router("/reply/delete",&controllers.ReplyController{},"get:Delete")

	//创建附件目录
	os.Mkdir("attachment",os.ModePerm)
	//作为静态文件处理
	//beego.SetStaticPath("/attachment","attachment")
	//作为一个单独的控制器来处理
	beego.Router("/attachment/:all",&controllers.AttachController{})

	beego.Run()
}
