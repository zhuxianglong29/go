package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	modbus "github.com/thinkgos/gomodbus/v2"
)

// modbus
func mbus(out chan<- []byte) {

	p := modbus.NewTCPClientProvider("192.168.142.1:502", modbus.WithEnableLogger())
	client := modbus.NewClient(p)
	err := client.Connect()
	if err != nil {
		fmt.Println("connect failed, ", err)
		return
	}
	defer client.Close()

	fmt.Println("starting")

	a, err := client.ReadHoldingRegisters(1, 0, 5) //从机,地址，数据量
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(a)
	//转map
	datamodbus := make(map[string]uint16)
	for index, value := range a {
		str := strconv.Itoa(index)
		datamodbus[str] = value
	}
	fmt.Printf("datamodbus %v\r\n", datamodbus)

	time.Sleep(time.Second)

	// json marshal
	json_datamodbus, err := json.Marshal(datamodbus)
	if err != nil {
		fmt.Println("json err=", err)
		return
	}
	fmt.Println("json:", string(json_datamodbus))
	out <- json_datamodbus

}

// http post client
func myclient(in <-chan []byte) {
	for {
		json_datemodbus := <-in //从channel中接收数据，并赋值给x

		fmt.Println(json_datemodbus)

		url := "http://192.168.142.128:8080/products/"

		contentType := "application/json"

		resp, err := http.Post(url, contentType, bytes.NewBuffer(json_datemodbus))
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

func main() {
	//先读配置文件
	c := make(chan []byte)
	go mbus(c)
	//go mbus(modbustcp_host,port,c)
	go myclient(c)
	//go myclient(httpserver_host,port,c)
	for {

	}

}
