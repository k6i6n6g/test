package controllers

import (
	"github.com/gomodule/redigo/redis"
	"github.com/astaxie/beego"
	"reflect"

)
type GoRedis struct {
	beego.Controller
}
func(this * GoRedis)ShowGet(){
	//打开     tcp协议      ：端口
	conn,err:=redis.Dial("tcp",":6379")
	//关闭
	defer conn.Close()
	if err!=nil{
		beego.Error("redis数据库链接失败",err)
	}

	//操作数据库                               属性
	resp,err:=conn.Do("mget","class1","aa")
	//万能接口
	re,err:=redis.Values(resp,err)
	//切片赋值的时侯声明变量
	var string1 string
	var int1  int
	//把接口函数放在string中
	redis.Scan(re,&string1,&int1)
	beego.Info("回复是：",string1,int1)
	beego.Info(reflect.TypeOf(int1))
}
