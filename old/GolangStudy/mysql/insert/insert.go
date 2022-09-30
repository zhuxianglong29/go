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

func insertsql(a string, b uint16) error {
	str := "insert into testdata1(mark,data) values (?,?)"
	r, err := tx.Exec(str, a, b)
	i, _ := r.LastInsertId()
	fmt.Printf("i: %v\n", i)
	return err
}

// 单行插入，多行可利用命令一次插入，具体见 https://blog.csdn.net/qq_39337886/article/details/123317292
func mysql_transaction(a string, b uint16) {
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

	insert_err := insertsql(a, b)
	if insert_err != nil {
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
	go mysql_transaction("1", 1)
	//ctrl+c退出
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	/*SIGINT:os.Interrupt=ctrl+c  SIGTERM:kill函数产生: 相当于shell> kill不加-9时 pid
	SIGKILL：相当于shell> kill -9 pid
	*/
	s := <-exit
	logrus.Infoln("sigal return=", s)
}
