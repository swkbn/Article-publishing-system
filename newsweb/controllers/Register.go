package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsweb/models"

)

type RegisterController struct {
		beego.Controller
}

func (this*RegisterController) Showregister() {

	this.TplName="register.html"
}
func (this*RegisterController)ShowrePost()  {

	userName:=this.GetString("userName")
	userPwd:=this.GetString("password")
	beego.Info(userName,userPwd)
	//判断数据
	if userPwd==""||userName=="" {

		beego.Error("数据不能为空")
		this.TplName="register.html"
		return
	}
	//把数据插入表中
	o:=orm.NewOrm()
	var sh models.NewWeb
	sh.Name=userName
	sh.Pow=userPwd
	o.Insert(&sh)
	//重定向到登陆页面
	this.Redirect("/login",302)
}
//登陆
func (this *RegisterController)Showlogin()  {

	this.TplName="login.html"
}

func (this*RegisterController)Hendellogin()  {
	//获取数据

	userName:=this.GetString("userName")
	pwd:=this.GetString("password")
	//校验数据
	if userName==""||pwd=="" {

		this.Data["err"] = "用户名和密码不能为空"
		this.TplName = "login.html"
		return
	}

	//操作数据
	var sh models.NewWeb
	sh.Name=userName
	//
	o:=orm.NewOrm()

	err:=o.Read(&sh,"name")
	if err!=nil {
		this.Data["err"] = "用户名错误"
		this.TplName = "login.html"
		return
	}
	//判断

	if pwd!=sh.Pow {
		this.Data["err"] = "密码错误"
		this.TplName = "login.html"
		return
	}

	//返回数据登陆成功
	//this.Ctx.WriteString("登陆成功")
	//进行判断
	this.SetSession("userName",userName)
	this.Redirect("/article/index",302)
}

func (this*RegisterController)Showlogout()  {

	this.Redirect("/login",302)
	this.DelSession("userName")

}
