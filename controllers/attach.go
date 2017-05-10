package controllers

import (
	"github.com/astaxie/beego"
	"net/url"
	"os"
	"io"
)

type AttachController struct {
	beego.Controller
}

func (this *AttachController)Get()  {
	//获取URL，如果有中文需要反编码
	filePath,err:=url.QueryUnescape(this.Ctx.Request.RequestURI[1:])
	if err!=nil{
		this.Ctx.WriteString(err.Error())
		return
	}
	f,err:=os.Open(filePath)
	if err!=nil{
		this.Ctx.WriteString(err.Error())
		return
	}
	beego.Info("打开文件：",filePath,"成功，显示之！")

	defer f.Close()
	_,err=io.Copy(this.Ctx.ResponseWriter,f)
}
