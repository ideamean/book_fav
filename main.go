package main

import (
	_ "book_fav/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.SetStaticPath("/", "webroot")
	beego.SetStaticPath("/barcode.html", "webroot/barcode.html")
	beego.SetStaticPath("/js", "webroot/js")
	beego.Run()
}
