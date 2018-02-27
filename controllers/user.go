package controllers

import (
	"book_fav/dao"
	"book_fav/model"
	"fmt"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"github.com/jordan-wright/email"
)

type UserController struct {
	BaseController
}

//初始化函数
func (c *UserController) Prepare() {
	var skipLogin = map[string]bool{
		"Login": true,
	}

	//跳过不需要登录的action
	_, actionName := c.GetControllerAndAction()
	_, ok := skipLogin[actionName]
	if ok {
		c.SkipLogin = true
	}

	c.BaseController.Prepare()
}

//获取用户信息
func (c *UserController) Info() {
	var state model.State
	info, err := dao.GetUserInfoByUid(c.TokenInfo.Uid)
	if err != nil {
		state.Errno = 100
		state.Errmsg = err.Error()
		c.Data["json"] = state
		c.ServeJSON()
		return
	}
	state.Data = info
	c.Data["json"] = state
	c.ServeJSON()
	return
}

//退出
func (c *UserController) Logout() {
	c.BaseController.Logout()
	var state model.State
	state.Errmsg = "success"
	c.Data["json"] = state
	c.ServeJSON()
}

//登陆
func (c *UserController) Login() {
	var state model.State
	email := strings.TrimSpace(c.GetString("email"))
	logType := strings.TrimSpace(c.GetString("log_type"))

	if email == "" {
		state.Errno = 100
		state.Errmsg = "请输入邮箱"
		c.Data["json"] = state
		c.ServeJSON()
		return
	}

	isvalid := checkEmail(email)
	if !isvalid {
		state.Errno = 101
		state.Errmsg = "邮箱格式不正确"
		c.Data["json"] = state
		c.ServeJSON()
		return
	}

	userinfo, err := dao.GetUserInfoByEmail(email)
	if err != nil {
		state.Errno = 102
		state.Errmsg = "无法获取用户信息"
		c.Data["json"] = state
		c.ServeJSON()
		return
	}

	type Result struct {
		LogType string `json:"log_type"`
	}

	var r Result
	r.LogType = logType

	state.Data = r

	if logType == "send_code" {
		genCode, err := dao.AddToken(userinfo.Id)
		if err != nil {
			state.Errno = 103
			state.Errmsg = "无法获取用户信息"
			c.Data["json"] = state
			c.ServeJSON()
			return
		}
		err = sendEmail(email, genCode)
		if err != nil {
			state.Errno = 103
			state.Errmsg = fmt.Sprintf("发送邮件失败: %s", err)
			c.Data["json"] = state
			c.ServeJSON()
			return
		}
		state.Errmsg = "验证码已发送到您的邮箱,请查收"
		c.Data["json"] = state
		c.ServeJSON()
		return
	}

	genCode, _ := c.GetInt("gen_code")
	if genCode <= 0 {
		state.Errno = 104
		state.Errmsg = "请输入邮箱验证码!"
		c.Data["json"] = state
		c.ServeJSON()
		return
	}
	tokenInfo, err := dao.CheckTokenByUid(userinfo.Id, genCode)
	if err != nil {
		state.Errno = 104
		state.Errmsg = err.Error()
		c.Data["json"] = state
		c.ServeJSON()
		return
	}
	if tokenInfo.Id <= 0 {
		state.Errno = 104
		state.Errmsg = "验证码不正确"
		c.Data["json"] = state
		c.ServeJSON()
		return
	}

	//setcookie
	maxAge := tokenInfo.Expire - time.Now().Unix()
	c.Ctx.SetCookie("token", tokenInfo.Token, maxAge)
	c.Data["json"] = tokenInfo
	c.ServeJSON()
}

func checkEmail(email string) (b bool) {
	if m, _ := regexp.MatchString("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+", email); !m {
		return false
	}
	return true
}

func sendEmail(to string, genCode int64) error {
	e := email.NewEmail()
	e.From = "chenlitao2008@qq.com"
	e.To = []string{to}
	e.Bcc = []string{"chenlitao2008@qq.com"}
	e.Subject = "家庭图书馆-登陆验证码"
	e.Text = []byte("验证码：" + fmt.Sprintf("%d", genCode))
	//e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	return e.Send("smtp.qq.com:587", smtp.PlainAuth("", "chenlitao2008@qq.com", "oeresnablptubgga", "smtp.qq.com"))
}
