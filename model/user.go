package model

type UserInfo struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	NickName string `json:"nickName"`
	Dateline int    `json:"dateline"`
}
