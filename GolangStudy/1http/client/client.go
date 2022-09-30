package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	modbus "github.com/thinkgos/gomodbus/v2"
)

// 配置参数
var modbustcp_host string
var modbustcp_port string
var httpserver_host string
var httpserver_port string

// 从toml文件中读取配置
type (
	config struct {
		Os      string
		Version string
		Server  server
	}
	server struct {
		Modbustcp_host  string
		Modbustcp_port  string
		Httpserver_host string
		Httpserver_port string
	}
)

/*
	mbus函数

@details 通过modbustcp得到从机寄存器的地址
@param modbustcp_host 从机地址
@param modbustco_port 从机端口
@param out 管道，输入给数据
@retvar 无
*/
func mbus(modbustcp_host string, modbustcp_port string, out chan<- []byte) {

	//p := modbus.NewTCPClientProvider(modbustcp_slave+modbustcp_port, modbus.WithEnableLogger())
	p := modbus.NewTCPClientProvider(fmt.Sprint(modbustcp_host, modbustcp_port), modbus.WithEnableLogger())
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

/*
	myclient函数

@details http客户端，向服务器发送post请求
@param httpserver_host服务器地址
@param httpserver_port服务器端口
@param in 管道输出，从管道得到数据
@retvar 无
*/
func myclient(httpserver_host string, httpserver_port string, in <-chan []byte) {
	for {

		json_datemodbus := <-in //从channel中接收数据，并赋值给json_datemodbus
		logrus.Info(json_datemodbus)
		//url := httpserver_host + ":" + httpserver_port + "/products/"
		url := fmt.Sprint(httpserver_host, ":", httpserver_port, "/products/")
		contentType := "application/json"
		//client := http.Client{Timeout: 5 * time.Second}
		client := &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          1000,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 10 * time.Second,
			},
		}
		resp, err := client.Post(url, contentType, bytes.NewBuffer(json_datemodbus))

		if err != nil {
			logrus.Panic("post failed, err:%v\n", err)
		}
		logrus.Infoln("返回值", resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		logrus.Info("err=nil执行 err=", err)
		if err != nil {
			logrus.Panic("http err=", err)
		}

		logrus.Infoln(string(b))
		resp.Body.Close() //for循环中哦defer无效
	}
}

func main() {
	//读配置文件
	var tomlconfig config
	filePath := "client.toml"
	if _, err := toml.DecodeFile(filePath, &tomlconfig); err != nil {
		panic(err)
	}
	fmt.Println(tomlconfig)
	modbustcp_host = tomlconfig.Server.Modbustcp_host
	modbustcp_port = tomlconfig.Server.Modbustcp_port
	httpserver_host = tomlconfig.Server.Httpserver_host
	httpserver_port = tomlconfig.Server.Httpserver_port

	//日志打印到文件和命令行
	Stdout_writer := os.Stdout
	log_writer, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0755) //os.O_APPEND设置为打开文件， os.O_APPEND 如果存在就在尾部添加
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}
	logrus.SetOutput(io.MultiWriter(Stdout_writer, log_writer))

	for i := 0; i < 1000; i++ {

		i := make(chan []byte)
		go mbus(modbustcp_host, modbustcp_port, i)
		go myclient(httpserver_host, httpserver_port, i)
	}
	// i := make(chan []byte)
	// go mbus(modbustcp_host, modbustcp_port, i)
	// go myclient(httpserver_host, httpserver_port, i)
	// a := make(chan []byte)
	// go mbus(modbustcp_host, modbustcp_port, a)
	// go myclient(httpserver_host, httpserver_port, a)
	// b := make(chan []byte)
	// go mbus(modbustcp_host, modbustcp_port, b)
	// go myclient(httpserver_host, httpserver_port, b)
	//ctrl+c退出
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	/*SIGINT:os.Interrupt=ctrl+c  SIGTERM:kill函数产生: 相当于shell> kill不加-9时 pid
	SIGKILL：相当于shell> kill -9 pid
	*/
	s := <-exit
	logrus.Info("sigal return=", s)

}
