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

func updateData(a int) error {
	strSql := "update testdata1 set mark=?,data=? where id=?"
	r, err := tx.Exec(strSql, "name_update", 9, a)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	i2, err2 := r.RowsAffected()
	if err != nil {
		fmt.Printf("err2: %v\n", err2)
		panic(err2)
	}
	fmt.Printf("i2: %v\n", i2)
	return nil
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

	insert_err := updateData(a)
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
	go mysql_transaction(1)
	//ctrl+c退出
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	/*SIGINT:os.Interrupt=ctrl+c  SIGTERM:kill函数产生: 相当于shell> kill不加-9时 pid
	SIGKILL：相当于shell> kill -9 pid
	*/
	s := <-exit
	logrus.Infoln("sigal return=", s)
}
