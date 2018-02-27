package controllers

import (
	"book_fav/douban"
	"book_fav/model"
)

type DoubanController struct {
	BaseController
}

func (c *DoubanController) Get() {
	isbn := c.Ctx.Input.Param(":id")

	var state model.State
	var result model.BookInfo
	result, err := douban.GetBookInfoByIsbn(isbn)
	if err != nil {
		state.Errno = 100
		state.Errmsg = err.Error()
	} else {
		state.Data = result
	}
	c.Data["json"] = state
	c.ServeJSON()
}
