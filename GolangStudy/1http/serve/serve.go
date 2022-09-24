package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sort"
	_ "time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/pelletier/go-toml"
)

// mysql
func insertsql(user string, key string, mysql_addr string, database_name string, a string, b uint16) {
	// dsn := "root:root(密码)@tcp(127.0.0.1:3306（连接地址)）/go_db（数据库名）?charset=utf8mb4&parseTime=True"

	// dsn := user + ":" + key + "@tcp(" + mysql_addr + ")/" + database_name + "?charset=utf8mb4"
	dsn := fmt.Sprint(user, ":", key, "@tcp(", mysql_addr, ")/", database_name, "?charset=utf8mb4")
	//func insertsql(a string, b uint16) {
	//dsn := "zxl:666666@tcp(192.168.1.48:3306)/godata?charset=utf8mb4"

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
	err = tx.Commit()
	if err != nil {
		fmt.Println("提交错误，需要回滚！")
		return
	}
	fmt.Println("transaction success")

}

// 子路由, 分组路由
func myserve(server_listen_port string) {
	r := mux.NewRouter()

	//PathPrefix() 可以设置路由前缀，设置路由前缀为products
	products := r.PathPrefix("/products").Subrouter()
	//"http://localhost:8080/products/", 最后面的斜线一定要，不然路由不正确，页面出现404
	products.HandleFunc("/", ProductsHandler)

	users := r.PathPrefix("/users").Subrouter()

	users.HandleFunc("/id/{id}/name/{name}", UserHandler)

	http.ListenAndServe(server_listen_port, r)
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", "products")
	b, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(b))
	//json转map
	var cmd map[string]uint16
	err := json.Unmarshal(b, &cmd)
	if err != nil {
		fmt.Println("unmasharl err=", err)
		return
	}

	//上传mysql
	fmt.Println("cmd=", cmd)

	//从配置文件读参数
	config, _ := toml.LoadFile("./serve.toml")
	user := config.Get("server.user").(string)
	passkey := config.Get("server.passkey").(string)
	mysql_addr := config.Get("server.mysql_addr").(string)
	database_name := config.Get("server.database_name").(string)

	//按顺序读map
	keys := []string{}
	for key := range cmd {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		insertsql(user, passkey, mysql_addr, database_name, string(key), cmd[key])
		//insertsql(string(key), cmd[key])
	}
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //获取值
	id := vars["id"]
	name := vars["name"]
	fmt.Fprintf(w, "id: %s, name: %s \r\n", id, name)
}

func main() {

	config, _ := toml.LoadFile("./serve.toml")
	server_listen_port := config.Get("server.server_listen_port").(string)
	go myserve(server_listen_port)

	//ctrl+c退出
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	/*SIGINT:os.Interrupt=ctrl+c  SIGTERM:kill函数产生: 相当于shell> kill不加-9时 pid
	SIGKILL：相当于shell> kill -9 pid
	*/
	s := <-exit
	fmt.Println("sigal return=", s)

}
