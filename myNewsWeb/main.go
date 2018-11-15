package main

import (
	_ "myNewsWeb/routers"
	"github.com/astaxie/beego"
	_"myNewsWeb/models"
)

func main() {
	beego.AddFuncMap("PrePage",PrePage)
	beego.AddFuncMap("NextPage",NextPage)
	beego.Run()
}


func PrePage(nowPage int)int{
	if nowPage==1{
		return 1
	}
	return nowPage-1
}

func NextPage(nowPage int,pageCount int)int{
	if nowPage==pageCount{
		return pageCount
	}
	return nowPage+1
}
