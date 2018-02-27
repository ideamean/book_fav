package controllers

import (
	"book_fav/dao"
	"book_fav/model"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	TokenInfo model.TokenInfo
	SkipLogin bool
}

func (c *BaseController) Prepare() {

	var state model.State

	token := strings.TrimSpace(c.Ctx.GetCookie("token"))

	if token != "" {
		tokenInfo, err := dao.GetTokenInfoByToken(token)
		if err != nil {
			state.Errno = 100
			state.Errmsg = err.Error()
			c.Data["json"] = state
			c.ServeJSON()
			return
		}
		if tokenInfo.Id > 0 {
			c.TokenInfo = tokenInfo
		}
	}

	//是否跳过登陆
	if c.SkipLogin {
		return
	}

	if c.TokenInfo.Id <= 0 {
		state.Errno = 100000
		state.Errmsg = "please login"
		c.Data["json"] = state
		c.ServeJSON()
		return
	}

	if c.TokenInfo.Expire < time.Now().Unix() {
		c.Logout()
		state.Errno = 100000
		state.Errmsg = "please login"
		c.Data["json"] = state
		c.ServeJSON()
		return
	}
}

//退出
func (c *BaseController) Logout() {
	c.Ctx.SetCookie("token", "", -1)
	return
}
