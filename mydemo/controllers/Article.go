package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"mydemo/models"
	"time"
	"path"
	"math"
	"bytes"
	"encoding/gob"
	"github.com/gomodule/redigo/redis"
)
type  ArticleController struct {
	beego.Controller
}
//显示登陆界面
func (this * ArticleController) ShowIndex(){
	userName:=this.GetSession("userName")
	if userName==nil{
		this.Redirect("/login",302)
		return
	}
	//更新内容到index中
	o:=orm.NewOrm()
	//指定表
	qs:=o.QueryTable("Article")


	//定义一个切片
	var articles []models.Article
	//qs里面的东西放在articles里面
	//_,err:=qs.All(&articles)
	//if  err !=nil{
	//	beego.Error("查询所有文章错误")
	//	this.TplName="index.html"
	//	return
	//}

	//每页分两个
	page :=2

	//首页和末页
	pageIndex,err:=this.GetInt("pageIndex")
	//如果页码点进是为没有 就改成第一页
	if err!=nil{
		pageIndex=1
	}
	//页码重0开始 按照一定规律进行的
	start :=page*(pageIndex-1)
	//获取属性标签
	typeName:=this.GetString("select")
	//添加科技娱乐属性
	var count int64
	var err1 error
	if typeName==""{
		//orm 中一对多的查询是惰性查询
		//1获取总记录和总数据    RelatedSel改变惰性    .Filter 相当与where    count是累计文章的个数
		count,err1=qs.RelatedSel("ArticleType").Count()
		qs.Limit(page,start).RelatedSel("ArticleType").All(&articles)
	}else {
		//1获取总记录和总数据

		count,err1=qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
		//获取数据（获取几条数据，从第几条开始获取数据）
		qs.Limit(page,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	}
	if err1!=nil {
		beego.Error("查询数据条目书错误")
		this.TplName="index.html"
		return
	}
	//页面向上取整  因为只有有小数时才会向上取整 所以要变成小数
	pageCount:=math.Ceil(float64(count)/float64(page))
	//传递数据
	this.Data["count"]=count
	this.Data["pageCount"]=int(pageCount)

	this.Data["typeName"]=typeName
	///获取所用类型数据
	var articleTypes []models.ArticleType
	//o.QueryTable("ArticleType").All(&articleTypes)
	conn,err:=redis.Dial("tcp","192.168.109.133:6379")
	if err!=nil{
		beego.Error("redis数据库链接失败",err)
	}
	//
	resp,err:=redis.Bytes(conn.Do("get","articleTypes"))
	//定义一个解码器
	dec:=gob.NewDecoder(bytes.NewReader(resp))
	//解码
	dec.Decode(&articleTypes)
	//如果切片的长度为0 就重新获取
	if len(articleTypes)==0{
		//获取所有类型
		o.QueryTable("articleType").All(&articleTypes)
		//序列化存储
		var buffer bytes.Buffer
		//进行编码
		enc:=gob.NewEncoder(&buffer)
		//进行编译
	    enc.Encode(&articleTypes)
	    //操作              建立         key            values    传到上面去
	    conn.Do("set","articleTypes",buffer.Bytes())
	    beego.Info("从mysq中获得数据")

	}
	/*///序列化和反序列化
	//要有一个容器，用来接受编码之后的字节流
	var buffer bytes.Buffer
	//要有一个编码器
	enc:=gob.NewEncoder(&buffer)
	//编码
	enc.Encode(&articleTypes)
	conn.Do("set","articleTypes",buffer.Bytes())
	resp,err:=conn.Do("get","articleTypes")
	//先获取字节流数据
	types,err:=redis.Bytes(resp,err)
	//获取解码器
	dec:=gob.NewDecoder(bytes.NewReader(types))
	//解码
	var testTypes []models.ArticleType
	dec.Decode(&testTypes)
	beego.Info(testTypes)
*/
	this.Data["articleTypes"]= articleTypes

	beego.Info("==========",articleTypes)

    this.Data["pageIndex"]=pageIndex
	this.Data["articles"]=articles
	this.TplName="index.html"
}
//添加文章
func(this * ArticleController) ShowAdd(){
	o:=orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"]= articleTypes
	this.TplName="add.html"
}
//处理添加文章
func (this * ArticleController)HandleAdd(){
	typeName:=this.GetString("select")
	//获取数据
	//标题
	title :=this.GetString("articleName")
	//主要内容
	content :=this.GetString("content")
	//图片路径   head 是图片的详细信息
	file,head,err:=this.GetFile("uploadname")
	if head.Filename == ""{
		//处理数据
		//数据库的插入操作
		//获取orm对象
		o :=orm.NewOrm()
		//获取一个插入对象
		var  article models.Article
		//给对象插入值
		article.Title=title
		article.Content=content
		//插入数据库中
		o.Insert(&article)
		//返回数据  用这个来直接进行跳转转换
		this.Redirect("/article/index",302)
        //读完了就直接退出，不执行下面的了
		return
	}
	//关闭文件对象
	defer file.Close()

	//校验数据  题目   内容
	 if title ==""|| content =="" || err != nil {
	 	//把后面的文字传到前面页面
		 this.Data["errmsg"]="添加文章失败，请重新添加！"
	 	this.TplName="add.html"
	 	return
	 }

	 //解决文件覆盖1问题
	 //时间格式的文件名    时间 格式
	 fileName :=time.Now().Format("2016-01-02-15-04-05")
	 //获取后缀名 ext       图片的文件名 head中有Filename
	 ext :=path.Ext(head.Filename)
	 beego.Info(head.Filename,ext)
	 //文件类型也需要校验
	 if ext !=".jpg"  && ext !=".png" && ext !=".jpeg"{
	 	this.Data["errmsg"]="上传图片格式不正确，请重新上传"
	 	this.TplName="add.html"
	 	return
	 }
	 //文件大小校验
	 //5000000单位b
	 //5MB 5*1024KB  1KB=1024B  5MB=5*1024*1024B
	 if head.Size>5*1024*1024{
	 	this.Data["errmsg"]="上传图片过大，请重新上传"
	 	this.TplName="add.html"
	 	return
	 }



	 //处理数据
	 //数据库的插入操作
	 //获取orm对象
	 o :=orm.NewOrm()
	 //获取一个插入对象
	 var  article models.Article
	 //给对象插入值
	 article.Title=title
	 article.Content=content
	 article.Img="/static/img/"+fileName+ext
	//外建
	var articleType models.ArticleType
	articleType.TypeName=typeName
    o.Read(&articleType,"TypeName")
    article.ArticleType=&articleType


	//插入数据库中
	 o.Insert(&article)
	 //返回数据  用这个来直接进行跳转转换
	 this.Redirect("/article/index",302)
 //读完了就直接退了，不执行下面的了
     return

}
//内容详情
func(this *ArticleController)ShowContent(){
	//文章的id号为主键
	articleId,err:=this.GetInt("articleId")
	if err !=nil{
		beego.Error("请求链接失败")
		this.TplName="index.html"
		return
	}
	//获取对象
	o:=orm.NewOrm()
	//获取数据库
	var article models.Article
	//赋值
	article.Id=articleId
	//读取文章
	err=o.Read(&article)
	if err!=nil{
		beego.Error("查询文章失败")
		this.TplName="index.html"
		return
	}
	//读取文章时计数
	article.Count+=1
	//更新文章
	o.Update(&article)
	//需要添加用户浏览记录
	//多对多数据添加
	m2m:=o.QueryM2M(&article,"Users")
	//插入的是用户的对象
	var user models.User1
	//记录用户名
	userName:=this.GetSession("userName")
	//赋值过来
	user.Name=userName.(string)
	o.Read(&user,"Name")
	//上面的主要是获取user 在这边添加
	m2m.Add(user)

	//获取浏览记录                为多对多新创建的新接口
	//o.LoadRelated(&article,"Users")
	//第二种多对多查询方法
	qs:=o.QueryTable("User1")
	//声明一个装函数的
	var users []models.User1
	//where  Articles []*Article `orm:"reverse(many)"` 第一个为条件  第二个具体的值
	//distinct 为去重   all 为装载上面声名工具
	qs.Filter("Articles__Article__Id",article.Id).Distinct().All(&users)
	//传到content。html中的最经浏览下面
	this.Data["users"]=users
	//把前面的Article通过RelatedSel关联ArticleType  通过Filter过滤当里面两个值相等时就好  第一个为ArticleType的id 第二个为Article的
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("Id", articleId).All(&article)
	//传送数据
	this.Data["article"]=article
	this.TplName="content.html"
}
//更新文章页面
func(this *ArticleController)ShowUpdate(){
	//获取数据
	articleId,err:=this.GetInt("articleId")
	//校验数据
	if err!=nil{
		beego.Error("请求链接错误")
		this.Redirect("/article/index",302)
		return
	}
	//处理数据
	//1.获取orm
	o:=orm.NewOrm()
	//获取更新对象
	var article models.Article
	article.Id=articleId
	err=o.Read(&article)
	if err!=nil{
		beego.Error("更新文章不存在")
		this.Redirect("/article/index",302)
		return
	}
	//返回数据
	this.Data["article"]=article
	this.TplName="update.html"

}
//包装函数
func UploadFunc(this *ArticleController,filepath string)string{
	//图片路径   head 是图片的详细信息
	file,head,err:=this.GetFile("uploadname")
	//关闭文件对象
	defer file.Close()
	if err!=nil{
		return ""
	}
	//解决文件覆盖1问题
	//时间格式的文件名    时间 格式
	fileName :=time.Now().Format("2016-01-02-15-04-05")
	//获取后缀名 ext       图片的文件名 head中有Filename
	ext :=path.Ext(head.Filename)
	beego.Info(head.Filename,ext)
	//文件类型也需要校验
	if ext !=".jpg"  && ext !=".png" && ext !=".jpeg"{
		this.Data["errmsg"]="上传图片格式不正确，请重新上传"
		this.TplName="add.html"
		return ""
	}
	//文件大小校验
	//5000000单位b
	//5MB 5*1024KB  1KB=1024B  5MB=5*1024*1024B
	if head.Size>5*1024*1024{
		this.Data["errmsg"]="上传图片过大，请重新上传"
		this.TplName="add.html"
		return ""
	}
	//把图片储存起来   名称  图片的路径  时间字符串  后缀名
	this.SaveToFile("uploadname","./static/img/"+fileName+ext)
	return "./static/img/"+fileName+ext
}
//添加文章
func (this *ArticleController)HandleUpdate()  {
	//获取数据
	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	articleId,err:=this.GetInt("articleId")
	fileAddr:=UploadFunc(this,"uploadname")
	//校验数据
	if articleName==""||content==""|| err!=nil{
		this.Data["errmsg"]="上传数据失败"
		this.TplName="update.html"
		return
	}

	//数据库更新处理数据
	//获取对象
	o:=orm.NewOrm()
	//获取数据
	var article models.Article
	//查询操作
	article.Id=articleId

	//看是否更新l
	err=o.Read(&article)
	if err!=nil{
		this.Data["errmsg"]="文章更新失败"
		this.TplName="update.html"
		return
	}

	//赋值
	article.Content=content
	article.Title=articleName
	article.Img=fileAddr
	//更新数据库
	o.Update(&article)
	this.Redirect("/article/index",302)
}
//删除文章
func (this *ArticleController)ShowDelete(){
	articleId,err:=this.GetInt("articleId")
	if err !=nil{
		this.Data["errmsg"]="文章获取失败"
		this.Redirect("/article/index",302)
		return
	}
	o:=orm.NewOrm()
	var article models.Article
	article.Id=articleId
	_,err=o.Delete(&article)
	if err!=nil {
		this.Data["errmsg"]="文章删除失败"
		this.Redirect("/article/index",302)
		return
	}
	this.Redirect("/article/index",302)
}
//添加文章类型页面
func (this *ArticleController)ShowAddType()  {
	o:=orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(& articleTypes)
	this.Data["articleTypes"]= articleTypes
	this.TplName="addType.html"
}
//添加文章类型
func (this *ArticleController)HandleAddType () {
	//获取
	typeName:=this.GetString("typeName")
	//校验
	if typeName==""{
		beego.Error("获取类型失败")
		this.TplName="addType.html"
		return
	}
	//处理
	o:=orm.NewOrm()
	var articleType  models.ArticleType
	articleType.TypeName=typeName
	_,err:=o.Insert(&articleType)
	if err !=nil{
		beego.Error("插入数据类型失败")
		this.TplName="addType.html"
		return
	}
	//返回
	this.Redirect("/article/addType",302)
}
//删除文章类型的属性
func (this *ArticleController)DeleteType(){
	typeId,err:=this.GetInt("typeId")
	if err!=nil{
		beego.Error("删除文章类型链接错误")
		this.Redirect("/article/addType",302)
		return
	}
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id=typeId
	_,err=o.Delete(&articleType)
	if err!=nil{
		beego.Error("删除失败")
		this.Redirect("/article/addType",302)
		return
	}
	this.Redirect("/article/addType",302)
}