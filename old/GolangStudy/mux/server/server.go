package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// 子路由, 分组路由
func a() {
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
	// for index, value := range cmd {
	// 	str := index
	// 	data := value
	// 	initDB(str, data)
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

}
