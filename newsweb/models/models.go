package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "database/sql"
	"time"
)

type NewWeb struct {
	Id       int
	Name     string     `orm:"unique"`
	Pow      string
	Articles []*Article `orm:"reverse(many)"`
}
type Article struct {
	Id          int
	Title       string       `orm:"size(100)"`
	Time        time.Time    `orm:"type(datetime);auto_now"`
	Count       int          `orm:"default(0)"`
	Counter     string
	Img         string       `orm:"null"`
	ArticleType *ArticleType `orm:"rel(fk);on_delete(do_nothing)"`
	Users       []*NewWeb    `orm:"rel(m2m)"`
	//Price float64 `orm:"digits(10);decimals(20)"`
}

type ArticleType struct {
	Id       int
	TypeName string     `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"`
}

func init() {

	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/under?charset=utf8")

	orm.RegisterModel(new(NewWeb), new(Article), new(ArticleType))
	orm.RunSyncdb("default", false, true)
}
