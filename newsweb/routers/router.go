package routers

import (
	"newsweb/controllers"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
)

func init() {

    //定义一个过滤
    beego.InsertFilter("/article/*",beego.BeforeExec,filterFunc)


    beego.Router("/", &controllers.MainController{})
    
    beego.Router("/register",&controllers.RegisterController{},"get:Showregister;post:ShowrePost")

    beego.Router("/login",&controllers.RegisterController{},"get:Showlogin;post:Hendellogin")

    beego.Router("/article/index",&controllers.AriacleController{},"get:ShowIndex")
    beego.Router("/article/add",&controllers.AriacleController{},"get:ShowAdd;post:HelderAdd")

    //查看详情
    beego.Router("/article/content",&controllers.AriacleController{},"get:ShowContent")
    //更新操作
    beego.Router("/article/update",&controllers.AriacleController{},"get:ShowUpdate;post:HeaderlUpdate")
    //删除操作
    beego.Router("/article/delete",&controllers.AriacleController{},"get:ShowDelete")
    //添加分类
    beego.Router("/article/addType",&controllers.AriacleController{},"get:ShowAddtype;post:HenderAddtype")

    //退出登录
    beego.Router("/article/logout",&controllers.RegisterController{},"get:Showlogout")
    //删除类型
    beego.Router("/article/deleteType",&controllers.AriacleController{},"get:DeleteType")

    //连接数据库
    beego.Router("/redis",&controllers.GoRedis{},"get:ShowGet")

}

var filterFunc  = func(ctx*context.Context) {

    userName:=ctx.Input.Session("userName")

    if userName==nil {
        ctx.Redirect(302,"/login")
        return
    }
}
