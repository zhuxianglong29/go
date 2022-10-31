package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/BurntSushi/toml"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

type config struct {
	Os      string
	Version string
	Local   local
	Remote  remote
}
type local struct {
	Addr string
	Port string
}
type remote struct {
	Mqttaddr string `toml:"mqttaddr"`
	Mqttport string
}

type EUI64 [16]byte // 格式会不对？？？
type DataDownPayload struct {
	ApplicationID int64           `json:"applicationID,string"`
	DevEUI        EUI64           `json:"devEUI"`
	Confirmed     bool            `json:"confirmed"`
	FPort         uint8           `json:"fPort"`
	Data          []byte          `json:"data"`
	Object        json.RawMessage `json:"object"`
}

var datachan = make(chan []byte)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {
	var tomlconfig config
	filePath := "zpub.toml"
	if _, err := toml.DecodeFile(filePath, &tomlconfig); err != nil {
		panic(err)
	}
	fmt.Println(tomlconfig)

	fmt.Println("服务器开始监")
	listen, err := net.Listen("tcp", fmt.Sprint(tomlconfig.Local.Addr, ":", tomlconfig.Local.Port))
	if err != nil {
		fmt.Println("listen err= ", err)
		return
	}
	defer listen.Close()
	fmt.Println("等待连接...")
	conn, err := listen.Accept()
	if err != nil {
		fmt.Println("Acceept err=", err)
	} else {
		fmt.Printf("Accept suc conn=%v 客户端ip=%v\n", conn, conn.RemoteAddr().String())
	}
	//等待链接
	go process(conn)

	//sub(client)

	go publish(tomlconfig.Remote.Mqttaddr, tomlconfig.Remote.Mqttport)

	//client.Disconnect(100000)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	/*SIGINT:os.Interrupt=ctrl+c  SIGTERM:kill函数产生: 相当于shell> kill不加-9时 pid
	SIGKILL：相当于shell> kill -9 pid
	*/
	s := <-exit
	logrus.Info("sigal return=", s)
}

func process(conn net.Conn) {

	for {
		//协程-客户

		buf := make([]byte, 20)
		fmt.Printf("服务器在等待客户%s 发送信息\n", conn.RemoteAddr().String())
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("客户端退出err=%v", err)
			return
		}
		fmt.Println(n)
		fmt.Println(buf)
		datachan <- buf[:n-1]
		//3.显示客户端发送内容
		//fmt.Print(string(buf[:n-1]))
		//fmt.Printf("buf格式%T %v\n", buf, buf[:n])

	}
}

func publish(confaddr string, confport string) {
	for {
		//会解析到models.DataDownPayload中，格式要匹配
		data := <-datachan
		var broker = confaddr
		var port = confport
		opts := mqtt.NewClientOptions()
		opts.AddBroker(fmt.Sprintf("tcp://%s:%s", broker, port))
		opts.SetClientID("2232330000888802")
		//opts.SetUsername("zxl")
		//opts.SetPassword("666666")
		opts.SetDefaultPublishHandler(messagePubHandler)
		opts.OnConnect = connectHandler
		opts.OnConnectionLost = connectLostHandler
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		pl := DataDownPayload{
			//DevEUI:    0000000000000000,
			Confirmed: false,
			FPort:     2,
			Data:      data,
		}
		fmt.Println("pl", pl)
		//b := pl //必须是json格式
		b, err := json.Marshal(pl)
		if err != nil {
			fmt.Println("marshal err ", err)
		}
		fmt.Println("b", b)

		token := client.Publish("application/1/device/2232330000888802/command/down", 0, false, b) //application/+/device/+/command/down
		token.Wait()
		client.Disconnect(1)
	}
}

// func sub(client mqtt.Client) {
// 	topic := "application/1/device/2232330000888802/event/up"
// 	token := client.Subscribe(topic, 2, nil)
// 	token.Wait()
// 	fmt.Printf("Subscribed to topic: %s", topic)
// }
