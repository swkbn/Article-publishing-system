package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsweb/models"
	"time"
	"path"
	"math"
	"github.com/gomodule/redigo/redis"
	"encoding/gob"
	"bytes"


)

type AriacleController struct {
	beego.Controller
}

//展示首页
func (this *AriacleController) ShowIndex() {
	//需要判断
	userName := this.GetSession("userName")

	if userName == nil {
		this.Redirect("/login", 302)
		return
	}

	//获取数据
	o := orm.NewOrm()
	//指定表
	qs := o.QueryTable("Article")
	//一个容器
	var rq []models.Article
	//查询
	/*_,err:=qs.All(&rq)
	if err!=nil {
		beego.Error("查询错误")
		this.TplName="index.html"
		return
	}*/
	//处理分页
	//先获取总记录的数量

	//定义一页现实多少条数据
	pageIndex := 2

	//this.Data[""]

	//首页和末页实现

	page, err := this.GetInt("pageIndex")
	if err != nil {
		//给他一个初始页
		page = 1
	}

	this.Data["page"] = page
	//获取数据
	start := pageIndex * (page - 1)
	//获取选中的类型
	typeName := this.GetString("select")
	//把typename传输到前段
	this.Data["typeName"] = typeName
	//orm中一对多的查询是惰性查询，
	var count int64
	var err1 error
	if typeName == "" {
		qs.Limit(pageIndex, start).RelatedSel("ArticleType").All(&rq)
		count, err1 = qs.RelatedSel("ArticleType").Count()
	} else {
		qs.Limit(pageIndex, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&rq)
		count, err1 = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
	}

	if err1 != nil {
		beego.Error("获取数量失败")
		this.TplName = "index.html"
		return
	}

	//获取总页数
	pageCount := math.Ceil(float64(count) / float64(pageIndex))

	beego.Info(count, "总数据数", pageCount, "页数")

	this.Data["pageCount"] = int(pageCount)
	this.Data["count"] = count

	//传递数据
	this.Data["articles"] = rq
	//this.Data["ArticleType"]=ArticleType
	//获取所有数据类型
	var sj [] models.ArticleType
	//连接数据库
	conn,err:=redis.Dial("tcp","192.168.189.11:6379")
	if err!=nil {
		beego.Error("连接redis数据库失败")
	}
	defer conn.Close()
	resp,err:=redis.Bytes(conn.Do("get","articleTypes"))
	//定义一个解码器
	dec:=gob.NewDecoder(bytes.NewReader(resp))
	//解码
	dec.Decode(&sj)


	beego.Info("获取数据为：",sj)

	if len(sj)==0 {
		//从mysql数据库中获取数据
		o.QueryTable("ArticleType").All(&sj)
		//序列化存储
		//定义一个容器
		var buffer bytes.Buffer
		//定义一个加码器
		enc:=gob.NewEncoder(&buffer)
		//解码
		enc.Encode(&sj)
		//把数据存入redis数据库中
		conn.Do("set","articleTypes",buffer.Bytes())
		//
		beego.Info("从数据库中获取数据")
	}







	//连接redis
	/*conn,err:=redis.Dial("tcp",":6379")
	if err!=nil {
		beego.Error("redis数据库连接失败")
	}
	//序列化和反序列化
	//要有一个容器
	var buffer bytes.Buffer
	//创建一个编码器
	enc:=gob.NewEncoder(&buffer)
	//编码
	enc.Encode(&sj)
	//存入数据库中
	conn.Do("set","sj",buffer.Bytes())
	//从数据库中读取
	resp,err:=conn.Do("get","sj")
	//先获取字节流数据
	types,err:=redis.Bytes(resp,err)
	//获取解码器
	dec:=gob.NewDecoder(bytes.NewReader(types))
	//解码
	//创建一个容器
	var testTypes []models.ArticleType
	dec.Decode(&testTypes)
	beego.Info("解码后数据",testTypes)*/
	beego.Info("req",rq)
	//把数据传输给前端
	this.Data["sj"] = sj

	this.TplName = "index.html"
}

//显示添加
func (this *AriacleController) ShowAdd() {
	//添加类型的数据
	o := orm.NewOrm()
	var Articletype []models.ArticleType
	o.QueryTable("ArticleType").All(&Articletype)
	this.Data["ArticleType"] = Articletype

	this.TplName = "add.html"
}

//添加插入数据
func (this *AriacleController) HelderAdd() {
	//获取类型数据
	selectName := this.GetString("select")

	//接受数据
	//标题
	article := this.GetString("articleName")
	//内容
	conten := this.GetString("content")
	//图片
	file, head, err := this.GetFile("uploadname")
	defer file.Close()
	//校验数据
	if article == "" || conten == "" || err != nil {
		this.Data["errmsg"] = "添加数据失败，请重新添加"
		this.TplName = "add.html"
		return
	}
	//文件存在覆盖的问题
	//利用当前时间
	fileName := time.Now().Format("2006-01-02-15-04-05")
	//取文件的后缀
	ext := path.Ext(head.Filename)
	//校验文件类型
	if ext != ".jpg" && ext != "png" {
		this.Data["errmsg"] = "上传文件格式不正确"
		this.TplName = "add.html"
		return
	}
	//校验文件大小
	if head.Size > 5000000 {
		this.Data["errmsg"] = "已超过文件大小上限，请重新添加"
		this.TplName = "add.html"
		return
	}
	//把图片存起来
	this.SaveToFile("uploadname", "./static/img/"+fileName+ext)
	//赋值
	o := orm.NewOrm()
	var sh models.Article
	sh.Title = article
	sh.Counter = conten
	sh.Img = "/static/img/" + fileName + ext
	//插入文章类型
	var articleType models.ArticleType
	//把获取到的类型名进行赋值
	articleType.TypeName = selectName
	//根据查询
	o.Read(&articleType, "TypeName")
	//把文章内容对应的对象进行赋值
	sh.ArticleType = &articleType

	//插入数据
	_, err = o.Insert(&sh)


	//添加成功后跳转到主页
	this.Redirect("/article/index", 302)
}

//显示详情页

func (this *AriacleController) ShowContent() {

	//获取id
	articleId, err := this.GetInt("articleId")
	if err != nil {
		beego.Error("获取数据失败")
		this.TplName = "index.html"
		return
	}

	//数据处理
	o := orm.NewOrm()

	var sj models.Article

	sj.Id = articleId

	err = o.Read(&sj)
	if err != nil {
		beego.Error("查询数据失败")
		this.Redirect("/index", 302)
		return
	}
	//阅读量处理
	sj.Count += 1
	//更新数据
	o.Update(&sj)

	//多对多的数据添加(多对多表的操作)
	m2m := o.QueryM2M(&sj, "Users")
	//插入用户对象
	var user models.NewWeb
	userName := this.GetSession("userName")
	user.Name = userName.(string)
	o.Read(&user, "Name")
	//插入数据
	m2m.Add(user)
	//获取浏览记录
	//o.LoadRelated(&sj,"Users")
	//多对多查询操作
	//指定表
	qs := o.QueryTable("NewWeb")

	var users []models.NewWeb
	//条件查询	//进行去重
	qs.Filter("Articles__Article__Id", sj.Id).Distinct().All(&users)
	//把数据传输到前端
	this.Data["users"] = users
	//返回数据
	this.Data["sj"] = sj
	this.TplName = "content.html"

}

//显示更新
func (this *AriacleController) ShowUpdate() {

	articleId, err := this.GetInt("articleId")
	if err != nil {
		beego.Error("获取id失败")
		this.Redirect("/article/index", 302)
		return
	}
	//获取
	o := orm.NewOrm()
	var sj models.Article
	//赋值
	sj.Id = articleId
	//读取数据
	o.Read(&sj)
	//返回数据
	this.Data["article"] = sj
	this.TplName = "update.html"

}

//处理图片进行封装
func Ltder(this *AriacleController, uploadname string) string {

	//图片
	file, head, err := this.GetFile(uploadname)
	if err != nil {
		beego.Error("处理图片出错")
		return ""
	}
	defer file.Close()
	//文件存在覆盖的问题
	//利用当前时间
	fileName := time.Now().Format("2006-01-02-15-04-05")
	//取文件的后缀
	ext := path.Ext(head.Filename)
	//校验文件类型
	if ext != ".jpg" && ext != "png" {
		this.Data["errmsg"] = "上传文件格式不正确"
		this.TplName = "add.html"
		return ""
	}
	//校验文件大小
	if head.Size > 5000000 {
		this.Data["errmsg"] = "已超过文件大小上限，请重新添加"
		this.TplName = "add.html"
		return ""
	}
	//把图片存起来
	this.SaveToFile(uploadname, "./static/img/"+fileName+ext)
	return "/static/img/" + fileName + ext
}

//对更新进行操作
func (this *AriacleController) HeaderlUpdate() {
	//获取数据和id
	articleId, err := this.GetInt("articleId")
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	//图片操作
	heder := Ltder(this, "uploadname")
	//读取数据
	o := orm.NewOrm()
	var sj models.Article
	sj.Id = articleId
	err = o.Read(&sj)
	if err != nil {
		beego.Error("操作的文章不存在")
		return
	}
	//对数据进行更改
	sj.Title = articleName
	sj.Counter = content
	sj.Img = heder
	o.Update(&sj)
	//添加成功后跳转到主页
	this.Redirect("/article/index", 302)
}

//删除操作

func (this *AriacleController) ShowDelete() {
	//获取id'

	articleId, err := this.GetInt("articleId")
	if err != nil {
		beego.Error("获取失败")
		this.Redirect("/article/index", 302)
		return
	}
	//操作数据
	o := orm.NewOrm()
	var sj models.Article
	sj.Id = articleId
	//删除数据
	_, err = o.Delete(&sj)

	if err != nil {
		beego.Error("删除失败")
		this.Redirect("/article/index", 302)
		return
	}

	this.Redirect("/article/index", 302)
}

//显示添加分类操作

func (this *AriacleController) ShowAddtype() {
	//显示数据
	o := orm.NewOrm()
	var addType []models.ArticleType
	o.QueryTable("ArticleType").All(&addType)
	//返回数据
	this.Data["addType"] = addType

	this.TplName = "addType.html"

}

//添加分类操作

func (this *AriacleController) HenderAddtype() {
	typeName := this.GetString("typeName")
	if typeName == "" {
		beego.Error("输入的数据为空，请重新插入")
		this.Redirect("/article/addType", 302)
		return
	}
	o := orm.NewOrm()
	var sj models.ArticleType
	sj.TypeName = typeName
	_, err := o.Insert(&sj)
	if err != nil {
		beego.Error("插入数据已存在", err)
		this.Redirect("/article/addType", 302)
		return
	}
	this.Redirect("/article/addType", 302)
}

//删除类型操作

func (this*AriacleController)DeleteType()  {
	//获取数据
	Typeid,err:=this.GetInt("Typeid")

	if err!=nil {
		beego.Error("删除数据不存在",err)
		this.Redirect("/article/addType",302)
		return
	}
	//
	o:=orm.NewOrm()
	var addType models.ArticleType
	addType.Id=Typeid
	o.Delete(&addType)
	//
	this.Redirect("/article/addType",302)

}