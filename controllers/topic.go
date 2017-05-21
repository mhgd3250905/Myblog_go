package controllers

import (
	"Myblog/models"
	"encoding/base64"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
	"fmt"
	"github.com/russross/blackfriday"
	"github.com/microcosm-cc/bluemonday"
)

type TopicController struct {
	beego.Controller
}

func (this *TopicController) Get() {
	this.Data["IsTopic"] = true
	this.TplName = "topic.html"
	this.Data["IsLogin"] = checkLoginAccount(this.Controller)
	topics, err := models.GetAllTopics("", "", false)
	beego.Warn(len(topics))
	if err != nil {
		beego.Error(err)
	} else {
		this.Data["Topics"] = topics
	}

}

func (this *TopicController) Post() {
	this.Data["IsLogin"] = checkLoginAccount(this.Controller)
	//判断是否登录
	if !checkLoginAccount(this.Controller) {
		beego.Warn("添加文章操作，没有登录，请登录！")
		this.Redirect("/login", 302)
		return
	}

	//将form提交数据映射到struct
	testTopic := models.Topic{}
	if err := this.ParseForm(&testTopic); err != nil {
		beego.Error(err)
	}

	beego.Info("获取到Topic-----------------------------")
	beego.Info(testTopic)
	beego.Info("获取到Topic-----------------------------")


	////保存图片
	//topicContent,err:=saveImgAndReplace(&testTopic)
	//if err != nil {
	//	beego.Warn("content转换失败！")
	//}
	//
	//testTopic.Content=topicContent


	//获取附件
	_, fh, err := this.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}

	var attachment string
	if fh != nil {
		//保存附件
		attachment = fh.Filename
		beego.Info("\n接收到文件: ", attachment, "\n")
		err = this.SaveToFile("attachment", path.Join("attachment", attachment))
		//filename:tmp.go--->attachment/tmp.go
		if err != nil {
			beego.Error(err)
		}

	}

	//设置附件
	testTopic.Attachment = attachment

	//判断是否上传tid：如果有，就是新增文章 否则就是修改文章
	if len(this.Input().Get("tid")) == 0 {
		err = models.AddTopic(&testTopic)
	} else {
		title := this.Input().Get("title")
		category := this.Input().Get("category")
		label := this.Input().Get("label")
		content := this.Input().Get("content")
		tid := this.Input().Get("tid")
		err = models.ModifyTopic(tid, title, category, label, content, attachment)
	}

	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/topic", 302)
}

//增加文章
func (this *TopicController) Add() {

	categories, err := models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
	this.Data["Categories"] = categories
	this.TplName = "topic_add.html"
	this.Data["IsLogin"] = checkLoginAccount(this.Controller)
}

//文章详情
func (this *TopicController) View() {
	this.TplName = "topic_view.html"
	topic, err := models.GetTopic(this.Ctx.Input.Param("0"))
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}

	unsafe := blackfriday.MarkdownCommon([]byte(topic.Content))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	topic.Content=string(html)
	this.Data["Topic"] = topic
	this.Data["Tid"] = this.Ctx.Input.Param("0")
	var splitLabel []string
	labels := make([]string, 0)
	splitLabel = strings.Split(topic.Label, " ")
	for i, value := range splitLabel {
		if len(value) > 0 {
			labels = append(labels, splitLabel[i])
		}
	}
	this.Data["Label"] = labels

	//获取回复内容
	replies, err := models.GetAllReplies(this.Ctx.Input.Param("0"))
	if err != nil {
		beego.Error(err)
		return
	}
	this.Data["Replies"] = replies

	this.Data["IsLogin"] = checkLoginAccount(this.Controller)
}

func (this *TopicController) Modify() {
	this.TplName = "topic_modify.html"
	tid := this.Input().Get("tid")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}

	this.Data["Topic"] = topic
	this.Data["Tid"] = tid
	this.Data["IsLogin"] = checkLoginAccount(this.Controller)
}

func (this *TopicController) Delete() {
	if !checkLoginAccount(this.Controller) {
		beego.Warn("删除文章操作，没有登录，请登录！")
		this.Redirect("/login", 302)
		return
	}
	categroy := this.Input().Get("category")

	err := models.DelTopic(categroy, this.Ctx.Input.Param("0"))
	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/topic", 302)
}



func saveImgAndReplace(topic *models.Topic) (string,error) {
	//解析html中的<img> src属性
	reader := strings.NewReader(topic.Content)
	doc, err := goquery.NewDocumentFromReader(reader)
	imgDom:= doc.Find("img").Eq(0)
	imageBase64, exist :=imgDom.Attr("src")
	//如果属性内容不存在
	if !exist {
		beego.Info("<img>的src属性不存在！")
		return "",err
	}
	//去除一般的前缀描述，获取base64字符串
	c := imageBase64[22:]
	imageStr := string(c)

	//解码base64 头png
	imgBuf, err := base64.StdEncoding.DecodeString(imageStr) //成图片文件并把文件写入到buffer
	if err != nil {
		beego.Error("解析base64失败")
		return "",err
	}

	exist, err = PathExists("./topic_img")
	if err != nil {
		return "",err
	}
	if !exist {
		err = os.Mkdir("./topic_img", 0777)
		if err != nil {
			beego.Info("创建topic_img目录失败！")

		}
	}


	timeStr:=time.Now().Format("20060102150405")
	imageName:=fmt.Sprintf("./static/topic_img/%d_%s.png",topic.Id,timeStr) //id不对！！！！！！！！！！！！！
	beego.Info("图片文件名字为："+imageName)

	err = ioutil.WriteFile(imageName, imgBuf, 0666) //buffer输出到jpg文件中（不做处理，直接写到文件）
	if err != nil {
		beego.Error("保存图片失败")
		return "",err
	}

	imgDom.SetAttr("src",imageName[1:])

	contentHtml,err:=doc.Html()
	if err != nil {
		beego.Info("从goquery解析doc转化content失败！")
	}

	beego.Info("转化后的content为： ",contentHtml)

	return contentHtml,nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
