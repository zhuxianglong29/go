/*
单行查询:sqlstr  db.QueryRow Scan
多行查询:sqlstr		db.Query  r.close r.next
*/
package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 定义一个全局对象db
var db *sql.DB

func initDB() (err error) {
	// dsn := "root:root(密码)@tcp(127.0.0.1:3306（连接地址)）/go_db（数据库名）?charset=utf8mb4&parseTime=True"
	dsn := "root:123456@tcp(192.168.40.128:3306)/go_test?charset=utf8mb4"
	// open函数只是验证格式是否正确，并不是创建数据库连接
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 与数据库建立连接
	err2 := db.Ping()
	if err2 != nil {
		return err2
	}
	return nil
}

type user struct {
	id   int
	data int
}

func queryRowData() {
	sqlStr := "select Id from testdata where id=?"
	var u user
	err := db.QueryRow(sqlStr, 1).Scan(&u.id, &u.data) //1表示查询id=1
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("Id:%d data:%v ", u.id, u.data)
	time.Sleep(time.Second * 30)
}

// func queryManyData() {
// 	sqlStr := "select Id,username,password,status from user_tb1 where status > ?" //>?  <?  或不给参数全查询
// 	r, err := db.Query(sqlStr, 0)
// 	if err != nil {
// 		fmt.Printf("err: %v\n", err)
// 		return
// 	}
// 	defer r.Close()
// 	// 循环读取结果集中的数据
// 	for r.Next() {
// 		var u user
// 		err2 := r.Scan(&u.id, &u.username, &u.password, &u.status)
// 		if err2 != nil {
// 			fmt.Printf("err2: %v\n", err2)
// 			return
// 		}
// 		fmt.Printf("Id:%d username:%v password:%v status:%d\n", u.id, u.username, u.password, u.status)
// 	}
// }

func main() {
	err := initDB()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Println("连接成功")
	}
	fmt.Printf("db: %v\n", db)
	queryRowData()
	//queryManyData()
}
