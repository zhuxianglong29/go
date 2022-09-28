package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	modbus "github.com/thinkgos/gomodbus/v2"
)

// modbus
func mbus(modbustcp_slave string, modbustcp_port string, out chan<- []byte) {

	//p := modbus.NewTCPClientProvider(modbustcp_slave+modbustcp_port, modbus.WithEnableLogger())
	p := modbus.NewTCPClientProvider(fmt.Sprint(modbustcp_slave, modbustcp_port), modbus.WithEnableLogger())
	client := modbus.NewClient(p)
	err := client.Connect()
	if err != nil {
		logrus.Error("connect failed, ", err)
		return
	}
	defer client.Close()

	logrus.Info("modbus starting")

	a, err := client.ReadHoldingRegisters(1, 0, 5) //从机编号,地址，数据量
	if err != nil {
		logrus.Error(err.Error())

	}
	logrus.Info(a)
	//转map
	datamodbus := make(map[string]uint16)
	for index, value := range a {
		str := strconv.Itoa(index)
		datamodbus[str] = value
	}
	logrus.Info("datamodbus %v\r\n", datamodbus)
	// json marshal
	json_datamodbus, err := json.Marshal(datamodbus)
	if err != nil {
		logrus.Panic("json err=", err)

	}
	logrus.Infoln("json:", string(json_datamodbus))
	out <- json_datamodbus
	time.Sleep(time.Second)
}

// http post client
func myclient(httpserver_host string, httpserver_port string, in <-chan []byte) {
	for {
		json_datemodbus := <-in //从channel中接收数据，并赋值给json_datemodbus
		logrus.Info(json_datemodbus)
		//url := httpserver_host + ":" + httpserver_port + "/products/"
		url := fmt.Sprint(httpserver_host, ":", httpserver_port, "/products/")
		contentType := "application/json"
		resp, err := http.Post(url, contentType, bytes.NewBuffer(json_datemodbus))
		if err != nil {
			logrus.Panic("post failed, err:%v\n", err)
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Panic(err)

		}
		logrus.Infoln(string(b))
		resp.Body.Close()
	}
}

func main() {
	//日志打印到文件和命令行
	Stdout_writer := os.Stdout
	log_writer, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0755) //os.O_APPEND设置为打开文件
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}
	logrus.SetOutput(io.MultiWriter(Stdout_writer, log_writer))

	config, _ := toml.LoadFile("./client.toml")
	//key := config.Get("os").(string)
	//fmt.Println(key)
	modbustcp_slave := config.Get("server.modbustcp_slave").(string)
	modbustcp_port := config.Get("server.modbustcp_port").(string)
	httpserver_host := config.Get("server.httpserver_host").(string)
	httpserver_port := config.Get("server.httpserver_port").(string)

	c := make(chan []byte)
	go mbus(modbustcp_slave, modbustcp_port, c)
	go myclient(httpserver_host, httpserver_port, c)

	//ctrl+c退出
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	/*SIGINT:os.Interrupt=ctrl+c  SIGTERM:kill函数产生: 相当于shell> kill不加-9时 pid
	SIGKILL：相当于shell> kill -9 pid
	*/
	s := <-exit
	logrus.Info("sigal return=", s)

}
