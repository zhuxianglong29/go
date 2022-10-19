package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sort"
	_ "time"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// 配置参数，通过读取配置文件赋值
var user string
var passkey string
var mysql_addr string
var database_name string
var server_listen_port string

type (
	config struct {
		Os      string
		Version string
		Server  server
	}
	server struct {
		User               string
		Passkey            string
		Mysql_adrr         string
		Database_name      string
		Server_listen_port string
	}
)

// mysql事务定义全局方便调用
var tx *sql.Tx

type mysql_port struct {
}

var p mysql_port

//mysql_port

/*
	insertsql函数

@details 在p方法中写mysql操作函数，具体如下
@param
@param
@param
@retvar 无
*/
//插入
func (p mysql_port) insertsql(a string, b uint16) error {
	str := "insert into testdata1(mark,data) values (?,?)"
	_, err := tx.Exec(str, a, b)
	return err
}

// 删除行
func (p mysql_port) delData(a int) error {
	strSql := "delete from testdata1 where id = ?"
	r, err := tx.Exec(strSql, a)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	i2, err2 := r.RowsAffected()
	if err2 != nil {
		fmt.Printf("err2: %v\n", err2)
		return err2
	}
	fmt.Printf("i2: %v\n", i2)
	return nil
}

// 查询行
func (p mysql_port) queryRowData(a int) error {
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

// 查询多行
func (p mysql_port) queryManyData1(a int) error {
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

// 更新数据单行
func (p mysql_port) updateData(a int) error {
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

/*
mysql_transaction函数
@details mysql事物函数，通过调用mysql服务函数实现
@param
@param
@param
@retvar 无
*/
func mysql_transaction(a string, b uint16) {
	//time.Sleep(time.Second * 60)//超时推出测试
	// dsn := "root:root(密码)@tcp(127.0.0.1:3306（连接地址)）/go_db（数据库名）?charset=utf8mb4&parseTime=True"
	dsn := fmt.Sprint(user, ":", passkey, "@tcp(", mysql_addr, ")/", database_name, "?charset=utf8mb4")
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
	//insert
	// str := "insert into testdata(id,data) values (?,?)"
	// r, err := tx.Exec(str, a, b)
	// if err != nil {
	// 	tx.Rollback()
	// 	logrus.Panic("tx exec err :", err)

	// }
	// i2, err2 := r.LastInsertId()
	// if err2 != nil {
	// 	tx.Rollback()
	// 	logrus.Panic("err2: %v\n", err2)
	// }
	// logrus.Info("i2: %v\n", i2)
	err_part := p.insertsql(a, b)
	if err_part != nil {
		tx.Rollback()
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		//logrus.Panic("提交错误，需要回滚！")
		logrus.Println("提交错误，需要回滚！")
		return //不能panic不然服务器会退从

	}
	logrus.Info("transaction success")

}

/*
	myserve函数

@details http服务端，监听请求，完成请求mysql服务
@param server_listen_port 监听的本机端口
@retvar 无
*/
func myserve(server_listen_port string) {
	r := mux.NewRouter()

	//PathPrefix() 可以设置路由前缀，设置路由前缀为products
	products := r.PathPrefix("/products").Subrouter()
	//"http://localhost:8080/products/", 最后面的斜线一定要，不然路由不正确，页面出现404
	products.HandleFunc("/", ProductsHandler)

	//users := r.PathPrefix("/users").Subrouter()

	//users.HandleFunc("/id/{id}/name/{name}", UserHandler)

	http.ListenAndServe(server_listen_port, r)
}

// products路由对应函数
func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	logrus.Info(w, "%s", "products")
	b, _ := io.ReadAll(r.Body)
	logrus.Info(string(b))
	//json转map
	var cmd map[string]uint16
	err := json.Unmarshal(b, &cmd)
	if err != nil {
		logrus.Panic("unmasharl err=", err)
	}
	logrus.Info("cmd=", cmd)

	// //上传mysql
	// //从配置文件读参数
	// config, _ := toml.LoadFile("./serve.toml")
	// user := config.Get("server.user").(string)
	// passkey := config.Get("server.passkey").(string)
	// mysql_addr := config.Get("server.mysql_addr").(string)
	// database_name := config.Get("server.database_name").(string)
	//按顺序读map
	keys := []string{}
	for key := range cmd {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		mysql_transaction(string(key), cmd[key])
		//insertsql(string(key), cmd[key])
	}

}

// func UserHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r) //获取值
// 	id := vars["id"]
// 	name := vars["name"]
// 	logrus.Info(w, "id: %s, name: %s \r\n", id, name)
// }

func main() {
	//读出配置文件
	var tomlconfig config
	filePath := "serve.toml"
	if _, err := toml.DecodeFile(filePath, &tomlconfig); err != nil {
		panic(err)
	}
	fmt.Println(tomlconfig) //调试用
	user = tomlconfig.Server.User
	passkey = tomlconfig.Server.Passkey
	mysql_addr = tomlconfig.Server.Mysql_adrr
	database_name = tomlconfig.Server.Database_name
	server_listen_port = tomlconfig.Server.Server_listen_port
	//日志打印到文件和命令行
	Stdout_writer := os.Stdout
	log_writer, err := os.OpenFile("serverlog.txt", os.O_WRONLY|os.O_CREATE, 0755) //os.O_APPEND设置为打开文件
	if err != nil {
		logrus.Panic("create file log.txt failed: %v", err)
	}
	logrus.SetOutput(io.MultiWriter(Stdout_writer, log_writer))

	//读配置文件
	// config, _ := toml.LoadFile("./serve.toml")
	// server_listen_port := config.Get("server.server_listen_port").(string)
	go myserve(server_listen_port)

	//ctrl+c退出
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	/*SIGINT:os.Interrupt=ctrl+c  SIGTERM:kill函数产生: 相当于shell> kill不加-9时 pid
	SIGKILL：相当于shell> kill -9 pid
	*/
	s := <-exit
	logrus.Infoln("sigal return=", s)

}
