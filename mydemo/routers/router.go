package routers

import (
	"mydemo/controllers"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
)

func init() {
    //添加规范性 ，不能随便登陆，要进入网页是必须要从登陆口走
    beego.InsertFilter("/article/*",beego.BeforeExec,filterFunc)
    beego.Router("/", &controllers.MainController{})
    //注册
    beego.Router("/web", &controllers.WebController1{},"get:Register1;post:Register2")
    //登陆
    beego.Router("/login", &controllers.WebController1{},"get:ShowLogin;post:HandleLogin")
    //文章首页
    beego.Router("/article/index",&controllers.ArticleController{},"get:ShowIndex")
    //添加文章
    beego.Router("/article/add",&controllers.ArticleController{},"get:ShowAdd;post:HandleAdd")
    //查看文章详情
    beego.Router("/article/content",&controllers.ArticleController{},"get:ShowContent")
    //编辑文件
    beego.Router("/article/update",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
    //删除文章
    beego.Router("/article/delete",&controllers.ArticleController{},"get:ShowDelete")
    //添加文章类型
    beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    //点击退出按钮
    beego.Router("/article/logout",&controllers.WebController1{},"get:Logout")
    //删除文章类型的属性
    beego.Router("/article/deleteType",&controllers.ArticleController{},"get:DeleteType")

    //写了一个redis
    beego.Router("/GoRedis",&controllers.GoRedis{},"get:ShowGet")

}
//第一步函数的调用        规定
var filterFunc= func(ctx *context.Context) {
    //获取数据
    userName :=ctx.Input.Session("userName")
    if   userName==nil{
        ctx.Redirect(302,"/login")
        return
    }
}
