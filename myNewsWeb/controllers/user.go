package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"myNewsWeb/models"
	"encoding/base64"
)

type UserController struct{
	beego.Controller
}

//显示注册页
func (this *UserController)ShowRegister(){

	this.TplName="register.html"
}

//处理注册请求
func (this *UserController)HandleRegister(){
	//1.获取数据
	userName:=this.GetString("userName")
	pwd:=this.GetString("password")

	//2.校验数据
	if userName=="" || pwd==""{
		this.Data["errmsg"]="账号或密码不能为空！"
		this.TplName="register.html"
		return
	}

	//3.处理数据
	o:=orm.NewOrm()
	var user models.User
	user.UserName=userName
	user.Pwd=pwd
	_,err:=o.Insert(&user)
	if err!=nil{
		this.Data["errmsg"]="注册失败！请重新注册。。。"
		this.TplName="register.html"
		return
	}

	//4.返回数据
	this.Redirect("/login",302)
	//this.Ctx.WriteString("注册成功!")
}

//显示登陆页
func (this *UserController)ShowLogin(){
	dec:=this.Ctx.GetCookie("userName")
	userName,_:=base64.StdEncoding.DecodeString(dec)
	if string(userName)!=""{

		this.Data["userName"]=string(userName)
		this.Data["checked"]="checked"
	}else{
		this.Data["userName"]=""
		this.Data["checked"]=""
	}

	this.TplName="login.html"
}

//处理登陆请求
func (this *UserController)HandleLogin(){
	//1.获取数据
	userName:=this.GetString("userName")
	pwd:=this.GetString("password")
	//2.校验数据
	if userName=="" || pwd==""{
		this.Data["errmsg"]="账号或密码不能为空！"
		this.TplName="login.html"
		return
	}
	//3.处理数据
	o:=orm.NewOrm()
	var user models.User
	user.UserName=userName
	err:=o.Read(&user,"UserName")
	if err!=nil{
		this.Data["errmsg"]="账号错误，请重新输入"
		this.TplName="login.html"
		return
	}

	if user.Pwd!=pwd{
		this.Data["errmsg"]="密码错误，请重新输入！"
		this.TplName="login.html"
		return
	}

	//等校验完成后处理cookie
	remember:=this.GetString("remember")


	if remember=="on"{
		enc:=base64.StdEncoding.EncodeToString([]byte(userName))
		this.Ctx.SetCookie("userName",enc,3600*1)
	}else{
		this.Ctx.SetCookie("userName",userName,-1)
	}

	//4.返回数据
	this.SetSession("userName",userName)
	this.Redirect("/article/articleList",302)
	//this.Ctx.WriteString("登陆成功！")
}

//处理退出请求
func (this *UserController)HandleLogout(){
	this.DelSession("userName")
	this.Redirect("/login",302)
}