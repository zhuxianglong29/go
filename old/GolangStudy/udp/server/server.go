package main

import (
	"fmt"
	"net"
	"strconv"
)

func jisheng(n int) (result int) {
	if n > 1 {
		result = n * jisheng(n-1)
		return result
	}
	return 1
}

func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 40000,
	})
	if err != nil {
		fmt.Println("listen UDP failed err=", err)
		return
	}
	defer conn.Close()

	//无需连接
	var buf [1024]byte
	for {
		n, addr, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			fmt.Println("read form UDP failed err=", err)
			return
		}
		fmt.Println(buf[:n])
		//数据处理
		strnub := string(buf[:n-1])
		fmt.Println("strnub=", strnub)
		k, _ := strconv.Atoi(strnub)
		fmt.Println("k=", k)
		j := strconv.Itoa(jisheng(k))
		//j:=jisheng(k)

		reply := []byte(j)
		fmt.Println("阶乘", jisheng(k))

		//发送数据
		//reply := strings.ToUpper(string(date[:n]))
		conn.WriteToUDP([]byte(reply), addr)
	}

}
