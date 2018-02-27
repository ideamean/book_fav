package controllers

import (
	"book_fav/dao"
	"book_fav/douban"
	"book_fav/model"
	"fmt"
	"os"
	"strings"
)

type BookController struct {
	BaseController
}

//保存书架
func (c *BookController) AddPurchase() {
	uid := c.TokenInfo.Uid
	bid, _ := c.GetInt64("bid")
	remark := c.GetString("remark")

	var state model.State

	if bid <= 0 {
		state.Errno = 100
		state.Errmsg = "bid参数不正确"
		c.Data["json"] = state
		c.ServeJSON()
		return
	}
	bookInfo, err := dao.GetBookInfoById(bid)
	if err != nil {
		state.Errno = 101
		state.Errmsg = fmt.Sprintf("%s", err)
		c.Data["json"] = state
		c.ServeJSON()
		return
	}

	if bookInfo.Bid <= 0 {
		state.Errno = 102
		state.Errmsg = "该已书籍不存在"
		c.Data["json"] = state
		c.ServeJSON()
		return
	}

	err = dao.AddPurchase(uid, bid, remark)
	if err != nil {
		state.Errno = 100
		state.Errmsg = fmt.Sprintf("%s", err)
	} else {
		state.Errmsg = "success"
	}
	c.Data["json"] = state
	c.ServeJSON()
}

//搜索书库
func (c *BookController) Search() {
	uid := c.TokenInfo.Uid
	isbn := strings.TrimSpace(c.GetString("isbn"))

	var state model.State

	type Result struct {
		Purchase struct {
			IsPurchase int                    `json:"is_purchase"`
			Info       model.BookPurchaseInfo `json:"info"`
		} `json:"purchase"`
		BookInfo model.BookInfo `json:"book_info"`
	}

	var r Result

	if isbn == "" {
		state.Errno = 100
		state.Errmsg = "isbn不能为空"
		c.Data["json"] = state
		c.ServeJSON()
		return
	}

	info, err := dao.GetBookInfoByIsbn(isbn)
	if err != nil {
		state.Errno = 101
		state.Errmsg = fmt.Sprintf("%s", err)
		c.Data["json"] = state
		c.ServeJSON()
		return
	}

	if info.Bid > 0 {
		//判断是否已购
		purcharseInfo, err := dao.IsPurchase(uid, info.Bid)
		if err != nil {
			state.Errno = 101
			state.Errmsg = fmt.Sprintf("%s", err)
			c.Data["json"] = state
			c.ServeJSON()
			return
		}
		if purcharseInfo.Id > 0 {
			r.Purchase.IsPurchase = 2 //已购
			r.Purchase.Info = purcharseInfo
		} else {
			r.Purchase.IsPurchase = 1 //未购
		}
		r.BookInfo = info
		state.Data = r
		c.Data["json"] = state
		c.ServeJSON()
		return
	}

	//从豆瓣获取信息
	doubanInfo, err := douban.GetBookInfoByIsbn(isbn)
	if err != nil {
		state.Errno = 102
		state.Errmsg = fmt.Sprintf("%s", err)
		c.Data["json"] = state
		c.ServeJSON()
		return
	}

	if doubanInfo.DoubanId <= 0 {
		state.Errno = 103
		state.Errmsg = "未查找到该书籍"
		c.Data["json"] = state
		c.ServeJSON()
		return
	}

	bookInfo, err := dao.AddBookList(doubanInfo)
	if err != nil {
		state.Errno = 104
		state.Errmsg = fmt.Sprintf("%s", err)
		c.Data["json"] = state
		c.ServeJSON()
		return
	}
	r.BookInfo = bookInfo
	state.Data = r
	c.Data["json"] = state
	c.ServeJSON()
}

//获取书籍列表
func (c *BookController) GetBookList() {
	uid := c.TokenInfo.Uid
	searchType := c.GetString("type")
	page, _ := c.GetInt("page")
	pageSize, _ := c.GetInt("size")

	if searchType != "1" && searchType != "2" {
		searchType = "1"
	}

	var state model.State
	result, err := dao.GetBookList(searchType, uid, page, pageSize)
	if err != nil {
		state.Errno = 100
		state.Errmsg = fmt.Sprintf("%s", err)
		c.Data["json"] = state
		c.ServeJSON()
		return
	}

	state.Data = result
	c.Data["json"] = state
	c.ServeJSON()
}

//获取条形码
func (c *BookController) GetBarCode() {
	var state model.State

	f, h, err := c.GetFile("barImageFile")
	if err != nil {
		state.Errno = 100
		state.Errmsg = err.Error()
		c.Data["json"] = state
		c.ServeJSON()
		return
	}
	f.Close()

	tmpFile := "/tmp/" + h.Filename
	err = c.SaveToFile("barImageFile", tmpFile)
	if err != nil {
		state.Errno = 101
		state.Errmsg = err.Error()
		c.Data["json"] = state
		c.ServeJSON()
		return
	}

	result, err := dao.BarDecodeByImage(tmpFile)
	if err != nil {
		state.Errno = 102
		state.Errmsg = err.Error()
		c.Data["json"] = state
		c.ServeJSON()
		return
	}
	os.RemoveAll(tmpFile)
	state.Data = result
	c.Data["json"] = state
	c.ServeJSON()
	return
}
