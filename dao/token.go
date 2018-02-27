package dao

import (
	"book_fav/model"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"time"
)

//添加token,返回genCode, error
func AddToken(uid int64) (int64, error) {
	db, err := GetDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	rand, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return 0, err
	}

	genCode := rand.Int64()

	tokenStr := fmt.Sprintf("uid_%d_%d", time.Now().UnixNano(), genCode)
	h := md5.New()
	io.WriteString(h, tokenStr)
	token := fmt.Sprintf("%x", h.Sum(nil))

	stmt, err := db.Prepare("insert into user_token set token=?,uid=?,gen_code=?,expire=?,status=0,dateline=?")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(token, uid, genCode, 0, time.Now().Unix())
	if err != nil {
		return 0, err
	}
	return genCode, nil
}

//根据uid,token验证登陆
func CheckTokenByUid(uid int64, genCode int) (model.TokenInfo, error) {
	var info model.TokenInfo
	db, err := GetDB()
	if err != nil {
		return info, err
	}

	defer db.Close()

	var id int64
	var token string
	var expire int64
	var status int
	var dateline int
	err = db.QueryRow("select id,token,expire,status,dateline from user_token where uid=? and gen_code=? and status=0", uid, genCode).Scan(&id, &token, &expire, &status, &dateline)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return info, nil
		} else {
			return info, err
		}
	}
	//1小时内过期自动延期到2个月过期
	if time.Now().Unix()-3600 > expire || status == 0 {
		expire = time.Now().Unix() + 86400*60
		stmt, err := db.Prepare("update user_token set status=1,expire=? where id=?")
		if err != nil {
			return info, err
		}
		_, err = stmt.Exec(expire, id)
		if err != nil {
			return info, err
		}
	}
	info.Id = id
	info.Token = token
	info.Uid = uid
	info.GenCode = genCode
	info.Expire = expire
	info.Status = status
	info.Dateline = dateline
	return info, nil
}

//获取token
func GetTokenInfoByToken(token string) (model.TokenInfo, error) {
	var info model.TokenInfo
	db, err := GetDB()
	if err != nil {
		return info, err
	}

	defer db.Close()

	var id int64
	var uid int64
	var genCode int
	var expire int64
	var status int
	var dateline int
	err = db.QueryRow("select id,uid,gen_code,expire,status,dateline from user_token where token=?", token).Scan(&id, &uid, &genCode, &expire, &status, &dateline)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return info, nil
		} else {
			return info, err
		}
	}
	info.Id = id
	info.Token = token
	info.Uid = uid
	info.GenCode = genCode
	info.Expire = expire
	info.Status = status
	info.Dateline = dateline
	return info, nil
}
