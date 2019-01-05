package main

import (
	_ "newsweb/routers"
	"github.com/astaxie/beego"
	//_"newsweb/models"
)
func main() {
	//关联
	beego.AddFuncMap("pre",getpre)
	beego.AddFuncMap("next",getnext)
	beego.Run()
}
func getpre(pre int)int  {
	if pre-1<=0 {
		return pre
	}
	return pre-1
}

func getnext(next,cont int)int  {
	if next+1>cont {

		return next
	}
	return next+1
}