package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"

	"GolangStudy/grpc/server/pb"

	_ "github.com/go-sql-driver/mysql"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// 配置参数，通过读取配置文件赋值
var user string
var passkey string
var mysql_addr string
var database_name string

//var server_listen_port string

type (
	config struct {
		Os      string
		Version string
		Server  server
	}
	server struct {
		User               string
		Passkey            string
		Mysql_adrr         string
		Database_name      string
		Server_listen_port string
	}
)

// mysql事务定义全局方便调用
var tx *sql.Tx

type mysql_port struct {
}

var p mysql_port

//mysql_port

/*
	insertsql函数

@details 在p方法中写mysql操作函数，具体如下
@param
@param
@param
@retvar 无
*/
//插入
func (p mysql_port) insertsql(b string) error {
	str := "insert into testdata2(data) values (?)"
	_, err := tx.Exec(str, b)
	return err
}

/*
mysql_transaction函数
@details mysql事物函数，通过调用mysql服务函数实现
@param
@param
@param
@retvar 无
*/
func mysql_transaction(a string) {
	//time.Sleep(time.Second * 60)//超时推出测试
	// dsn := "root:root(密码)@tcp(127.0.0.1:3306（连接地址)）/go_db（数据库名）?charset=utf8mb4&parseTime=True"
	dsn := fmt.Sprint(user, ":", passkey, "@tcp(", mysql_addr, ")/", database_name, "?charset=utf8mb4")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logrus.Panic(err)
	}
	// 与数据库建立连接
	err2 := db.Ping()
	if err2 != nil {
		logrus.Panic(err2)
	}
	logrus.Info("连接成功")
	//开启事务
	tx, err = db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		logrus.Panic("tx start err :", err)
	}
	logrus.Infof("insert a=%v %T", a, a)
	err_part := p.insertsql(a)
	if err_part != nil {
		tx.Rollback()
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logrus.Panic("提交错误，需要回滚！")
		//logrus.Info("提交错误，需要回滚！")

	}
	logrus.Info("transaction success")

}

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type grpc_server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *grpc_server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	mysql_transaction(in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func myserver() {
	// flag 读取命令行参数
	//flag.Parse()
	// 建立 TCP 连接
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 创建 gRPC 服务
	s := grpc.NewServer()
	// 注册服务，两个参数，将 server 结构体的方法进行注册
	pb.RegisterGreeterServer(s, &grpc_server{})
	log.Printf("server listening at %v", lis.Addr())
	// 运行服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
func main() {
	//读出配置文件
	var tomlconfig config
	filePath := "serve.toml"
	if _, err := toml.DecodeFile(filePath, &tomlconfig); err != nil {
		panic(err)
	}
	fmt.Println(tomlconfig) //调试用
	user = tomlconfig.Server.User
	passkey = tomlconfig.Server.Passkey
	mysql_addr = tomlconfig.Server.Mysql_adrr
	database_name = tomlconfig.Server.Database_name
	//server_listen_port = tomlconfig.Server.Server_listen_port
	//日志打印到文件和命令行
	Stdout_writer := os.Stdout
	log_writer, err := os.OpenFile("serverlog.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755) //os.O_APPEND设置为打开文件
	if err != nil {
		logrus.Panic("create file log.txt failed: %v", err)
	}
	logrus.SetOutput(io.MultiWriter(Stdout_writer, log_writer))

	//读配置文件
	// config, _ := toml.LoadFile("./serve.toml")
	// server_listen_port := config.Get("server.server_listen_port").(string)
	go myserver()

	//ctrl+c退出
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	/*SIGINT:os.Interrupt=ctrl+c  SIGTERM:kill函数产生: 相当于shell> kill不加-9时 pid
	SIGKILL：相当于shell> kill -9 pid
	*/
	s := <-exit
	logrus.Infoln("sigal return=", s)
}
