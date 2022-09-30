package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.40.128:8888")
	if err != nil {
		fmt.Println("client dial err=", err)
		return
	}

	for {
		reader := bufio.NewReader(os.Stdin)

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("readString err=", err)
		}

		n, err := conn.Write([]byte(line))
		if err != nil {
			fmt.Println("conn.Write err=", err)
		}
		fmt.Printf("接收数据占位%d\n", n)
		//**
		buf := make([]byte, 1024)
		conn.Read(buf)
		fmt.Println(string(buf))

		//fmt.Printf("\n%v",buf)
	}

}
