package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var tx *sql.Tx

func queryRowData(a int) error {
	type user struct {
		id   int
		mark string
		data int
	}
	sqlStr := "select id,mark,data from testdata1 where id=?"
	var u user
	err := tx.QueryRow(sqlStr, a).Scan(&u.id, &u.mark, &u.data)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("Id:%d mark:%v data:%v", u.id, u.mark, u.data)
	return err
}
func queryManyData1(a int) error {
	type user struct {
		id   int
		mark string
		data int
	}
	sqlStr := "select id,mark,data from testdata1 where id>?" //从id大于开始读
	r, err := tx.Query(sqlStr, a)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	defer r.Close()
	users := make([]user, 0)
	// 循环读取结果集中的数据
	for r.Next() {
		var u user
		err2 := r.Scan(&u.id, &u.mark, &u.data)
		if err2 != nil {
			fmt.Printf("err2: %v\n", err2)
			return err2
		}
		users = append(users, u)
	}
	fmt.Printf("users: %+v\n", users)
	return err
}

func mysql_transaction(a int) {
	//time.Sleep(time.Second * 60)//超时推出测试
	// dsn := "root:root(密码)@tcp(127.0.0.1:3306（连接地址)）/go_db（数据库名）?charset=utf8mb4&parseTime=True"
	//dsn := fmt.Sprint(user, ":", passkey, "@tcp(", mysql_addr, ")/", database_name, "?charset=utf8mb4")
	dsn := "root:666666@tcp(192.168.0.181:3306)/godata?charset=utf8mb4"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logrus.Panic(err)
	}
	// 与数据库建立连接
	err2 := db.Ping()
	if err2 != nil {
		logrus.Panic(err2)
	}
	logrus.Info("连接成功")
	//开启事务
	tx, err = db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		logrus.Panic("tx start err :", err)
	}
	err_part := queryManyData1(a)
	if err_part != nil {
		tx.Rollback()
	}
	err = tx.Commit()
	if err != nil {
		logrus.Panic("提交错误，需要回滚！")
		tx.Rollback()
	}
	logrus.Info("transaction success")

}
func main() {
	go mysql_transaction(0)
	//ctrl+c退出
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	/*SIGINT:os.Interrupt=ctrl+c  SIGTERM:kill函数产生: 相当于shell> kill不加-9时 pid
	SIGKILL：相当于shell> kill -9 pid
	*/
	s := <-exit
	logrus.Infoln("sigal return=", s)
}
