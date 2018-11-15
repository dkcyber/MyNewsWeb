package routers

import (
	"myNewsWeb/controllers"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
)

func init() {
    beego.InsertFilter("/article/*",beego.BeforeExec,funcFilter)
    beego.Router("/", &controllers.MainController{})
    beego.Router("/register", &controllers.UserController{},"get:ShowRegister;post:HandleRegister")
    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    beego.Router("/article/logout",&controllers.UserController{},"get:HandleLogout")
    beego.Router("/article/articleList",&controllers.ArticleController{},"get:ShowArticleList")
    beego.Router("/article/articleAdd",&controllers.ArticleController{},"get:ShowArticleAdd;post:HandleArticleAdd")
    beego.Router("/article/articleDetail",&controllers.ArticleController{},"get:ShowArticleDetail")
    beego.Router("/article/articleUpdate",&controllers.ArticleController{},"get:ShowArticleUpdate;post:HandleArticleUpdate")
    beego.Router("/article/articleDelete",&controllers.ArticleController{},"get:HandleArticleDelete")
    beego.Router("/article/articleTypeAdd",&controllers.ArticleController{},"get:ShowArticleTypeAdd;post:HandleArticleTypeAdd")
    beego.Router("/article/articleTypeDelete",&controllers.ArticleController{},"get:HandleArticleTypeDelete")
}

var funcFilter = func(ctx *context.Context) {
    userName:=ctx.Input.Session("userName")
    if userName==nil{
        ctx.Redirect(302,"/login")
    }else{

        ctx.Input.SetData("userName",userName.(string))
    }

}