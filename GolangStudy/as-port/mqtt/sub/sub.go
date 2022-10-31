package main

import (
	"encoding/json"
	"fmt"
	"log"
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

type EUI64 [8]byte
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
	panic(err)
}

func main() {
	var tomlconfig config
	filePath := "zsub.toml"
	if _, err := toml.DecodeFile(filePath, &tomlconfig); err != nil {
		panic(err)
	}
	//fmt.Println(tomlconfig)

	conn, err := net.Dial("tcp", fmt.Sprint(tomlconfig.Local.Addr, ":", tomlconfig.Local.Port))

	if err != nil {
		fmt.Println("client dial err=", err)
		return
	}
	fmt.Println("client dial sucessful")
	go sub(tomlconfig.Remote.Mqttaddr, tomlconfig.Remote.Mqttport)

	go func() {
		for {

			data := <-datachan
			//fmt.Println("data=", data)

			n, err := conn.Write(data)
			if err != nil {
				log.Panic("conn.Write err=", err)
			}
			fmt.Printf("接收数据占位%d\n", n)
			//**

		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	/*SIGINT:os.Interrupt=ctrl+c  SIGTERM:kill函数产生: 相当于shell> kill不加-9时 pid
	SIGKILL：相当于shell> kill -9 pid
	*/
	s := <-exit
	logrus.Info("sigal return=", s)

}

func sub(a string, b string) {

	//会解析到models.DataDownPayload中，格式要匹配

	var broker = a
	var port = b
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
	// opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
	// 	fmt.Printf("Subscribe: Topic is [%s]; msg is [%s]\n", msg.Topic(), string(msg.Payload()))
	// 	datachan <- msg.Payload()
	// })

	topic := "application/1/device/2232330000888802/event/up"
	token := client.Subscribe(topic, 0, subCallBackFunc)
	token.Wait()
	//fmt.Printf("Subscribed to topic: %s", topic)

}

func subCallBackFunc(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Subscribe: Topic is [%s]; msg is [%s]\n", msg.Topic(), string(msg.Payload()))
	datachan <- msg.Payload()
}
