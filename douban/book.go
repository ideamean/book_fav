package douban

import (
	"book_fav/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetBookInfoByIsbn(isbn string) (model.BookInfo, error) {
	api := fmt.Sprintf("https://api.douban.com/v2/book/isbn/%s", isbn)
	var info model.BookInfo
	resp, err := http.Get(api)
	if err != nil {
		return info, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return info, err
	}

	err = json.Unmarshal(data, &info)
	if err != nil {
		return info, err
	}
	return info, err
}
