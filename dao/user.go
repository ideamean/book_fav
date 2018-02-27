package dao

import (
	"book_fav/model"
	"errors"
	"strings"
	"time"
)

//根据邮箱获取用户信息
func GetUserInfoByEmail(email string) (model.UserInfo, error) {
	var info model.UserInfo
	db, err := GetDB()
	if err != nil {
		return info, err
	}

	defer db.Close()

	var id int64
	var nickname string
	var dateline int
	err = db.QueryRow("select id,nickname,dateline from user_list where email=?", email).Scan(&id, &nickname, &dateline)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err := AddUser(email)
			if err != nil {
				return info, errors.New("添加用户失败")
			}
			return GetUserInfoByEmail(email)
		} else {
			return info, err
		}
	}
	info.Id = id
	info.Email = email
	info.NickName = nickname
	info.Dateline = dateline
	return info, nil
}

//根据邮箱获取用户信息
func GetUserInfoByUid(uid int64) (model.UserInfo, error) {
	var info model.UserInfo
	db, err := GetDB()
	if err != nil {
		return info, err
	}

	defer db.Close()

	var id int64
	var nickname string
	var email string
	var dateline int
	err = db.QueryRow("select id,email,nickname,dateline from user_list where id=?", uid).Scan(&id, &email, &nickname, &dateline)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return info, nil
		} else {
			return info, err
		}
	}
	info.Id = id
	info.Email = email
	info.NickName = nickname
	info.Dateline = dateline
	return info, nil
}

//添加用户
func AddUser(email string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	defer db.Close()

	stmt, err := db.Prepare("insert into user_list set email=?,nickname=?,dateline=?")
	if err != nil {
		return err
	}

	defer stmt.Close()
	nicknameArr := strings.Split(email, "@")
	_, err = stmt.Exec(email, nicknameArr[0], time.Now().Unix())
	return err
}
