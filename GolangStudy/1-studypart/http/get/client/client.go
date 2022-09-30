package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, err := http.Get("http://127.0.0.1:8000/123456")
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	resCode := resp.StatusCode
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read from resp.Body failed, err:%v\n", err)
		return
	}
	fmt.Println(resCode, "/n")
	fmt.Print(string(body))
}
