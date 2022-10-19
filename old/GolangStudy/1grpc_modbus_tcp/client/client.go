package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	modbus "github.com/thinkgos/gomodbus/v2"
	"google.golang.org/grpc"

	"GolangStudy/grpc/client/pb"

	"google.golang.org/grpc/credentials/insecure"
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
type Mymodbus struct {
}

var mbus Mymodbus
var string_datamodbus string

func (mbus Mymodbus) tcp(modbustcp_host string, modbustcp_port string, out chan<- string) {

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

	datamodbus, err := client.ReadHoldingRegisters(1, 0, 5) //从机编号,地址，数据量
	if err != nil {
		logrus.Error(err.Error())

	}
	logrus.Infof("a=%v  %T", datamodbus, datamodbus)
	//转map

	logrus.Info("datamodbus %v\r\n", datamodbus)
	// json marshal
	//json_datamodbus, err := json.Marshal(datamodbus)
	// if err != nil {
	// 	logrus.Panic("json err=", err)

	// }
	// logrus.Infoln("json:", string(json_datamodbus))
	str_datamodbus := fmt.Sprint(datamodbus)
	out <- str_datamodbus

}

// const (
// 	defaultName = "world"
// )

// var (
// 	addr = flag.String("addr", "localhost:50051", "the address to connect to")
// 	name = flag.String("name", defaultName, "Name to greet")
// )

// func mbus(modbustcp_host string, modbustcp_port string, out chan<- []byte) {

// 	//p := modbus.NewTCPClientProvider(modbustcp_slave+modbustcp_port, modbus.WithEnableLogger())
// 	p := modbus.NewTCPClientProvider(fmt.Sprint(modbustcp_host, modbustcp_port), modbus.WithEnableLogger())
// 	fmt.Println(fmt.Sprint(modbustcp_host, modbustcp_port))
// 	client := modbus.NewClient(p)
// 	err := client.Connect()
// 	if err != nil {
// 		logrus.Error("connect failed, ", err)
// 		return
// 	}
// 	defer client.Close()

// 	logrus.Info("modbus starting")

// 	a, err := client.ReadHoldingRegisters(1, 0, 5) //从机编号,地址，数据量
// 	if err != nil {
// 		logrus.Error(err.Error())

// 	}
// 	logrus.Info(a)
// 	//转map
// 	datamodbus := make(map[string]uint16)
// 	for index, value := range a {
// 		str := strconv.Itoa(index)
// 		datamodbus[str] = value
// 	}
// 	logrus.Info("datamodbus %v\r\n", datamodbus)
// 	// json marshal
// 	json_datamodbus, err := json.Marshal(datamodbus)
// 	if err != nil {
// 		logrus.Panic("json err=", err)

// 	}
// 	logrus.Infoln("json:", string(json_datamodbus))
// 	out <- json_datamodbus
// 	time.Sleep(time.Second)
// }

func grpc_client(httpserver_host string, httpserver_port string, in <-chan string) {
	//flag.Parse()
	// Set up a connection to the server.
	string_datamodbus = <-in
	addr := flag.String("addr", "localhost:50051", "the address to connect to")
	//name := flag.String("name", string_datamodbus, "Name to greet")

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: string_datamodbus})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
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
	log_writer, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755) //os.O_APPEND设置为打开文件， os.O_APPEND 如果存在就在尾部添加
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}
	logrus.SetOutput(io.MultiWriter(Stdout_writer, log_writer))

	c := make(chan string)

	go mbus.tcp(modbustcp_host, modbustcp_port, c)

	go grpc_client(httpserver_host, httpserver_port, c)

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
