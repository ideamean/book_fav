package model

type ApiBookListResult struct {
	Page         int             `json:"page"`
	PageSize     int             `json:"page_size"`
	BookList     []BookInfo      `json:"book_list"`
	PurchaseList map[int64]int64 `json:"purchase_list"`
}
