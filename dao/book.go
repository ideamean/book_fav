package dao

import (
	"book_fav/model"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//根据bid获取书信息
func GetBookInfoById(bid int64) (model.BookInfo, error) {
	var info model.BookInfo
	db, err := GetDB()
	if err != nil {
		return info, err
	}

	defer db.Close()

	fields := "id,douban_id,title,origin_title,alt_tite,publisher,pubdate,isbn13,isbn10,translator,pages,author,author_intro,catalog,summary,price,binding,image_large,image_medium,image_small,dateline"

	var imageLarge string
	var imageMedium string
	var imageSmall string
	var translatorStr string
	var authorStr string
	var dateline int
	err = db.QueryRow(fmt.Sprintf("select %s from book_list where id=?", fields), bid).Scan(&info.Bid, &info.DoubanId, &info.Title, &info.OriginTitle, &info.AltTitle, &info.Publisher, &info.Pubdate, &info.Isbn13, &info.Isbn10, &translatorStr, &info.Pages, &authorStr, &info.AuthorIntro, &info.Catalog, &info.Summary, &info.Price, &info.Binding, &imageLarge, &imageMedium, &imageSmall, &dateline)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return info, nil
		} else {
			return info, err
		}
	}

	info.Images = map[string]string{"large": imageLarge, "medium": imageMedium, "small": imageSmall}
	info.Author = strings.Split(authorStr, ",")
	info.Translator = strings.Split(translatorStr, ",")
	return info, nil
}

//根据isbn获取书信息
func GetBookInfoByIsbn(isbn string) (model.BookInfo, error) {
	var info model.BookInfo
	db, err := GetDB()
	if err != nil {
		return info, err
	}

	defer db.Close()

	isbnLen := len(isbn)

	var col string
	if isbnLen >= 13 {
		col = "isbn13"
	} else {
		col = "isbn10"
	}

	fields := "id,douban_id,title,origin_title,alt_tite,publisher,pubdate,isbn13,isbn10,translator,pages,author,author_intro,catalog,summary,price,binding,image_large,image_medium,image_small,dateline"
	sql := fmt.Sprintf("select %s from book_list where %s=?", fields, col)

	var imageLarge string
	var imageMedium string
	var imageSmall string
	var translatorStr string
	var authorStr string
	var dateline int
	err = db.QueryRow(sql, isbn).Scan(&info.Bid, &info.DoubanId, &info.Title, &info.OriginTitle, &info.AltTitle, &info.Publisher, &info.Pubdate, &info.Isbn13, &info.Isbn10, &translatorStr, &info.Pages, &authorStr, &info.AuthorIntro, &info.Catalog, &info.Summary, &info.Price, &info.Binding, &imageLarge, &imageMedium, &imageSmall, &dateline)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return info, nil
		} else {
			return info, err
		}
	}

	info.Images = map[string]string{"large": imageLarge, "medium": imageMedium, "small": imageSmall}
	info.Author = strings.Split(authorStr, ",")
	info.Translator = strings.Split(translatorStr, ",")
	return info, nil
}

//添加书库
func AddBookList(info model.BookInfo) (model.BookInfo, error) {
	db, err := GetDB()
	if err != nil {
		return info, err
	}

	defer db.Close()

	stmt, err := db.Prepare("insert into book_list set douban_id=?,title=?,origin_title=?,alt_tite=?,publisher=?,pubdate=?,isbn13=?,isbn10=?,translator=?,pages=?,author=?,author_intro=?,catalog=?,summary=?,price=?,binding=?,image_large=?,image_medium=?,image_small=?,dateline=?")
	if err != nil {
		return info, err
	}

	defer stmt.Close()

	translator := strings.Join(info.Translator, ",")
	author := strings.Join(info.Author, ",")
	pages, _ := strconv.ParseInt(info.Pages, 10, 64)
	ret, err := stmt.Exec(info.DoubanId, info.Title, info.OriginTitle, info.AltTitle, info.Publisher, info.Pubdate, info.Isbn13, info.Isbn10, translator, pages, author, info.AuthorIntro, info.Catalog, info.Summary, info.Price, info.Binding, info.Images["large"], info.Images["medium"], info.Images["small"], time.Now().Unix())
	if err != nil {
		return info, err
	}
	newId, err := ret.LastInsertId()
	if err != nil {
		return info, err
	}
	info.Bid = newId
	return info, nil
}

//判断是否已购
func IsPurchase(uid int64, bid int64) (model.BookPurchaseInfo, error) {
	var info model.BookPurchaseInfo
	db, err := GetDB()
	if err != nil {
		return info, err
	}

	defer db.Close()

	var dataId, dataUid, dataBid int64
	var dateline int
	var remark string
	err = db.QueryRow("select id,uid,bid,remark,dateline from book_purchase where uid=? and bid=?", uid, bid).Scan(&dataId, &dataUid, &dataBid, &remark, &dateline)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return info, nil
		} else {
			return info, err
		}
	}
	info.Id = dataId
	info.Uid = dataUid
	info.Bid = dataBid
	info.Remark = remark
	info.Dateline = dateline
	return info, nil
}

//添加已购
func AddPurchase(uid int64, bid int64, remark string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	defer db.Close()

	stmt, err := db.Prepare("insert ignore into book_purchase set uid=?,bid=?,remark=?,dateline=?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(uid, bid, remark, time.Now().Unix())
	return err
}

//获取书籍列表
func GetBookList(searchType string, uid int64, page int, pageSize int) (model.ApiBookListResult, error) {
	var result model.ApiBookListResult

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 || pageSize > 200 {
		pageSize = 200
	}

	result.Page = page
	result.PageSize = pageSize
	result.PurchaseList = make(map[int64]int64)

	db, err := GetDB()
	if err != nil {
		return result, err
	}

	defer db.Close()

	var sql string
	fields := "id,douban_id,title,origin_title,alt_tite,publisher,pubdate,isbn13,isbn10,translator,pages,author,author_intro,catalog,summary,price,binding,image_large,image_medium,image_small,dateline"
	fields = "a." + strings.Join(strings.Split(fields, ","), ",a.")
	if searchType == "1" {
		sql = fmt.Sprintf("select %s,ifnull(b.id, 0) as purcharse_id from book_list a left join book_purchase b on a.id=b.bid and b.uid=%d order by id desc limit %d,%d", fields, uid, (page-1)*pageSize, pageSize)
	} else {
		sql = fmt.Sprintf("select %s,b.id as purcharse_id from book_purchase b left join book_list a on b.bid=a.id where b.uid=%d order by b.id desc limit %d,%d", fields, uid, (page-1)*pageSize, pageSize)
	}

	rows, err := db.Query(sql)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var info model.BookInfo

		var imageLarge string
		var imageMedium string
		var imageSmall string
		var translatorStr string
		var authorStr string
		var dateline int
		var purcharseId int64
		err = rows.Scan(&info.Bid, &info.DoubanId, &info.Title, &info.OriginTitle, &info.AltTitle, &info.Publisher, &info.Pubdate, &info.Isbn13, &info.Isbn10, &translatorStr, &info.Pages, &authorStr, &info.AuthorIntro, &info.Catalog, &info.Summary, &info.Price, &info.Binding, &imageLarge, &imageMedium, &imageSmall, &dateline, &purcharseId)

		if err != nil {
			return result, err
		}

		info.Images = map[string]string{"large": imageLarge, "medium": imageMedium, "small": imageSmall}
		info.Author = strings.Split(authorStr, ",")
		info.Translator = strings.Split(translatorStr, ",")

		result.BookList = append(result.BookList, info)
		if purcharseId > 0 {
			result.PurchaseList[info.Bid] = purcharseId
		}
	}

	return result, nil
}
