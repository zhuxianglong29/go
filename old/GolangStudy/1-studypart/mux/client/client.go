package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// 自定义tranpot
var (
	client = http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: false,
		},
	}
)

func main() {
	for {
		url := "http://192.168.40.128:8080/products/"
		// 表单数据
		//contentType := "application/x-www-form-urlencoded"
		//data := "name=小王子&age=18"
		// json
		reader := bufio.NewReader(os.Stdin)

		line, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("readString err=", err)
		}

		contentType := "application/json"
		//data := `{"name":"小王子","age":18}`
		resp, err := http.Post(url, contentType, strings.NewReader(line))
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
}
