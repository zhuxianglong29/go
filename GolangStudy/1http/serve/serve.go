package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// mysql
// 定义一个全局对象db
var db *sql.DB

func insertsql(a string, b uint16) {
	// dsn := "root:root(密码)@tcp(127.0.0.1:3306（连接地址)）/go_db（数据库名）?charset=utf8mb4&parseTime=True"
	dsn := "root:666666@tcp(127.0.0.1:3306)/godata?charset=utf8mb4"
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
	r, err := db.Exec(sqlStr, a, b)
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

// 子路由, 分组路由
func myserve() {
	r := mux.NewRouter()

	//PathPrefix() 可以设置路由前缀，设置路由前缀为products
	products := r.PathPrefix("/products").Subrouter()
	//"http://localhost:8080/products/", 最后面的斜线一定要，不然路由不正确，页面出现404
	products.HandleFunc("/", ProductsHandler)
	//"http://localhost:8080/products/{key}"
	products.HandleFunc("/{key}", ProductHandler)

	users := r.PathPrefix("/users").Subrouter()
	// "/users"
	users.HandleFunc("/", UsersHandler)
	// "/users/id/参数/name/参数"
	users.HandleFunc("/id/{id}/name/{name}", UserHandler)

	http.ListenAndServe(":8080", r)
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

	for index, value := range cmd {
		a := index
		b := value
		insertsql(a, b)
		fmt.Println("a=", a, "b=", b)
		time.Sleep(time.Second)
	}
	// 	keys:=[]uint16{}
	// 	for key, _ := range cmd {
	// 		keys = append(keys, key)
	// 	}
	// 	sort.Strings(keys)
	// 	for _,key:= range keys{
	// 		insertsql(key, cmd[key])
	// 	}

	// 	keys := []string{}
	// 	for key, _ := range mp {
	// 		keys = append(keys, key)
	// 	}
	// 	sort.Strings(keys)
	// 	for _, key := range keys {
	// 		fmt.Println(key, " ----> ", mp[key])

	// }
}
func ProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //获取路由的值
	fmt.Fprintf(w, "key: %s", vars["key"])
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, " %s \r\n", "users handler")
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //获取值
	id := vars["id"]
	name := vars["name"]
	fmt.Fprintf(w, "id: %s, name: %s \r\n", id, name)
}

func main() {

	go myserve()

	for {

	}

}
