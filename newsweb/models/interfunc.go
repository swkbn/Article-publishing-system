package models

import "github.com/astaxie/beego"

type Sunder interface {
	Undel()

}
func Undel(a * string)  {
	beego.Info("打印数据为",*a)
}

