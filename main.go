package main

import (
	_ "book_fav/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.SetStaticPath("/", "webroot")
	beego.SetStaticPath("/js", "webroot/js")
	beego.Run()
}
