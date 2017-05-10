package controllers

import (
	"github.com/astaxie/beego"
	"Myblog/models"
)

type ReplyController struct {
	beego.Controller
}

func (this *ReplyController) Add()  {
	tid:=this.Input().Get("tid")
	nickName:=this.Input().Get("nickname")
	content:=this.Input().Get("content")
	err:=models.AddReply(tid,nickName,content)
	if err !=nil{
		beego.Error(err)
	}
	this.Redirect("/topic/view/"+tid,302)
}

func (this *ReplyController) Delete()  {

	tid:=this.Input().Get("tid")
	rid:=this.Input().Get("rid")
	beego.Info("tid: "+tid+" rid"+rid+"\n")
	err:=models.DeleteReply(rid)
	if err!=nil {
		beego.Info("删除出现问题！")
		beego.Error(err)
	}
	beego.Info("此刻Tid为："+tid)
	this.Redirect("/topic/view/"+tid,302)
}