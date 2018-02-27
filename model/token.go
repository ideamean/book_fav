package model

type TokenInfo struct {
	Id       int64  `json:"id"`
	Token    string `json:"token"`
	Uid      int64  `json:"uid"`
	GenCode  int    `json:"gen_code"`
	Expire   int64  `json:"expire"`
	Status   int    `json:"status"`
	Dateline int    `json:"dateline"`
}
