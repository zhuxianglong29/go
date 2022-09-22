package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func client(a string) {
	url := "http://192.168.142.128:9090/post"
	// 表单数据
	//contentType := "application/x-www-form-urlencoded"
	//data := "name=小王子&age=18"
	// json
	contentType := "application/json"
	data := a

	resp, err := http.Post(url, contentType, strings.NewReader(data))
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v\n", err)
		return
	}
	fmt.Println(string(b))

}

func main() {
	go client("21")
	go client("22")
	go client("23")
	go client("24")
	go client("25")

	for {

	}

}
