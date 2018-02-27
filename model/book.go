package model

type BookInfo struct {
	Bid int64 `json:"bid"`
	DoubanBookInfo
}

type BookPurchaseInfo struct {
	Id       int64  `json:"id"`
	Uid      int64  `json:"uid"`
	Bid      int64  `json:"Bid"`
	Remark   string `json:"remark"`
	Dateline int    `json:"dateline"`
}
