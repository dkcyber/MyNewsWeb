package models

import (
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"time"
)
//用户结构体
type User struct {
	Id int
	UserName string `orm:"unique"`
	Pwd string
	Articles []*Article `orm:"rel(m2m)"`

}

//文章结构体
type Article struct {
	Id int `orm:"pk;auto"`
	Title string `orm:"size(100)"`
	Content string `orm:"size(500)"`
	Time time.Time `orm:"type(datetime);auto_now"`
	ReadCount int  `orm:"default(0)"`
	Image string `orm:"null"`
	ArticleType *ArticleType `orm:"rel(fk)"`
	Users []*User `orm:"reverse(many)"`
}


type ArticleType struct {
	Id int
	TypeName string `orm:"size(100)"`
	Articles []*Article `orm:"reverse(many)"`


}

//初始化数据库，在“main.go中导入自动运行包下init函数”
func init(){
	//1.注册数据库,需导入mysql数据库驱动
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/mynewsweb?charset=utf8&loc=Local")
	//2.注册对象，每新增一个orm映射数据库表均在此添加
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	//3.运行同步数据库
	orm.RunSyncdb("default",false,true)

}