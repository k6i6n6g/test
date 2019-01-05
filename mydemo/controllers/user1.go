package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"mydemo/models"
)
type WebController1 struct {
	beego.Controller
}
 func (this *WebController1)Register1(){
 	this.TplName="register.html"
 }
 func (this *WebController1)Register2(){
 	Name:=this.GetString("userName")
 	Pwd :=this.GetString("password")
 	beego.Info(Name,Pwd)

 	if Name== ""||Pwd== ""{
 		beego.Error("用户名或不能为空")
 		this.TplName="register.html"
 		return
	}

	o :=orm.NewOrm()
	var user1 models.User1
	user1.Name=Name
	user1.Pwd=Pwd
	connt,err :=o.Insert(&user1)
	 if err !=nil {
		 beego.Error("用户注册失败")
		 this.TplName="register.html"
		 return
	 }
	 beego.Info("条数为：",connt)
	 //this.Ctx.WriteString("注册成功")
	 this.Redirect("/login",302)


 }
//登陆业务
func(this * WebController1)ShowLogin(){
	//进行点击记住用户名
	//获取下面的key
	userName:=this.Ctx.GetCookie("userName")
	//进行判断如果为空会是什么情况不为空为什么情况
	if userName !=""{
		this.Data["userName"]=userName
		this.Data["checked"]="checked"
	}else{
		this.Data["userName"]=userName
		this.Data["checked"]=""
	}

	//一开始渲染的
	this.TplName="login.html"
}

//处理登陆页面
func(this * WebController1)HandleLogin(){
	//获取用户名
	user1Name :=this.GetString("userName")
	pwd :=this.GetString("password")

	if  user1Name=="" ||pwd == ""{
		this.Data["err"]="用户名或者密码不能为空"
		this.TplName="login.html"
		return
	}
	o:=orm.NewOrm()
	var user1  models.User1
	user1.Name=user1Name
	err:=o.Read(&user1,"Name")
	if err !=nil{
		this.Data["err"]="用户名不存在"
		this.TplName="login.html"
		return
	}
	if user1.Pwd !=pwd {
		this.Data["err"]="密码错误"
		this.TplName="login.html"
		return
	}
	//作用是设置cookie点击记住按钮   key    value          时间
	this.Ctx.SetCookie("userName",user1Name,3600*24)

	this.SetSession("userName",user1Name)
	//传递数据的
	//this.Ctx.WriteString("登陆成功")

	this.Redirect("/article/index",302)
}
//退出登陆
func(this *WebController1)Logout(){
	this.DelSession("userName")
	this.Redirect("/login",302)
}



