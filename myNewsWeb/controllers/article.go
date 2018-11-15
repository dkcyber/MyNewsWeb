package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"myNewsWeb/models"
	"math"
	"strconv"
)

type ArticleController struct {
	beego.Controller
}

//上传文件处理函数封装,返回值：上传到数据库的路径名
func UploadFileToServer(this *ArticleController,uploadname string)string{
	//校验上传的文件
	file,header,err:=this.GetFile(uploadname)
	if err!=nil{
		this.Data["errmsg"]="获取上传图片失败！请重新操作"
		this.TplName="add.html"
		return ""
	}
	defer file.Close()
	//1.图片大小
	if header.Size>500000{
		beego.Error("图片过大，请重新上传！")
		return ""
	}

	//2.图片格式
	fileExt:=path.Ext(header.Filename)
	if fileExt!=".jpg" && fileExt!=".png" && fileExt!=",jpeg"{
		beego.Error("图片格式错误，请重新上传！")
		return ""
	}

	//3.图片重名
	fileName:=time.Now().Format("20060102150405")+fileExt
	err=this.SaveToFile("uploadname","./static/image/"+fileName)
	if err!=nil{
		beego.Error("上传图片失败，请重新上传！")
		return ""
	}
	return "/static/image/"+fileName
}

//显示文章列表主页信息
func (this *ArticleController)ShowArticleList(){
	//判断登陆session是否存在，无则跳转回登陆页
	userName:=this.GetSession("userName")
	if userName==nil{
		this.Redirect("/login",302)
		return
	}
	this.Data["userName"]=userName.(string)

	//获取数据
	o:=orm.NewOrm()
	typeName:=this.GetString("select")
	qs:=o.QueryTable("Article")
	var articleCount int64
	if typeName==""{
		articleCount,_=qs.Count()
	}else{
		articleCount,_=qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
	}

	pageSize:=2
	pageCount:=int(math.Ceil(float64(articleCount)/float64(pageSize)))
	nowPage,err:=this.GetInt("nowPage")
	if err!=nil{
		nowPage=1
	}

	var articles []models.Article
	if typeName==""{
		qs.Limit(pageSize,(nowPage-1)*2).RelatedSel("ArticleType").All(&articles)
	}else {
		qs.Limit(pageSize,(nowPage-1)*2).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	}

	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)

	this.Data["typeName"]=typeName
	this.Data["articleTypes"]=articleTypes
	this.Data["articleCount"]=articleCount
	this.Data["pageCount"]=pageCount
	this.Data["nowPage"]=nowPage
	this.Data["articles"]=articles
	this.Layout="layout.html"
	this.TplName="index.html"
}

//显示添加文章页面
func (this *ArticleController)ShowArticleAdd(){
	o:=orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"]=articleTypes
	this.Layout="layout.html"
	this.TplName="add.html"

}

//处理添加文章请求
func (this *ArticleController)HandleArticleAdd(){
	//1.获取数据
	//2.校验数据
	//3.处理数据
	//4.返回数据

	//校验标题和内容
	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	articleTypeName:=this.GetString("select")

	if articleName=="" || content=="" || articleTypeName==""{
		this.Data["errmsg"]="标题和内容不能为空!！"
		this.TplName="add.html"
		return
	}



	//校验上传的文件
	file,header,err:=this.GetFile("uploadname")
	if err!=nil{
		this.Data["errmsg"]="获取上传图片失败！请重新操作"
		this.TplName="add.html"
		return
	}
	defer file.Close()
	//1.图片大小
	if header.Size>500000{
		this.Data["errmsg"]="图片过大，请重新上传！"
		this.TplName="add.html"
		return
	}

	//2.图片格式
	fileExt:=path.Ext(header.Filename)
	if fileExt!=".jpg" && fileExt!=".png" && fileExt!=",jpeg"{
		this.Data["errmsg"]="图片格式错误，请重新上传！"
		this.TplName="add.html"
		return
	}

	//3.图片重名
	fileName:=time.Now().Format("20060102150405")+fileExt
	err=this.SaveToFile("uploadname","./static/image/"+fileName)
	if err!=nil{
		this.Data["errmsg"]="上传图片失败，请重新上传！"
		this.TplName="add.html"
		return
	}

	o:=orm.NewOrm()
	var article models.Article
	article.Title=articleName
	article.Content=content
	article.Image="/static/image/"+fileName

	var articleType models.ArticleType
	articleType.TypeName=articleTypeName
	o.Read(&articleType,"TypeName")
	article.ArticleType=&articleType

	_,err=o.Insert(&article)
	if err!=nil{
		this.Data["errmsg"]="添加文章失败，请重新操作！"
		this.TplName="add.html"
		return
	}

	this.Redirect("/article/articleList",302)
}

//显示文章详情页
func (this *ArticleController)ShowArticleDetail(){
	//获取数据
	articleId,err:=this.GetInt("id")

	//校验数据
	if err!=nil{
		this.Data["errmsg"]="请求路径异常"
		this.TplName="index.html"
		return
	}

	//处理数据
	//获取orm对象
	o:=orm.NewOrm()
	//获取查询对象
	var article models.Article
	//给查询条件赋值
	article.Id=articleId
	//进行查询操作
	err=o.Read(&article)  //如果查询条件是ID的话，可以省略查询条件
	if err!=nil{
		this.Data["errmsg"]="获取文章信息异常"
		this.TplName="index.html"
		return
	}

	//------------多对多插入-------------

	//1.获取article对象

	//o:=orm.NewOrm()
	//var article models.Article
	//article.Id=articleId
	//o.Read(&article)

	//2.获取多对多m2m操作对象
	m2m:=o.QueryM2M(&article,"Users")

	//3.获取要插入的user对象
	var user models.User
	userName:=this.GetSession("userName").(string)
	user.UserName=userName
	o.Read(&user,"UserName")

	//4.插入user对象到m2m对象中,形成多对多关系
	m2m.Add(user)
	//------------多对多插入-------------

	//------------多对多查询-------------

	////第一种：简便方式，但是不返回querySeter对象,不能操作高级查询，如去重操作
	//o.LoadRelated(&article,"Users")

	//第二种：利用高级查询和过滤器，去重操作
	var users []models.User
	//利用反向查询过滤，查询有哪些用户读过这篇文章,Filter参数（"user表中条件__要查的表格__对应查询字段",查询字段的值）
	o.QueryTable("User").Filter("Articles__Article__Id",articleId).Distinct().All(&users)

	this.Data["users"]=users

	//------------多对多查询-------------

	//更新阅读量
	article.ReadCount++
	o.Update(&article)



	//返回数据
	this.Data["article"]=article
	this.TplName="content.html"
}

//显示文章编辑页
func (this *ArticleController)ShowArticleUpdate(){
	//获取数据
	articleId,err:=this.GetInt("id")

	errmsg:=this.GetString("errmsg")

	if errmsg!=""{
		this.Data["errmsg"]=errmsg
	}

	//校验数据
	if err!=nil{
		beego.Error("请求路径错误")
		this.Redirect("/article/articleList",302)
		return
	}
	//处理数据
	o:=orm.NewOrm()
	var article models.Article
	article.Id=articleId
	err=o.Read(&article)
	if err!=nil{
		beego.Error("请求路径错误！")
		this.Redirect("/article/articleList",302)
	}

	//返回数据

	this.Data["article"]=article
	this.TplName="update.html"
}

//处理文章编辑更新
func (this *ArticleController)HandleArticleUpdate(){
	//获取数据
	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	filePath:=UploadFileToServer(this,"uploadname")
	articleId,err2:=this.GetInt("id")
	//todo err2出错，应该提示页面错误，跳转到index.html

	//校验数据
	if articleName=="" || content=="" || filePath==""|| err2!=nil{
		errmsg:="内容不能为空！"
		this.Redirect("/article/articleUpdate?id="+strconv.Itoa(articleId)+"&errmsg="+errmsg,302)
		return
	}
	//处理数据
	//获取orm对象
	o:=orm.NewOrm()
	//获取查询对象
	var article models.Article
	//设置查询条件
	article.Id=articleId
	//查询
	err:=o.Read(&article)
	if err!=nil{
		errmsg:="获取更新文章失败"
		this.Redirect("/article/articleUpdate?id="+strconv.Itoa(articleId)+"&errmsg="+errmsg,302)
		return
	}
	//修改对象属性
	article.Title=articleName
	article.Content=content
	article.Image=filePath
	//更新数据库表对象
	_,err=o.Update(&article)
	if err!=nil{
		errmsg:="更新文章失败"
		this.Redirect("/article/articleUpdate?id="+strconv.Itoa(articleId)+"&errmsg="+errmsg,302)
		return
	}
	//返回数据
	this.Redirect("/article/articleList",302)
}

//处理文章删除
func (this *ArticleController)HandleArticleDelete(){
	//获取数据
	articleId,err:=this.GetInt("id")

	//校验数据
	if err!=nil{
		beego.Error("传递id错误")
		this.Redirect("/article/articleList",302)
		return
	}
	//处理数据
	//获取orm对象
	o:=orm.NewOrm()
	//获取删除对象
	var article models.Article
	//删除条件赋值
	article.Id=articleId
	//删除操作
	_,err=o.Delete(&article)
	if err!=nil{
		beego.Error("删除失败!")
		this.Redirect("/article/articleList",302)
		return
	}

	//返回数据
	this.Redirect("/article/articleList",302)
}

//显示添加分类页面
func (this *ArticleController)ShowArticleTypeAdd(){
	//获取数据
	var articleTypes []models.ArticleType
	o:=orm.NewOrm()
	_,err:=o.QueryTable("ArticleType").All(&articleTypes)
	//校验数据
	if err!=nil{
		beego.Error("获取文章分类失败！")
	}
	//处理数据

	//返回数据
	this.Data["articleTypes"]=articleTypes
	this.Layout="layout.html"
	this.TplName="addType.html"
}

//处理添加分类
func (this *ArticleController)HandleArticleTypeAdd(){
	//获取数据
	articleTypeName:=this.GetString("typeName")

	//校验数据
	if articleTypeName==""{
		beego.Error("分类名不能为空！")
		this.Redirect("/article/articleTypeAdd",302)
		return
	}
	//处理数据
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName=articleTypeName
	_,err:=o.Insert(&articleType)
	if err!=nil{
		beego.Error("添加失败！")
		this.Redirect("/article/articleTypeAdd",302)
		return
	}
	//返回数据
	this.Redirect("/article/articleTypeAdd",302)
}

//处理删除文章分类
func (this *ArticleController)HandleArticleTypeDelete(){
	//获取数据
	articleTypeId,err:=this.GetInt("id")

	//校验数据
	if err!=nil{
		beego.Error("请求路径id错误")
		this.Redirect("/article/articleTypeAdd",302)
		return
	}
	//处理数据
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id=articleTypeId
	_,err=o.Delete(&articleType)
	if err!=nil{
		beego.Error("删除失败！")
		this.Redirect("/article/articleTypeAdd",302)
		return
	}
	//返回数据
	this.Redirect("/article/articleTypeAdd",302)
}