package controllers

import (
	"github.com/astaxie/beego"
	"Myblog/models"
	"strings"
	"path"
)

type TopicController struct {
	beego.Controller
}

func (this *TopicController)Get()  {
	this.Data["IsTopic"]=true
	this.TplName="topic.html"
	uname:=this.GetSession("uname")
	pwd:=this.GetSession("pwd")
	this.Data["IsLogin"] = checkAccount(uname, pwd)
	topics,err:=models.GetAllTopics("","",false)
	beego.Warn(len(topics))
	if err != nil {
		beego.Error(err)
	}else {
		this.Data["Topics"]=topics
	}

}

func (this *TopicController) Post()  {
	uname:=this.GetSession("uname")
	pwd:=this.GetSession("pwd")
	this.Data["IsLogin"] = checkAccount(uname, pwd)
	//判断是否登录
	if !checkAccount(uname, pwd) {
		beego.Warn("添加文章操作，没有登录，请登录！")
		this.Redirect("/login",302)
		return
	}

	//将form提交数据映射到struct
	testTopic:=models.Topic{}
	if err:=this.ParseForm(&testTopic);err!=nil{
		beego.Error(err)
	}
	beego.Info("获取到Topic-----------------------------")
	beego.Info(testTopic)
	beego.Info("获取到Topic-----------------------------")


	//获取附件
	_,fh,err:=this.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}

	var attachment string
	if fh != nil {
		//保存附件
		attachment=fh.Filename
		beego.Info("\n接收到文件: ",attachment,"\n")
		err=this.SaveToFile("attachment",path.Join("attachment",attachment))
		//filename:tmp.go--->attachment/tmp.go
		if err != nil {
			beego.Error(err)
		}

	}

	//设置附件
	testTopic.Attachment=attachment

	//判断是否上传tid：如果有，就是新增文章 否则就是修改文章
	if len(this.Input().Get("tid"))==0 {
		err=models.AddTopic(&testTopic)
	}else {
		title:=this.Input().Get("title")
		category:=this.Input().Get("category")
		label:=this.Input().Get("label")
		content:=this.Input().Get("content")
		tid:=this.Input().Get("tid")
		err=models.ModifyTopic(tid,title,category,label,content,attachment)
	}

	if err!=nil{
		beego.Error(err)
	}
	this.Redirect("/topic",302)
}


//增加文章
func (this *TopicController)Add()  {
	categories,err:=models.GetAllCategories()
	if err!=nil {
		beego.Error(err)
	}
	this.Data["Categories"]=categories
	this.TplName="topic_add.html"
}

//文章详情
func (this *TopicController)View()  {
	this.TplName="topic_view.html"
	topic,err:=models.GetTopic(this.Ctx.Input.Param("0"))
	if err != nil{
		beego.Error(err)
		this.Redirect("/",302)
		return
	}
	this.Data["Topic"]=topic
	this.Data["Tid"]=this.Ctx.Input.Param("0")
	var splitLabel []string
	labels:=make([]string,0)
	splitLabel=strings.Split(topic.Label," ")
	for i,value:=range splitLabel{
		if len(value)>0 {
			labels=append(labels,splitLabel[i])
		}
	}
	this.Data["Label"]=labels

	//获取回复内容
	replies,err:=models.GetAllReplies(this.Ctx.Input.Param("0"))
	if err != nil {
		beego.Error(err)
		return
	}
	this.Data["Replies"]=replies
	uname:=this.GetSession("uname")
	pwd:=this.GetSession("pwd")
	this.Data["IsLogin"] = checkAccount(uname, pwd)
}

func (this *TopicController)Modify(){
	this.TplName="topic_modify.html"
	tid:=this.Input().Get("tid")
	topic,err:=models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		this.Redirect("/",302)
		return
	}

	this.Data["Topic"]=topic
	this.Data["Tid"]=tid
}

func (this *TopicController)Delete()  {
	uname:=this.GetSession("uname")
	pwd:=this.GetSession("pwd")
	if !checkAccount(uname,pwd) {
		beego.Warn("删除文章操作，没有登录，请登录！")
		this.Redirect("/login",302)
		return
	}
	categroy:=this.Input().Get("category")

	err:=models.DelTopic(categroy,this.Ctx.Input.Param("0"))
	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/topic",302)
}


