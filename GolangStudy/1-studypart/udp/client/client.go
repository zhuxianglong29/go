package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 40000,
	})
	if err != nil {
		fmt.Println("conn failed err=", err)
		return
	}
	defer conn.Close()

	var reply [1024]byte
	reader := bufio.NewReader(os.Stdin)
	for {
		//输入数据
		fmt.Print("请输入内容：")
		msg, _ := reader.ReadString('\n')
		conn.Write([]byte(msg))
		//接受数据
		n, _, err := conn.ReadFromUDP(reply[:])
		if err != nil {
			fmt.Println("read msg err=", err)
			return
		}
		fmt.Println("收到信息：", string(reply[:n]))
	}
}
