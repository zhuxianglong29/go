package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type user struct {
	id   int
	data int
}

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
	//开启事务
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {

			tx.Rollback()
		}
		fmt.Println("tx start err :", err)
		return
	}

	sqlStr := "select Id from testdata where id=?"
	var u user
	_ = tx.QueryRow(sqlStr, 1).Scan(&u.id, &u.data) //1表示查询id=1

	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("Id:%d data:%v ", u.id, u.data)
	time.Sleep(time.Second * 30)

	str := "insert into testdata(id,data) values (?,?)"
	data := map[string]uint16{
		"1": 1,
		"2": 2,
		"3": 3,
	}
	for key, vulues := range data {
		_, err := tx.Exec(str, key, vulues)
		if err != nil {
			fmt.Println("tx exec err :", err)
			tx.Rollback()
			return
		}
	}
	time.Sleep(time.Second * 3)
	err = tx.Commit()
	if err != nil {
		fmt.Println("提交错误，需要回滚！")
		return
	}
	fmt.Println("transaction success")
}

func main() {
	go insertsql()
	for {
	}

}
