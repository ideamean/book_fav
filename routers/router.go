package routers

import (
	"book_fav/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/douban/book/isbn/:id", &controllers.DoubanController{})
	beego.Router("/api/book/search", &controllers.BookController{}, "*:Search")
	beego.Router("/api/book/addpurchase", &controllers.BookController{}, "*:AddPurchase")
	beego.Router("/api/book/list", &controllers.BookController{}, "*:GetBookList")
	beego.Router("/api/book/grtbarcode", &controllers.BookController{}, "*:GetBarCode")

	beego.Router("/api/user/login", &controllers.UserController{}, "*:Login")
	beego.Router("/api/user/logout", &controllers.UserController{}, "*:Logout")
	beego.Router("/api/user/info", &controllers.UserController{}, "*:Info")
}
