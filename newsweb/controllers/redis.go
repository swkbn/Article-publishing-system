package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"newsweb/models"
)

type GoRedis struct {
	beego.Controller
}

func (this*GoRedis)ShowGet()  {
	//连接数据库
	conn,err:=redis.Dial("tcp","192.168.189.11:6379")
	if err!=nil {
		beego.Error("连接数据库失败")
		return
	}

	 da:="aaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	models.Undel(&da)

	//操作数据库

	resp,err:=conn.Do("mget","a1")
	//返回值处理
	//获取同一类型自动存为一个切片
	//re,err:=redis.Strings(resp,err)
	////获取不同类型
	var a int
	var b string
	re,err:=redis.Values(resp,err)

	redis.Scan(re,&b)
	beego.Info("返回值为：",b,a)


}

