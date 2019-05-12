package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init()  {
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(192.168.43.115:3306)/fileserver?charset=utf8")
	if err != nil{
		panic(err)
	}
	db.SetConnMaxLifetime(1000)
	err = db.Ping()
	if err != nil{
		panic(err)
	}
}

// 返回数据库连接对象
func DBConn()*sql.DB  {
	return db
}