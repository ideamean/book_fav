package dao

import (
	"database/sql"

	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego/config"
)

func GetDB() (*sql.DB, error) {
	iniconf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		return nil, err
	}
	host := iniconf.String("mysql::host")
	port := iniconf.String("mysql::port")
	user := iniconf.String("mysql::user")
	passwd := iniconf.String("mysql::passwd")
	dbname := iniconf.String("mysql::db")
	charset := iniconf.String("mysql::charset")
	timeout := iniconf.String("mysql::timeout")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&timeout=%ss", user, passwd, host, port, dbname, charset, timeout)
	db, err := sql.Open("mysql", dsn)
	return db, err
}
