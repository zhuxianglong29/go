/*
连接:sql.Open  db.ping
insert操作: sqlstr  db.Exec
*/
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// 定义一个全局对象db
var db *sql.DB

func insertsql() {
	// dsn := "root:root(密码)@tcp(127.0.0.1:3306（连接地址)）/go_db（数据库名）?charset=utf8mb4&parseTime=True"
	dsn := "zxl:666666@tcp(192.168.1.48:3306)/godata?charset=utf8mb4"
	// open函数只是验证格式是否正确，并不是创建数据库连接

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 与数据库建立连接
	err2 := db.Ping()
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println("连接成功")

	fmt.Printf("db: %v\n", db)
	sqlStr := "insert into testdata(id,data) values (?,?)"
	r, err := db.Exec(sqlStr, "99", 99)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	i2, err2 := r.LastInsertId()
	if err2 != nil {
		fmt.Printf("err2: %v\n", err2)
		return
	}
	fmt.Printf("i2: %v\n", i2)

}
func main() {
	go insertsql()
	for {

	}

}
